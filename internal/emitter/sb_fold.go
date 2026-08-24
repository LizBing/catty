package emitter

import (
	"fmt"
	"os"
	"strings"

	"catty/internal/classfile"
)

// StringBuilder-chain folding (P-0009 U1): recognize the javac concat
// shape
//
//	new StringBuilder; dup; [<push>; invokespecial <init>(String|)]…
//	(<push>; invokevirtual append(X)…)* invokevirtual toString()
//
// and emit one Go concatenation ending in a fresh JString (non-interned,
// identity semantics preserved). Guarded windows only: branch targets,
// exception ranges, stack-slot operands and float/double formatting abort.
//
// Temporaries (`_sb<pc>`) are declared in the method prologue — the
// generated body is a flat goto graph, so mid-function declarations would
// be jumped over illegally.

type sbArg struct {
	lit    string // non-empty ⇒ direct Go string constant
	opnd   string // typed Go expression producing the operand
	typ    byte   // 'S' | 'I' | 'J' | 'C' | 'Z'
	size   int    // instruction bytes
	cat2   bool
	_guard bool
}

type sbFold struct {
	end  int    // last pc of the window (toString invoke)
	out  int    // stack depth before `new`; result lands at s[out]
	stmt string // assignments to _sb<start>
}

func (a sbArg) render() string {
	switch a.typ {
	case 'S':
		if a.lit != "" {
			return "_sb += " + a.lit + "\n"
		}
		return "_sb += genrt.StrOf(" + a.opnd + ")\n"
	case 'I':
		return "_sb += genrt.ItoA(" + a.opnd + ")\n"
	case 'J':
		return "_sb += genrt.JtoA(" + a.opnd + ")\n"
	case 'C', 'Z':
		// opnd keeps its ".(int32)" assertion — the helpers take int32.
		return "_sb += genrt." + map[byte]string{'C': "CtoA", 'Z': "ZtoA"}[a.typ] +
			"(" + a.opnd + ")\n"
	}
	return ""
}

// pushSBArg classifies whitelisted single-push instructions as concat
// arguments.
func (e *methodEmitter) pushSBArg(pc int) (sbArg, bool) {
	code := e.code
	op := code[pc]
	switch op {
	case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08:
		return sbArg{opnd: fmt.Sprint(int32(op - 3)), typ: 'I', size: 1}, true
	case 0x10:
		return sbArg{opnd: fmt.Sprint(int32(int8(code[pc+1]))), typ: 'I', size: 2}, true
	case 0x11:
		v := int32(int16(uint16(code[pc+1])<<8 | uint16(code[pc+2])))
		return sbArg{opnd: fmt.Sprint(v), typ: 'I', size: 3}, true
	case 0x12:
		cnst := &e.cf.Constants[code[pc+1]]
		switch cnst.Tag {
		case classfile.CInteger:
			return sbArg{opnd: fmt.Sprint(cnst.IntVal), typ: 'I', size: 2}, true
		case classfile.CString:
			if sv, err := e.cf.UTF8(cnst.Idx1); err == nil {
				return sbArg{lit: fmt.Sprintf("%q", sv), typ: 'S', size: 2}, true
			}
		}
		return sbArg{}, false
	case 0x13:
		idx := uint16(code[pc+1])<<8 | uint16(code[pc+2])
		cnst := &e.cf.Constants[idx]
		if cnst.Tag == classfile.CString {
			if sv, err := e.cf.UTF8(cnst.Idx1); err == nil {
				return sbArg{lit: fmt.Sprintf("%q", sv), typ: 'S', size: 3}, true
			}
		}
		return sbArg{}, false
	case 0x15:
		return sbArg{opnd: fmt.Sprintf("l%d.(int32)", code[pc+1]), typ: 'I', size: 2}, true
	case 0x1a, 0x1b, 0x1c, 0x1d:
		return sbArg{opnd: fmt.Sprintf("l%d.(int32)", op-0x1a), typ: 'I', size: 1}, true
	case 0x19:
		return sbArg{opnd: fmt.Sprintf("l%d", code[pc+1]), typ: 'S', size: 2}, true
	case 0x2a, 0x2b, 0x2c, 0x2d:
		return sbArg{opnd: fmt.Sprintf("l%d", op-0x2a), typ: 'S', size: 1}, true
	case 0x16:
		return sbArg{opnd: fmt.Sprintf("l%d.(int64)", code[pc+1]), typ: 'J', size: 2, cat2: true}, true
	case 0x1e, 0x1f, 0x20, 0x21:
		return sbArg{opnd: fmt.Sprintf("l%d.(int64)", op-0x1e), typ: 'J', size: 1, cat2: true}, true
	}
	return sbArg{}, false
}

// sbValue parses one concat argument under bytecode evaluation order:
// push {push binop}* — left operand, right operand, operator, …
//
// If a trailing construct isn't ours (e.g. the "operator" slot holds
// something else), we return just what matched; the caller then sees a
// non-invoke byte and bails the window — no wrong folds, no rollback.
// Division/remainder require nonzero literal rhs so ArithmeticException
// origins stay put.
func (e *methodEmitter) sbValue(pc, depth int) (sbArg, bool) {
	if depth > 3 {
		return sbArg{}, false
	}
	cur, ok := e.pushSBArg(pc)
	if !ok {
		return sbArg{}, false
	}
	pos := pc + cur.size
	for depth <= 3 && pos < len(e.code) && isPushOpcode(e.code[pos]) {
		rhs, ok2 := e.sbValue(pos, depth+1)
		if !ok2 {
			break
		}
		op := byte(0xFF)
		if p := pos + rhs.size; p < len(e.code) {
			op = e.code[p]
		}
		var opStr string
		switch op {
		case 0x60:
			opStr = "+"
		case 0x64:
			opStr = "-"
		case 0x68:
			opStr = "*"
		case 0x70:
			opStr = "%"
		case 0x7e:
			opStr = "&"
		case 0x80:
			opStr = "|"
		case 0x82:
			opStr = "^"
		case 0x7a, 0x7c:
			if rhs.typ != 'I' {
				return cur, true
			}
			var expr string
			if op == 0x7a {
				expr = "(" + cur.opnd + " >> (" + rhs.opnd + " & 31))"
			} else {
				expr = "(int32)(uint32(" + cur.opnd + ") >> (uint32(" + rhs.opnd + ") & 31))"
			}
			cur = sbArg{opnd: expr, typ: 'I', size: cur.size + rhs.size + 1}
			pos += rhs.size + 1
			continue
		default:
			return cur, true // not an operator: trailing bytes are outer context
		}
		if rhs.typ != 'I' {
			return cur, true
		}
		if (op == 0x70 || op == 0x68) &&
			!(len(rhs.opnd) > 0 && rhs.opnd[0] >= '0' && rhs.opnd[0] <= '9') {
			return cur, true // possibly-zero divisor: keep native semantics
		}
		cur = sbArg{opnd: "(" + cur.opnd + " " + opStr + " " + rhs.opnd + ")",
			typ: 'I', size: cur.size + rhs.size + 1}
		pos += rhs.size + 1
	}
	return cur, true
}

func isPushOpcode(op byte) bool {
	switch op {
	case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x10, 0x11, 0x12, 0x13,
		0x15, 0x16, 0x19,
		0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21,
		0x2a, 0x2b, 0x2c, 0x2d:
		return true
	}
	return false
}
func (e *methodEmitter) findSBFolds() map[int]*sbFold {
	folds := map[int]*sbFold{}
	code := e.code
	if os.Getenv("CATTY_SBDBG") != "" {
		fmt.Fprintf(os.Stderr, "[sbfold] scan %s.%s\n", e.cf.ThisClass, e.m.Name)
	}

	inHandlerRange := func(pc int) bool {
		for _, h := range e.handlers {
			if pc >= int(h.StartPc) && pc < int(h.EndPc) {
				return true
			}
		}
		return false
	}
	guarded := func(pc int) bool {
		return !e.reach[pc] || e.targets[pc] || inHandlerRange(pc)
	}
	refInvoke := func(pc int) (name, desc string, ok bool) {
		if op := code[pc]; op != 0xb6 && op != 0xb7 {
			return "", "", false
		}
		_, nm, ds, _, err := refAt(e.cf, code, pc)
		if err != nil {
			return "", "", false
		}
		return nm, ds, true
	}

	for start := range e.reach {
		if code[start] != 0xbb {
			continue
		}
		clsName, cerr := classAt(e.cf, code, start)
		if cerr != nil || clsName != "java/lang/StringBuilder" {
			continue
		}
		pc := start + 3
		if os.Getenv("CATTY_SBDBG") != "" {
			fmt.Fprintf(os.Stderr, "[sbfold] try %s.%s start=%d\n", e.cf.ThisClass, e.m.Name, start)
		}
		if pc >= len(code) || code[pc] != 0x59 || guarded(pc) {
			if os.Getenv("CATTY_SBDBG") != "" {
				fmt.Fprintf(os.Stderr, "[sbfold]   dup reject at pc=%d code=%#x guarded=%v\n", pc, safeOp(code, pc), guarded(pc))
			}
			continue // expect dup
		}
		pc++

		var stmts []string
		if nm, ds, ok := refInvoke(pc); ok && nm == "<init>" && ds == "()V" {
			pc += 3
		} else {
			// Growth-chain guard: a non-literal initial string means the
			// chain ACCUMULATES (acc + delta per loop iteration). Go `+=`
			// re-copies the whole prefix each round — strictly worse than
			// the incremental UTF-16 StringBuilder natives (measured 23x
			// regression on Bench.strcat). Fold only fresh-build chains.
			arg, ok2 := e.sbValue(pc, 0)
			if !ok2 || arg.typ != 'S' || arg.lit == "" || guarded(pc) {
				continue
			}
			pc += arg.size
			nm, ds, ok3 := refInvoke(pc)
			if !ok3 || nm != "<init>" || ds != "(Ljava/lang/String;)V" || guarded(pc) {
				continue
			}
			stmts = append(stmts, arg.render())
			pc += 3
		}

		okFold := true
		for {
			if os.Getenv("CATTY_SBDBG") != "" {
				fmt.Fprintf(os.Stderr, "[sbfold]   loop enter pc=%d op=%#x\n", pc, safeOp(code, pc))
			}
			if nm, _, ok := refInvoke(pc); ok && nm == "toString" {
				break
			}
			arg, ok2 := e.sbValue(pc, 0)
			if !ok2 || guarded(pc) {
				okFold = false
				break
			}
			pc += arg.size
			nm, ds, ok := refInvoke(pc)
			if !ok || nm != "append" || guarded(pc) || !hasSuffixBuilder(ds) {
				okFold = false
				break
			}
			closeIdx := indexByteStr(ds, ')')
			if closeIdx < 0 {
				okFold = false
				break
			}
			argPart := ds[1:closeIdx]
			// Normalize descriptor form → sbArg type letter.
			want := argPart
			switch argPart {
			case "Ljava/lang/String;":
				want = "S"
			case "B", "S":
				want = "I" // sub-int args arrive as int pushes
			}
			if want == "D" || want == "F" || strings.HasPrefix(argPart, "[") ||
				strings.HasPrefix(want, "L") || want == "R" {
				okFold = false
				break // formatting / identity risk
			}
			if (arg.cat2 && want != "J") || (!arg.cat2 && want == "J") {
				okFold = false
				break
			}
			if want != string(arg.typ) {
				// append(C)/append(Z) receive an int-typed push; render()
				// strips the ".(int32)" cast for those.
				if !(want == "C" || want == "Z") || arg.typ != 'I' {
					okFold = false
					break
				}
				arg.typ = want[0]
			}
			stmts = append(stmts, arg.render())
			pc += 3
		}
		if !okFold {
			continue
		}
		nm, ds, ok := refInvoke(pc)
		if !ok || nm != "toString" || ds != "()Ljava/lang/String;" || guarded(pc) {
			continue
		}
		d, hasD := e.canonDepths[start]
		if !hasD {
			continue
		}
		folds[start] = &sbFold{end: pc + 3 - 1, out: d, stmt: strings.Join(stmts, "")}
	}
	return folds
}

func hasSuffixBuilder(ds string) bool {
	const suf = ")Ljava/lang/StringBuilder;"
	return len(ds) >= len(suf) && ds[len(ds)-len(suf):] == suf
}

func indexByteStr(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func safeOp(code []byte, pc int) byte {
	if pc >= 0 && pc < len(code) {
		return code[pc]
	}
	return 0xFF
}
