package kernel

import "catty/internal/classfile"

// bootstrap assembles the synthesized java.* surface required by M0.
// Definition order matters: supers and field types before dependents.
func bootstrap(k *Kernel) {
	mustDefine(k, &ClassDef{
		Name:   "java/lang/Object",
		Flags:  classfile.AccPublic,
		Fields: nil,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natObjectInit},
			{Name: "hashCode", Desc: "()I", Flags: classfile.AccPublic, Native: natObjectHashCode},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natObjectEquals},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natObjectToString},
			{Name: "wait", Desc: "(J)V", Flags: classfile.AccPublic|classfile.AccFinal, Native: natObjectWaitMillis},
			{Name: "notify", Desc: "()V", Flags: classfile.AccPublic|classfile.AccFinal, Native: natObjectNotify},
			{Name: "notifyAll", Desc: "()V", Flags: classfile.AccPublic|classfile.AccFinal, Native: natObjectNotifyAll},
		},
	})

	mustDefine(k, &ClassDef{
		Name:  "java/lang/String",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Methods: []MethodDef{
			{Name: "length", Desc: "()I", Flags: classfile.AccPublic, Native: natStringLength},
			{Name: "charAt", Desc: "(I)C", Flags: classfile.AccPublic, Native: natStringCharAt},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natStringEquals},
			{Name: "hashCode", Desc: "()I", Flags: classfile.AccPublic, Native: natStringHashCode},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringToString},
		},
	})

	mustDefine(k, &ClassDef{
		Name:  "java/lang/StringBuilder",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natStringBuilderInit},
			appendSB("Ljava/lang/String;", natSBAppendString),
			appendSB("I", natSBAppendInt),
			appendSB("J", natSBAppendLong),
			appendSB("C", natSBAppendChar),
			appendSB("Z", natSBAppendBool),
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natSBToString},
		},
	})

	throwables := []struct{ name, super string }{
		{"java/lang/Throwable", "java/lang/Object"},
		{"java/lang/Exception", "java/lang/Throwable"},
		{"java/lang/Error", "java/lang/Throwable"},
		{"java/lang/RuntimeException", "java/lang/Exception"},
		{"java/lang/ArithmeticException", "java/lang/RuntimeException"},
		{"java/lang/IndexOutOfBoundsException", "java/lang/RuntimeException"},
		{"java/lang/ArrayIndexOutOfBoundsException", "java/lang/IndexOutOfBoundsException"},
		{"java/lang/StringIndexOutOfBoundsException", "java/lang/IndexOutOfBoundsException"},
		{"java/lang/NullPointerException", "java/lang/RuntimeException"},
		{"java/lang/ClassCastException", "java/lang/RuntimeException"},
		{"java/lang/NegativeArraySizeException", "java/lang/RuntimeException"},
		{"java/lang/IllegalArgumentException", "java/lang/RuntimeException"},
		{"java/lang/IllegalStateException", "java/lang/RuntimeException"},
		{"java/lang/UnsupportedOperationException", "java/lang/RuntimeException"},
		{"java/lang/IllegalMonitorStateException", "java/lang/RuntimeException"},
		{"java/lang/InterruptedException", "java/lang/Exception"},
	}
	for _, t := range throwables {
		def := &ClassDef{
			Name:  t.name,
			Super: t.super,
			Flags: classfile.AccPublic,
		}
		if t.name == "java/lang/Throwable" {
			def.Fields = []FieldDef{
				{Name: "detailMessage", Desc: "Ljava/lang/String;", Flags: classfile.AccPrivate},
			}
			def.Methods = append(def.Methods,
				MethodDef{Name: "getMessage", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natThrowableGetMessage},
				MethodDef{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natThrowableToString},
			)
		}
		def.Methods = append(def.Methods,
			MethodDef{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natThrowableInitVoid},
			MethodDef{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natThrowableInitString},
		)
		mustDefine(k, def)
	}

	mustDefine(k, &ClassDef{
		Name:  "java/lang/Integer",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Fields: []FieldDef{
			{Name: "value", Desc: "I", Flags: classfile.AccPrivate | classfile.AccFinal},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(I)V", Flags: classfile.AccPublic, Native: natIntegerInit},
			{Name: "intValue", Desc: "()I", Flags: classfile.AccPublic, Native: natIntegerIntValue},
			{Name: "hashCode", Desc: "()I", Flags: classfile.AccPublic, Native: natIntegerHashCode},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natIntegerEquals},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natIntegerToString},
			{Name: "valueOf", Desc: "(I)Ljava/lang/Integer;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticIntegerValueOf},
			{Name: "toString", Desc: "(I)Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticIntegerToString},
		},
	})

	mustDefine(k, &ClassDef{
		Name:  "java/util/ArrayList",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natArrayListInit},
			{Name: "add", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natArrayListAdd},
			{Name: "get", Desc: "(I)Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natArrayListGet},
			{Name: "set", Desc: "(ILjava/lang/Object;)Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natArrayListSet},
			{Name: "size", Desc: "()I", Flags: classfile.AccPublic, Native: natArrayListSize},
			{Name: "isEmpty", Desc: "()Z", Flags: classfile.AccPublic, Native: natArrayListIsEmpty},
			{Name: "contains", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natArrayListContains},
		},
	})

	mustDefine(k, &ClassDef{
		Name:  "java/io/PrintStream",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "println", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natPrintlnString},
			{Name: "println", Desc: "(I)V", Flags: classfile.AccPublic, Native: natPrintlnInt},
			{Name: "println", Desc: "(J)V", Flags: classfile.AccPublic, Native: natPrintlnLong},
			{Name: "println", Desc: "(C)V", Flags: classfile.AccPublic, Native: natPrintlnChar},
			{Name: "println", Desc: "(Z)V", Flags: classfile.AccPublic, Native: natPrintlnBool},
			{Name: "println", Desc: "(Ljava/lang/Object;)V", Flags: classfile.AccPublic, Native: natPrintlnObject},
			{Name: "println", Desc: "()V", Flags: classfile.AccPublic, Native: natPrintlnVoid},
			{Name: "print", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natPrintString},
			{Name: "print", Desc: "(Ljava/lang/Object;)V", Flags: classfile.AccPublic, Native: natPrintObject},
		},
	})

	mustDefine(k, &ClassDef{
		Name:       "java/lang/System",
		Super:      "java/lang/Object",
		Flags:      classfile.AccPublic | classfile.AccFinal,
		StaticInit: systemOutInit,
		Fields: []FieldDef{
			{Name: "out", Desc: "Ljava/io/PrintStream;",
				Flags: classfile.AccPublic | classfile.AccStatic | classfile.AccFinal},
		},
		Methods: []MethodDef{
			{Name: "arraycopy", Desc: "(Ljava/lang/Object;ILjava/lang/Object;II)V",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natSystemArraycopy},
			{Name: "currentTimeMillis", Desc: "()J",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natCurrentTimeMillis},
			{Name: "nanoTime", Desc: "()J",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natNanoTime},
			{Name: "identityHashCode", Desc: "(Ljava/lang/Object;)I",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natIdentityHashCode},
		},
	})
}

func appendSB(desc string, nat NativeFunc) MethodDef {
	return MethodDef{
		Name:   "append",
		Desc:   "(" + desc + ")Ljava/lang/StringBuilder;",
		Flags:  classfile.AccPublic,
		Native: nat,
	}
}

func mustDefine(k *Kernel, def *ClassDef) {
	if _, err := k.DefineClass(def); err != nil {
		panic("kernel bootstrap: " + err.Error())
	}
}
