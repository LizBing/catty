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
			{Name: "wait", Desc: "(J)V", Flags: classfile.AccPublic | classfile.AccFinal, Native: natObjectWaitMillis},
			{Name: "notify", Desc: "()V", Flags: classfile.AccPublic | classfile.AccFinal, Native: natObjectNotify},
			{Name: "notifyAll", Desc: "()V", Flags: classfile.AccPublic | classfile.AccFinal, Native: natObjectNotifyAll},
			{Name: "getClass", Desc: "()Ljava/lang/Class;", Flags: classfile.AccPublic | classfile.AccFinal, Native: natObjectGetClass},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/Serializable",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/CharSequence",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Comparable",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/String",
		Super: "java/lang/Object",
		Ifaces: []string{"java/io/Serializable", "java/lang/CharSequence",
			"java/lang/Comparable"},
		Flags: classfile.AccPublic | classfile.AccFinal,
		Methods: []MethodDef{
			{Name: "length", Desc: "()I", Flags: classfile.AccPublic, Native: natStringLength},
			{Name: "format", Desc: "(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccStatic | classfile.AccVarargs, Native: natStringFormat},
			// clone() — arrays and Strings only (gson internals clone
			// buffers); Object.clone semantics deferred.
			{Name: "clone", Desc: "()Ljava/lang/Object;", Flags: 0x0001, Native: natStringClone},
			{Name: "charAt", Desc: "(I)C", Flags: classfile.AccPublic, Native: natStringCharAt},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natStringEquals},
			{Name: "hashCode", Desc: "()I", Flags: classfile.AccPublic, Native: natStringHashCode},
			{Name: "indexOf", Desc: "(Ljava/lang/String;)I", Flags: classfile.AccPublic, Native: natStringIndexOf},
			{Name: "getBytes", Desc: "()[B", Flags: classfile.AccPublic, Native: natStringGetBytes},
			{Name: "<init>", Desc: "([BII)V", Flags: classfile.AccPublic, Native: natStringInitBytesRange},
			{Name: "getChars", Desc: "(II[CI)V", Flags: classfile.AccPublic, Native: natStringGetChars},
			{Name: "<init>", Desc: "([C)V", Flags: classfile.AccPublic, Native: natStringInitChars},
			{Name: "<init>", Desc: "([CII)V", Flags: classfile.AccPublic, Native: natStringInitCharsRange},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringToString},
		},
	})

	mustDefine(k, &ClassDef{
		Name:  "java/lang/StringBuilder",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natStringBuilderInit},
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natSBInitString},
			appendSB("Ljava/lang/Object;", natSBAppendObject),
			appendSB("Ljava/lang/String;", natSBAppendString),
			appendSB("I", natSBAppendInt),
			appendSB("J", natSBAppendLong),
			appendSB("C", natSBAppendChar),
			appendSB("Z", natSBAppendBool),
			{Name: "length", Desc: "()I", Flags: classfile.AccPublic, Native: natSBLength},
			{Name: "append", Desc: "([CII)Ljava/lang/StringBuilder;", Flags: classfile.AccPublic, Native: natSBAppendChars},
			{Name: "append", Desc: "(D)Ljava/lang/StringBuilder;", Flags: classfile.AccPublic, Native: natSBAppendDouble},
			{Name: "setLength", Desc: "(I)V", Flags: classfile.AccPublic, Native: natSBSetLength},
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
		{"java/io/IOException", "java/lang/Exception"},
		{"java/io/InterruptedIOException", "java/io/IOException"},
		{"java/net/SocketException", "java/lang/Exception"},
		{"java/net/BindException", "java/net/SocketException"},
		{"java/lang/IllegalThreadStateException", "java/lang/IllegalArgumentException"},
		{"java/lang/InterruptedException", "java/lang/Exception"},
		{"java/lang/VirtualMachineError", "java/lang/Error"},
		{"java/lang/StackOverflowError", "java/lang/VirtualMachineError"},
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

	// Reflection classes first: wrapper TYPE statics and any getClass/ldc
	// usage depend on java/lang/Class existing before other bootstrap pieces.
	mustDefine(k, &ClassDef{
		Name:  "java/lang/ThreadLocal",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: 0x0001, Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
				recv.(*Instance).Payload = &tlSlot{}
				return nil, nil
			}},
			{Name: "get", Desc: "()Ljava/lang/Object;", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Payload.(*tlSlot).v, nil
				}},
			{Name: "set", Desc: "(Ljava/lang/Object;)V", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload.(*tlSlot).v = args[0]
					return nil, nil
				}},
			{Name: "remove", Desc: "()V", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					recv.(*Instance).Payload.(*tlSlot).v = nil
					return nil, nil
				}},
		},
	})
	// java.lang.StringBuffer — StringBuilder semantics over the same buf.
	mustDefine(k, &ClassDef{
		Name:  "java/lang/StringBuffer",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natStringBuilderInit},
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natSBInitString},
			{Name: "append", Desc: "(Ljava/lang/String;)Ljava/lang/StringBuilder;", Flags: classfile.AccPublic, Native: natSBAppendString},
			{Name: "append", Desc: "(I)Ljava/lang/StringBuilder;", Flags: classfile.AccPublic, Native: natSBAppendInt},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natSBToString},
			{Name: "length", Desc: "()I", Flags: classfile.AccPublic, Native: natSBLength},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Enum",
		Super: "java/lang/Object",
		Fields: []FieldDef{
			{Name: "name", Desc: "Ljava/lang/String;", Flags: classfile.AccPrivate | classfile.AccFinal},
			{Name: "ordinal", Desc: "I", Flags: classfile.AccPrivate | classfile.AccFinal},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;I)V", Flags: 0x0004, // protected
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					in := recv.(*Instance)
					// Defensive: anonymous-subclass chains may reach here
					// before generic allocation filled inherited slots.
					for len(in.Fields) < 2 {
						in.Fields = append(in.Fields, nil)
					}
					in.Fields[0] = args[0]
					in.Fields[1] = args[1]
					return nil, nil
				}},
			{Name: "valueOf", Desc: "(Ljava/lang/Class;Ljava/lang/String;)Ljava/lang/Enum;",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natEnumValueOf},
			{Name: "name", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccFinal,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Fields[0], nil
				}},
			{Name: "ordinal", Desc: "()I", Flags: classfile.AccPublic | classfile.AccFinal,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Fields[1], nil
				}},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Fields[0], nil
				}},
		},
	})
	if err := registerReflection(k); err != nil {
		panic("kernel bootstrap reflection: " + err.Error())
	}
	mustDefine(k, &ClassDef{
		Name:  "java/io/Closeable",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/Flushable",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Cloneable",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/reflect/InvocationTargetException",
		Super: "java/lang/Exception",
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/InstantiationException",
		Super: "java/lang/Exception",
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Void",
		Super: "java/lang/Object",
		Fields: []FieldDef{{Name: "TYPE", Desc: "Ljava/lang/Class;", Flags: 0x0019}},
		StaticInit: func(k *Kernel, c *Class) error {
			prim, err := k.primitiveClass("V", "void")
			if err != nil { return err }
			if f := c.fieldsByKey[memberKey("TYPE", "Ljava/lang/Class;")]; f != nil {
				c.Statics[f.StaticSlot] = prim
			}
			return nil
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Number",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Integer",
		Super: "java/lang/Number",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Fields: []FieldDef{
			{Name: "value", Desc: "I", Flags: classfile.AccPrivate | classfile.AccFinal},
			{Name: "TYPE", Desc: "Ljava/lang/Class;", Flags: classfile.AccPublic | classfile.AccStatic | classfile.AccFinal},
		},
		StaticInit: func(k *Kernel, c *Class) error {
			prim, err := k.primitiveClass("I", "int")
			if err != nil {
				return err
			}
			if f := c.fieldsByKey[memberKey("TYPE", "Ljava/lang/Class;")]; f != nil {
				c.Statics[f.StaticSlot] = prim
			}
			return nil
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(I)V", Flags: classfile.AccPublic, Native: natIntegerInit},
			{Name: "intValue", Desc: "()I", Flags: classfile.AccPublic, Native: natIntegerIntValue},
			{Name: "hashCode", Desc: "()I", Flags: classfile.AccPublic, Native: natIntegerHashCode},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natIntegerEquals},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natIntegerToString},
			{Name: "valueOf", Desc: "(I)Ljava/lang/Integer;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticIntegerValueOf},
			{Name: "toString", Desc: "(I)Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticIntegerToString},
			{Name: "parseInt", Desc: "(Ljava/lang/String;)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticIntegerParseInt},
			{Name: "parseInt", Desc: "(Ljava/lang/String;I)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticIntegerParseIntRadix},
			{Name: "toString", Desc: "(II)Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticIntegerToStringRadix},
		},
	})

	mustDefine(k, &ClassDef{
		Name:  "java/lang/Runnable",
		Super: "",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})

	mustDefine(k, &ClassDef{
		Name:   "java/lang/Thread",
		Super:  "java/lang/Object",
		Flags:  classfile.AccPublic,
		Ifaces: []string{"java/lang/Runnable"},
		Fields: []FieldDef{
			{Name: "name", Desc: "Ljava/lang/String;", Flags: classfile.AccPrivate},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natThreadInitVoid},
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natThreadInitName},
			{Name: "run", Desc: "()V", Flags: classfile.AccPublic, Native: natThreadRunDefault},
			{Name: "start", Desc: "()V", Flags: classfile.AccPublic, Native: natThreadStart},
			{Name: "join", Desc: "()V", Flags: classfile.AccPublic, Native: natThreadJoinForever},
			{Name: "join", Desc: "(J)V", Flags: classfile.AccPublic, Native: natThreadJoinMillis},
			{Name: "isAlive", Desc: "()Z", Flags: classfile.AccPublic, Native: natThreadIsAlive},
			{Name: "setName", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natThreadSetName},
			{Name: "getName", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natThreadGetName},
			{Name: "interrupt", Desc: "()V", Flags: classfile.AccPublic, Native: natThreadInterrupt},
			{Name: "isInterrupted", Desc: "()Z", Flags: classfile.AccPublic, Native: natThreadIsInterrupted},
			{Name: "currentThread", Desc: "()Ljava/lang/Thread;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticCurrentThread},
			{Name: "sleep", Desc: "(J)V", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticSleep},
			{Name: "interrupted", Desc: "()Z", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticInterruptedFlag},
		},
	})

	bootstrapP7(k)
	bootstrapCollectionsP7(k)
	bootstrapStringP7(k)
	bootstrapFileIO(k)
	bootstrapRouteC(k)

	// ---- java.net / java.io streams (payload-backed, M2) ----
	mustDefine(k, &ClassDef{
		Name:  "java/io/InputStream",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natObjectInit},
			{Name: "read", Desc: "([B)I", Flags: classfile.AccPublic, Native: natStreamReadB},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natStreamClose},
		},
	})

	mustDefine(k, &ClassDef{
		Name:  "java/io/OutputStream",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natObjectInit},
			{Name: "write", Desc: "([B)V", Flags: classfile.AccPublic, Native: natStreamWriteB},
			{Name: "write", Desc: "([BII)V", Flags: classfile.AccPublic, Native: natStreamWriteBII},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natStreamClose},
		},
	})

	mustDefine(k, &ClassDef{
		Name:   "java/net/Socket",
		Super:  "java/lang/Object",
		Flags:  classfile.AccPublic,
		Fields: []FieldDef{},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natObjectInit},
			{Name: "getInputStream", Desc: "()Ljava/io/InputStream;", Flags: classfile.AccPublic, Native: natSocketGetInputStream},
			{Name: "getOutputStream", Desc: "()Ljava/io/OutputStream;", Flags: classfile.AccPublic, Native: natSocketGetOutputStream},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natSocketClose},
		},
	})

	mustDefine(k, &ClassDef{
		Name:   "java/net/ServerSocket",
		Super:  "java/lang/Object",
		Flags:  classfile.AccPublic,
		Fields: []FieldDef{},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(I)V", Flags: classfile.AccPublic, Native: natServerSocketInit},
			{Name: "accept", Desc: "()Ljava/net/Socket;", Flags: classfile.AccPublic, Native: natServerSocketAccept},
			{Name: "getLocalPort", Desc: "()I", Flags: classfile.AccPublic, Native: natServerSocketLocalPort},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natServerSocketClose},
		},
	})

	// ArrayList moved to bootstrapCollectionsP7 (Route C consolidation).

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
			{Name: "nanoTime", Desc: "()J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticNanoTime},
			{Name: "tickMillis", Desc: "()I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticTickMillis},
			{Name: "currentTimeMillis", Desc: "()J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natStaticCurrentTimeMillis},
			{Name: "arraycopy", Desc: "(Ljava/lang/Object;ILjava/lang/Object;II)V",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natSystemArraycopy},
			{Name: "currentTimeMillis", Desc: "()J",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natCurrentTimeMillis},
			{Name: "nanoTime", Desc: "()J",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natNanoTime},
			{Name: "identityHashCode", Desc: "(Ljava/lang/Object;)I",
				Flags: classfile.AccPublic | classfile.AccStatic, Native: natIdentityHashCode},
			{Name: "getProperty", Desc: "(Ljava/lang/String;)Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccStatic,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					// Synthesized JDK-8 view so third-party version probes
					// take the JDK8 code paths this runtime implements.
					switch js, _ := args[0].(*JString); js.Go() {
					case "java.version":
						return ctx.NewStringGo("1.8.0_402"), nil
					case "java.specification.version":
						return ctx.NewStringGo("1.8"), nil
					}
					return nil, nil
				}},
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
