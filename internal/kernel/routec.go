package kernel

import (
	"catty/internal/classfile"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// bootstrapRouteC closes the JDK-surface gaps surfaced by the first
// third-party library (minimal-json): Serializable marker, String reader/
// writer, Double wrapper, List/Iterator protocol, Arrays/Collections utils.

func bootstrapRouteC(k *Kernel) {
	// Serializable: pure marker.

	// Writer abstract + StringWriter (payload-backed *strings.Builder).
	mustDefine(k, &ClassDef{
		Name:  "java/io/Writer",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccAbstract,
		Methods: []MethodDef{
			{Name: "write", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natWriterWriteString},
			{Name: "write", Desc: "(I)V", Flags: classfile.AccPublic, Native: natWriterWriteCharDefault},
			{Name: "append", Desc: "(Ljava/lang/CharSequence;)Ljava/io/Writer;", Flags: classfile.AccPublic, Native: natWriterAppendCS},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic | classfile.AccAbstract},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/StringWriter",
		Super: "java/io/Writer",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natStringWriterInit},
			{Name: "write", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natStringWriterWrite},
			{Name: "write", Desc: "([CII)V", Flags: classfile.AccPublic, Native: natStringWriterWriteCII},
			{Name: "write", Desc: "(I)V", Flags: classfile.AccPublic, Native: natStringWriterWriteChar},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringWriterToString},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natObjectInit},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/StringReader",
		Super: "java/io/Reader",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natStringReaderInit},
			{Name: "read", Desc: "([CII)I", Flags: classfile.AccPublic, Native: natStringReaderRead},
			{Name: "read", Desc: "()I", Flags: classfile.AccPublic, Native: natStringReaderRead1},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natFileClose},
		},
	})

	// Float wrapper (cat1 field).
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Float",
		Super: "java/lang/Number",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Fields: []FieldDef{
			{Name: "value", Desc: "F", Flags: classfile.AccPrivate | classfile.AccFinal},
			{Name: "TYPE", Desc: "Ljava/lang/Class;", Flags: 0x0019},
		},
		StaticInit: func(k *Kernel, c *Class) error {
			prim, err := k.primitiveClass("F", "float")
			if err != nil {
				return err
			}
			if f := c.fieldsByKey[memberKey("TYPE", "Ljava/lang/Class;")]; f != nil {
				c.Statics[f.StaticSlot] = prim
			}
			return nil
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(F)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Fields[0] = args[0]
					return nil, nil
				}},
			{Name: "parseFloat", Desc: "(Ljava/lang/String;)F", Flags: classfile.AccPublic | classfile.AccStatic, Native: natFloatParseString},
			{Name: "valueOf", Desc: "(F)Ljava/lang/Float;", Flags: classfile.AccPublic | classfile.AccStatic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					in, _ := ctx.K.NewInstance(mustLookup(ctx.K, "java/lang/Float"))
					in.Fields[0] = args[0]
					return in, nil
				}},
		},
	})
	// Double wrapper (cat2 field).
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Double",
		Super: "java/lang/Number",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Fields: []FieldDef{
			{Name: "value", Desc: "D", Flags: classfile.AccPrivate | classfile.AccFinal},
			{Name: "TYPE", Desc: "Ljava/lang/Class;", Flags: 0x0019},
		},
		StaticInit: func(k *Kernel, c *Class) error {
			prim, err := k.primitiveClass("D", "double")
			if err != nil {
				return err
			}
			if f := c.fieldsByKey[memberKey("TYPE", "Ljava/lang/Class;")]; f != nil {
				c.Statics[f.StaticSlot] = prim
			}
			return nil
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(D)V", Flags: classfile.AccPublic, Native: natDoubleInit},
			{Name: "doubleValue", Desc: "()D", Flags: classfile.AccPublic, Native: natWrapperValue("D")},
			{Name: "isNaN", Desc: "(D)Z", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleIsNaN},
			{Name: "isInfinite", Desc: "(D)Z", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleIsInfinite},
			{Name: "parseDouble", Desc: "(Ljava/lang/String;)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleParse},
			{Name: "toString", Desc: "(D)Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleToString},
			{Name: "valueOf", Desc: "(D)Ljava/lang/Double;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natWrapperValueOf("java/lang/Double")},
			{Name: "parseFloat", Desc: "(Ljava/lang/String;)F", Flags: classfile.AccPublic | classfile.AccStatic, Native: natFloatParseString},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natWrapperToString},
		},
	})

	// List interface + iterator protocol on ArrayList.
	mustDefine(k, &ClassDef{
		Name:   "java/util/ArrayList$Itr",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/util/Iterator"},
		Flags:  classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "hasNext", Desc: "()Z", Flags: classfile.AccPublic, Native: natItrHasNext},
			{Name: "next", Desc: "()Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natItrNext},
		},
	})

	// Arrays.asList / Collections.unmodifiableList → pragmatic aliases.
	mustDefine(k, &ClassDef{
		Name:  "java/util/Arrays",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "asList", Desc: "([Ljava/lang/Object;)Ljava/util/List;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natArraysAsList},
			{Name: "fill", Desc: "([CC)V", Flags: classfile.AccPublic | classfile.AccStatic, Native: natArraysFillChar},
			{Name: "fill", Desc: "([II)V", Flags: classfile.AccPublic | classfile.AccStatic, Native: natArraysFillInt},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Math",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Methods: []MethodDef{
			{Name: "max", Desc: "(II)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMaxInt},
			{Name: "min", Desc: "(II)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMinInt},
			{Name: "max", Desc: "(JJ)J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMaxLong},
			{Name: "min", Desc: "(JJ)J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMinLong},
			{Name: "max", Desc: "(DD)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMaxDouble},
			{Name: "min", Desc: "(DD)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMinDouble},
			{Name: "abs", Desc: "(I)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathAbsInt},
			{Name: "abs", Desc: "(J)J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathAbsLong},
			{Name: "abs", Desc: "(D)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathAbsDouble},
			{Name: "floor", Desc: "(D)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathFloor},
			{Name: "ceil", Desc: "(D)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathCeil},
			{Name: "sqrt", Desc: "(D)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathSqrt},
			{Name: "rint", Desc: "(D)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathRint},
			{Name: "floor", Desc: "(F)F", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathFloorF},
		},
	})
	// java.util.concurrent.ConcurrentHashMap backed by Go's sync.Map —
	// sufficient for gson's reflective metadata caches.
	mustDefine(k, &ClassDef{
		Name:  "java/util/concurrent/ConcurrentHashMap",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload = &sync.Map{}
					return nil, nil
				}},
			{Name: "get", Desc: "(Ljava/lang/Object;)Ljava/lang/Object;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					if v, ok := recv.(*Instance).Payload.(*sync.Map).Load(hashKey(args[0])); ok {
						return v, nil
					}
					return nil, nil
				}},
			{Name: "put", Desc: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					prev, loaded := recv.(*Instance).Payload.(*sync.Map).Swap(hashKey(args[0]), args[1])
					if loaded {
						return prev, nil
					}
					return nil, nil
				}},
			{Name: "containsKey", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					_, ok := recv.(*Instance).Payload.(*sync.Map).Load(hashKey(args[0]))
					return boolV(ok), nil
				}},
			{Name: "size", Desc: "()I", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					n := 0
					recv.(*Instance).Payload.(*sync.Map).Range(func(any, any) bool { n++; return true })
					return int32(n), nil
				}},
		},
	})
	// java.util.Locale + java.util.Calendar — minimal surfaces (gson ctor
	// path: TypeAdapters.CALENDAR_FACTORY references Calendar at Gson() time).
	mustDefine(k, &ClassDef{
		Name:  "java/util/Locale",
		Super: "java/lang/Object",
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/Calendar",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getInstance", Desc: "(Ljava/util/Locale;)Ljava/util/Calendar;",
				Flags: classfile.AccPublic | classfile.AccStatic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					c, _ := ctx.K.ClassByName("java/util/GregorianCalendar")
					if c == nil {
						return nil, ctx.Throw("java/lang/RuntimeException", "GregorianCalendar missing")
					}
					return ctx.K.NewInstance(c)
				}},
			{Name: "getTimeInMillis", Desc: "()J", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return int64(0), nil // stub
				}},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/GregorianCalendar",
		Super: "java/util/Calendar",
	})
	// java.util.Currency — minimal surface (gson DefaultDateTypeAdapter? no —
	// TypeAdapters.CURRENCY factory references it at Gson() construction).
	mustDefine(k, &ClassDef{
		Name:  "java/util/Currency",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getInstance", Desc: "(Ljava/util/Locale;)Ljava/util/Currency;",
				Flags: classfile.AccPublic | classfile.AccStatic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					c, _ := ctx.K.ClassByName("java/util/Currency")
					return ctx.K.NewInstance(c)
				}},
			{Name: "getCurrencyCode", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo("XXX"), nil
				}},
		},
	})
	// java.util.Date — minimal surface (gson DefaultDateTypeAdapter probes).
	mustDefine(k, &ClassDef{
		Name:  "java/util/Date",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload = time.Now().UnixMilli()
					return nil, nil
				}},
			{Name: "getTime", Desc: "()J", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Payload.(int64), nil
				}},
			{Name: "setTime", Desc: "(J)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload = args[0]
					return nil, nil
				}},
		},
	})
	// java.util.UUID — randomUUID + toString/equals (gson internal caches)
	mustDefine(k, &ClassDef{
		Name:  "java/util/UUID",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(JJ)V", Flags: 0x0000,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload = [2]int64{args[0].(int64), args[1].(int64)}
					return nil, nil
				}},
			{Name: "randomUUID", Desc: "()Ljava/util/UUID;", Flags: classfile.AccPublic | classfile.AccStatic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					cls, _ := ctx.K.ClassByName("java/util/UUID")
					in, ierr := ctx.K.NewInstance(cls)
					if ierr != nil {
						return nil, ierr
					}
					hi := int64(time.Now().UnixNano())
					lo := int64(time.Now().UnixNano() << 13)
					in.Payload = [2]int64{hi, lo}
					return in, nil
				}},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					p := recv.(*Instance).Payload.([2]int64)
					return ctx.NewStringGo(fmt.Sprintf("%x-%x", p[0], p[1])), nil
				}},
		},
	})
	// java.net.InetAddress — identity-only surface (getHostName/getHostAddress
	// return the stored string; getByName resolves via net.LookupHost).
	mustDefine(k, &ClassDef{
		Name:  "java/net/InetAddress",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getHostName", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(inetAddressString(recv)), nil
				}},
			{Name: "getHostAddress", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(inetAddressString(recv)), nil
				}},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(inetAddressString(recv)), nil
				}},
		},
	})
	// java.net.URI — parse + getters only (create/getScheme/getHost/getPath).
	mustDefine(k, &ClassDef{
		Name:  "java/net/URI",
		Super: "java/lang/Object",
		Fields: []FieldDef{
			{Name: "scheme", Desc: "Ljava/lang/String;", Flags: 0x0002},
			{Name: "host", Desc: "Ljava/lang/String;", Flags: 0x0002},
			{Name: "path", Desc: "Ljava/lang/String;", Flags: 0x0002},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natURICtor},
			{Name: "getScheme", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(uriField(recv.(*Instance), 0)), nil
				}},
			{Name: "getHost", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(uriField(recv.(*Instance), 1)), nil
				}},
			{Name: "getPath", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(uriField(recv.(*Instance), 2)), nil
				}},
		},
	})
	// java.net.URL — parse + getters only; openStream throws.
	mustDefine(k, &ClassDef{
		Name:  "java/net/URL",
		Super: "java/lang/Object",
		Fields: []FieldDef{
			{Name: "protocol", Desc: "Ljava/lang/String;", Flags: 0x0002},
			{Name: "host", Desc: "Ljava/lang/String;", Flags: 0x0002},
			{Name: "path", Desc: "Ljava/lang/String;", Flags: 0x0002},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic,
				Native: natURLInit},
			{Name: "getProtocol", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(recv.(*Instance).Fields[0].(string)), nil
				}},
			{Name: "getHost", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(recv.(*Instance).Fields[1].(string)), nil
				}},
			{Name: "getPath", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return ctx.NewStringGo(recv.(*Instance).Fields[2].(string)), nil
				}},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					in := recv.(*Instance)
					return ctx.NewStringGo(
						in.Fields[0].(string) + "://" + in.Fields[1].(string) + in.Fields[2].(string)), nil
				}},
			{Name: "getByName", Desc: "(Ljava/lang/String;)Ljava/net/InetAddress;",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natInetAddressGetByName},
		},
	})
	// java.util.BitSet — minimal surface for gson's Excluder internals.
	// java.util.concurrent.atomic.AtomicInteger — minimal surface
	mustDefine(k, &ClassDef{
		Name:  "java/util/concurrent/atomic/AtomicInteger",
		Super: "java/lang/Object",
		Fields: []FieldDef{
			{Name: "value", Desc: "I", Flags: 0x0002},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(I)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Fields[0] = args[0]
					return nil, nil
				}},
			{Name: "incrementAndGet", Desc: "()I", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					n := recv.(*Instance).Fields[0].(int32) + 1
					recv.(*Instance).Fields[0] = n
					return n, nil
				}},
			{Name: "getAndIncrement", Desc: "()I", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					n := recv.(*Instance).Fields[0].(int32)
					recv.(*Instance).Fields[0] = n + 1
					return n, nil
				}},
			{Name: "get", Desc: "()I", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Fields[0], nil
				}},
			{Name: "set", Desc: "(I)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Fields[0] = args[0]
					return nil, nil
				}},
		},
	})
	// java.util.concurrent.atomic.AtomicBoolean — minimal surface
	// AtomicIntegerArray / AtomicLongArray / AtomicLong / AtomicLongArray minimal
	mustDefine(k, &ClassDef{
		Name:  "java/util/concurrent/atomic/AtomicIntegerArray",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(I)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					arr, _ := ctx.K.NewArray("I", int(argI(args, 0)))
					recv.(*Instance).Payload = arr
					return nil, nil
				}},
			{Name: "get", Desc: "(I)I", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Payload.(*ArrayObj).Elems[argI(args, 0)], nil
				}},
			{Name: "set", Desc: "(II)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload.(*ArrayObj).Elems[argI(args, 0)] = argI(args, 1)
					return nil, nil
				}},
			{Name: "incrementAndGet", Desc: "(I)I", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					arr := recv.(*Instance).Payload.(*ArrayObj)
					i := argI(args, 0)
					n := arr.Elems[i].(int32) + 1
					arr.Elems[i] = n
					return n, nil
				}},
		},
	})
	// java.math.BigDecimal / BigInteger — minimal surfaces (gson numeric paths).
	mustDefine(k, &ClassDef{
		Name:  "java/math/BigDecimal",
		Super: "java/lang/Number",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(D)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload = args[0]
					return nil, nil
				}},
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					if js, ok := args[0].(*JString); ok && js != nil {
						if f, perr := strconv.ParseFloat(strings.TrimSpace(js.Go()), 64); perr == nil {
							recv.(*Instance).Payload = f
							return nil, nil
						}
					}
					return nil, ctx.Throw("java/lang/NumberFormatException", "BigDecimal")
				}},
			{Name: "doubleValue", Desc: "()D", Flags: classfile.AccPublic,
				Native: natBigDecimalDoubleValue},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: natBigDecimalToString},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/math/BigInteger",
		Super: "java/lang/Number",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					js, _ := args[0].(*JString)
					v, perr := strconv.ParseInt(strings.TrimSpace(js.Go()), 10, 64)
					if perr != nil {
						return nil, ctx.Throw("java/lang/NumberFormatException", js.Go())
					}
					recv.(*Instance).Payload = v
					return nil, nil
				}},
			{Name: "doubleValue", Desc: "()D", Flags: classfile.AccPublic,
				Native: natBigIntegerDoubleValue},
			{Name: "longValue", Desc: "()J", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Payload.(int64), nil
				}},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: natBigIntegerToString},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/concurrent/atomic/AtomicLongArray",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(I)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					arr, _ := ctx.K.NewArray("J", int(argI(args, 0)))
					recv.(*Instance).Payload = arr
					return nil, nil
				}},
			{Name: "get", Desc: "(I)J", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Payload.(*ArrayObj).Elems[argI(args, 0)], nil
				}},
			{Name: "set", Desc: "(IJ)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload.(*ArrayObj).Elems[argI(args, 0)] = argL(args, 1)
					return nil, nil
				}},
			{Name: "incrementAndGet", Desc: "(I)J", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					arr := recv.(*Instance).Payload.(*ArrayObj)
					i := argI(args, 0)
					n := arr.Elems[i].(int64) + 1
					arr.Elems[i] = n
					return n, nil
				}},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/concurrent/atomic/AtomicLong",
		Super: "java/lang/Object",
		Fields: []FieldDef{
			{Name: "value", Desc: "J", Flags: 0x0002},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(J)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Fields[0] = args[0]
					return nil, nil
				}},
			{Name: "get", Desc: "()J", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Fields[0], nil
				}},
			{Name: "incrementAndGet", Desc: "()J", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					n := recv.(*Instance).Fields[0].(int64) + 1
					recv.(*Instance).Fields[0] = n
					return n, nil
				}},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/concurrent/atomic/AtomicBoolean",
		Super: "java/lang/Object",
		Fields: []FieldDef{
			{Name: "value", Desc: "Z", Flags: 0x0002},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Z)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Fields[0] = args[0]
					return nil, nil
				}},
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Fields[0] = int32(0)
					return nil, nil
				}},
			{Name: "get", Desc: "()Z", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Fields[0], nil
				}},
			{Name: "set", Desc: "(Z)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Fields[0] = args[0]
					return nil, nil
				}},
			{Name: "compareAndSet", Desc: "(ZZ)Z", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					in := recv.(*Instance)
					if in.Fields[0].(int32) == args[0] {
						in.Fields[0] = args[1]
						return boolV(true), nil
					}
					return boolV(false), nil
				}},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/BitSet",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload = make([]bool, 0, 8)
					return nil, nil
				}},
			{Name: "set", Desc: "(I)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					bs := recv.(*Instance).Payload.([]bool)
					i := int(argI(args, 0))
					for len(bs) <= i {
						bs = append(bs, false)
					}
					bs[i] = true
					recv.(*Instance).Payload = bs
					return nil, nil
				}},
			{Name: "get", Desc: "(I)Z", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					bs := recv.(*Instance).Payload.([]bool)
					i := int(argI(args, 0))
					return boolV(i < len(bs) && bs[i]), nil
				}},
			{Name: "clear", Desc: "(I)V", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					bs := recv.(*Instance).Payload.([]bool)
					if i := int(argI(args, 0)); i < len(bs) {
						bs[i] = false
					}
					return nil, nil
				}},
			{Name: "isEmpty", Desc: "()Z", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					for _, b := range recv.(*Instance).Payload.([]bool) {
						if b {
							return boolV(false), nil
						}
					}
					return boolV(true), nil
				}},
			{Name: "cardinality", Desc: "()I", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					n := int32(0)
					for _, b := range recv.(*Instance).Payload.([]bool) {
						if b {
							n++
						}
					}
					return n, nil
				}},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/Collections",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "unmodifiableList", Desc: "(Ljava/util/List;)Ljava/util/List;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natCollectionsUnmodList},
			// Deviation (registered): empty* return mutable empties —
			// immutability not enforced; sufficient for gson's defaults.
			{Name: "emptyList", Desc: "()Ljava/util/List;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natCollectionsEmptyList},
			{Name: "emptySet", Desc: "()Ljava/util/Set;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natCollectionsEmptySet},
			{Name: "emptyMap", Desc: "()Ljava/util/Map;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natCollectionsEmptyMap},
			{Name: "singletonList", Desc: "(Ljava/lang/Object;)Ljava/util/List;", Flags: classfile.AccPublic | classfile.AccStatic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					list := newBareInstance(ctx.K, "java/util/ArrayList")
					list.Payload = &alBuf{data: []Value{args[0]}}
					return list, nil
				}},
		},
	})
}

func natStringWriterInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Payload = &strings.Builder{}
	return nil, nil
}

func natStringWriterWrite(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := recv.(*Instance).Payload.(*strings.Builder)
	b.WriteString(javaStr(args[0]))
	return nil, nil
}

func natStringWriterWriteCII(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := recv.(*Instance).Payload.(*strings.Builder)
	arr := args[0].(*ArrayObj)
	off, length := int(argI(args, 1)), int(argI(args, 2))
	for i := off; i < off+length && i < len(arr.Elems); i++ {
		if c, ok := arr.Elems[i].(int32); ok {
			b.WriteRune(rune(c))
		}
	}
	return nil, nil
}

func natStringWriterWriteChar(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := recv.(*Instance).Payload.(*strings.Builder)
	if c, ok := args[0].(int32); ok {
		b.WriteRune(rune(c))
	}
	return nil, nil
}

func natStringWriterToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.InternGo(recv.(*Instance).Payload.(*strings.Builder).String()), nil
}

type strRdr struct {
	cur  int
	s    string
	done bool
}

func natStringReaderInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	recv.(*Instance).Payload = &strRdr{s: js.String()}
	return nil, nil
}

func (sr *strRdr) readOne() int {
	if sr.done {
		return -1
	}
	r := []rune(sr.s)
	if sr.cur >= len(r) {
		sr.done = true
		return -1
	}
	ch := r[sr.cur]
	sr.cur++
	return int(ch)
}

func natStringReaderRead1(ctx *CallContext, recv Value, args []Value) (Value, error) {
	sr := recv.(*Instance).Payload.(*strRdr)
	return int32(sr.readOne()), nil
}

func natStringReaderRead(ctx *CallContext, recv Value, args []Value) (Value, error) {
	sr := recv.(*Instance).Payload.(*strRdr)
	arr := args[0].(*ArrayObj)
	off := argI(args, 1)
	length := argI(args, 2)
	n := 0
	if os.Getenv("CATTY_SRDBG") != "" {
		js, _ := AsJString(args[0])
		_ = js
	}
	for i := 0; i < int(length); i++ {
		c := sr.readOne()
		if c == -1 {
			break
		}
		arr.Elems[off+int32(i)] = int32(c)
		n++
	}
	if n == 0 {
		return int32(-1), nil // JDK Reader contract: EOF is -1, never 0
	}
	if os.Getenv("CATTY_SRDBG") != "" {
		first := ""
		for i := 0; i < n && i < 12; i++ {
			first += string(rune(arr.Elems[int(off)+i].(int32)))
		}
		println("[SR] read off=", int(off), "len=", int(length), "->n=", n, "first=12:", first)
	}
	return int32(n), nil
}

func natDoubleInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Fields[0] = args[0]
	return nil, nil
}

func natDoubleIsNaN(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := args[0].(float64)
	return boolV(v != v), nil
}

func natDoubleIsInfinite(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := args[0].(float64)
	return boolV(v > 1.7976931348623157e308 || v < -1.7976931348623157e308), nil
}

func natDoubleParse(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	v, err := strconv.ParseFloat(strings.TrimSpace(js.String()), 64)
	if err != nil {
		return nil, ctx.Throw("java/lang/NumberFormatException", js.String())
	}
	return v, nil
}

func natDoubleToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.InternGo(fmt.Sprintf("%v", args[0].(float64))), nil
}

// ---- ArrayList iterator ----

type itrBuf struct {
	list *alBuf
	idx  int
}

func natArrayListIterator(ctx *CallContext, recv Value, args []Value) (Value, error) {
	itrCls, ok := ctx.K.ClassByName("java/util/ArrayList$Itr")
	if !ok {
		return nil, fmt.Errorf("ArrayList$Itr missing")
	}
	obj, err := ctx.K.NewInstance(itrCls)
	if err != nil {
		return nil, err
	}
	obj.Payload = &itrBuf{list: alOf(recv)}
	return obj, nil
}

func natItrHasNext(ctx *CallContext, recv Value, args []Value) (Value, error) {
	it := recv.(*Instance).Payload.(*itrBuf)
	return boolV(it.idx < len(it.list.data)), nil
}

func natItrNext(ctx *CallContext, recv Value, args []Value) (Value, error) {
	it := recv.(*Instance).Payload.(*itrBuf)
	if it.idx >= len(it.list.data) {
		return nil, ctx.Throw("java/util/NoSuchElementException", "iterator exhausted")
	}
	v := it.list.data[it.idx]
	it.idx++
	return v, nil
}

// ---- arrays-as-list / unmodifiable passthrough ----

func natArraysAsList(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arrCls, ok := ctx.K.ClassByName("java/util/ArrayList")
	if !ok {
		return nil, fmt.Errorf("ArrayList missing")
	}
	obj, err := ctx.K.NewInstance(arrCls)
	if err != nil {
		return nil, err
	}
	buf := &alBuf{}
	if len(args) > 0 {
		if a, ok := args[0].(*ArrayObj); ok {
			buf.data = append(buf.data, a.Elems...)
		}
	}
	obj.Payload = buf
	return obj, nil
}

func natCollectionsUnmodList(ctx *CallContext, recv Value, args []Value) (Value, error) {
	// Deviation (registered): returns the list itself — mutability not
	// enforced; sufficient for read-only third-party consumers.
	return args[0], nil
}

// natWriterWriteString mirrors java.io.Writer.write(String)'s concrete
// default: convert to chars and delegate to the subclass write([CII).
func natWriterWriteString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	rs := []rune(js.String())
	arr := &ArrayObj{CompDesc: "C", Elems: make([]Value, len(rs))}
	for i, r := range rs {
		arr.Elems[i] = int32(r)
	}
	wm, err := ctx.K.ResolveMethod(recv.(*Instance).Class, "write", "([CII)V")
	if err != nil {
		return nil, err
	}
	_, terr := ctx.K.InvokeAs(ctx.Owner, wm, recv, []Value{arr, int32(0), int32(len(rs))})
	if terr != nil {
		return nil, terr
	}
	return nil, nil
}

func natWriterWriteCharDefault(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arr := &ArrayObj{CompDesc: "C", Elems: []Value{args[0]}}
	wm, err := ctx.K.ResolveMethod(recv.(*Instance).Class, "write", "([CII)V")
	if err != nil {
		return nil, err
	}
	_, terr := ctx.K.InvokeAs(ctx.Owner, wm, recv, []Value{arr, int32(0), int32(1)})
	return nil, terr
}

func natWriterAppendCS(ctx *CallContext, recv Value, args []Value) (Value, error) {
	wm, err := ctx.K.ResolveMethod(recv.(*Instance).Class, "write", "(Ljava/lang/String;)V")
	if err != nil {
		return nil, err
	}
	if _, terr := ctx.K.InvokeAs(ctx.Owner, wm, recv, args); terr != nil {
		return nil, terr
	}
	return recv, nil
}

// ---- java.lang.Math double surface (P-0009 U4 fixture needs) ----

func natMathFloor(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return math.Floor(argD(args, 0)), nil
}

func natMathCeil(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return math.Ceil(argD(args, 0)), nil
}

func natMathSqrt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	d := argD(args, 0)
	if d < 0 || d != d {
		return math.NaN(), nil
	}
	return math.Sqrt(d), nil
}

func natMathRint(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return math.RoundToEven(argD(args, 0)), nil
}

func natMathFloorF(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return float32(math.Floor(float64(argF(args, 0)))), nil
}

func natURLInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	spec := ""
	if js, ok := args[0].(*JString); ok && js != nil {
		spec = js.Go()
	}
	in := recv.(*Instance)
	rest := spec
	proto := ""
	if i := strings.Index(rest, "://"); i >= 0 {
		proto = rest[:i]
		rest = rest[i+3:]
	} else {
		return nil, ctx.Throw("java/net/MalformedURLException", "no protocol: "+spec)
	}
	host, path := "", "/"
	if i := strings.Index(rest, "/"); i >= 0 {
		host, path = rest[:i], rest[i:]
	} else {
		host = rest
	}
	set := func(idx int, s string) {
		f := in.Class.FindField([]string{"protocol", "host", "path"}[idx],
			"Ljava/lang/String;")
		if f != nil {
			in.Fields[f.Slot] = ctx.K.MakeJStringFromGo(s)
		}
	}
	set(0, proto)
	set(1, host)
	set(2, path)
	return nil, nil
}

func uriField(in *Instance, which int) string {
	names := []string{"scheme", "host", "path"}
	f := in.Class.FindField(names[which], "Ljava/lang/String;")
	if f == nil {
		return ""
	}
	if js, ok := in.Fields[f.Slot].(*JString); ok && js != nil {
		return js.Go()
	}
	return ""
}

func natURICtor(ctx *CallContext, recv Value, args []Value) (Value, error) {
	spec := ""
	if js, ok := args[0].(*JString); ok && js != nil {
		spec = js.Go()
	}
	in := recv.(*Instance)
	scheme, host, path := "", "", ""
	rest := spec
	if i := strings.Index(rest, ":"); i >= 0 {
		scheme = rest[:i]
		rest = rest[i+1:]
	}
	if strings.HasPrefix(rest, "//") {
		rest = rest[2:]
		if i := strings.Index(rest, "/"); i >= 0 {
			host, path = rest[:i], rest[i:]
		} else {
			host = rest
		}
	} else {
		path = rest
	}
	set := func(n, v string) {
		f := in.Class.FindField(n, "Ljava/lang/String;")
		if f != nil {
			in.Fields[f.Slot] = ctx.K.MakeJStringFromGo(v)
		}
	}
	set("scheme", scheme)
	set("host", host)
	set("path", path)
	return nil, nil
}

func inetAddressString(recv Value) string {
	in, _ := recv.(*Instance)
	if in == nil || in.Payload == nil {
		return ""
	}
	s, _ := in.Payload.(string)
	return s
}

func natInetAddressGetByName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := args[0].(*JString)
	host := ""
	if js != nil {
		host = js.Go()
	}
	in, err := ctx.K.NewInstance(mustLookup(ctx.K, "java/net/InetAddress"))
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	in.Payload = host
	return in, nil
}

func natBigDecimalDoubleValue(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return recv.(*Instance).Payload.(float64), nil
}

func natBigDecimalToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.NewStringGo(fmtFloatStr(recv.(*Instance).Payload.(float64))), nil
}

func natBigIntegerDoubleValue(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return float64(recv.(*Instance).Payload.(int64)), nil
}

func natBigIntegerToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.NewStringGo(fmtIntStr(recv.(*Instance).Payload.(int64))), nil
}

func fmtFloatStr(f float64) string {
	return fmt.Sprintf("%v", f)
}

func fmtIntStr(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	u := uint64(v)
	if neg {
		u = uint64(-v)
	}
	var buf [24]byte
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
