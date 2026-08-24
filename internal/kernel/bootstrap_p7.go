package kernel

import (
	"catty/internal/classfile"
)

// bootstrapP7 extends the synthesized java.* surface with collections,
// wrapper types, and wider String methods needed by non-trivial programs.

func wrapperDef(name string, primDesc string, extraMethods []MethodDef) *ClassDef {
	valueMethod := primDesc + "Value"
	switch name {
	case "java/lang/Long":
		valueMethod = "longValue"
	case "java/lang/Boolean":
		valueMethod = "booleanValue"
	case "java/lang/Character":
		valueMethod = "charValue"
	}

	def := &ClassDef{
		Name:  name,
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Fields: []FieldDef{
			{Name: "value", Desc: primDesc, Flags: classfile.AccPrivate | classfile.AccFinal},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(" + primDesc + ")V", Flags: classfile.AccPublic, Native: natWrapperInit(primDesc)},
			{Name: valueMethod, Desc: "()" + primDesc, Flags: classfile.AccPublic, Native: natWrapperValue(primDesc)},
			{Name: "hashCode", Desc: "()I", Flags: classfile.AccPublic, Native: natWrapperHashCode},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natWrapperEquals},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natWrapperToString},
			{Name: "valueOf", Desc: "(" + primDesc + ")L" + name + ";", Flags: classfile.AccPublic | classfile.AccStatic, Native: natWrapperValueOf(name)},
		},
	}
	def.Methods = append(def.Methods, extraMethods...)
	return def
}

func bootstrapP7(k *Kernel) {
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Number",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccAbstract,
	})

	mustDefine(k, wrapperDef("java/lang/Long", "J",
		[]MethodDef{
			{Name: "parseLong", Desc: "(Ljava/lang/String;)J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natWrapperParse("java/lang/Long")},
		}))
	mustDefine(k, wrapperDef("java/lang/Boolean", "Z",
		[]MethodDef{
			{Name: "parseBoolean", Desc: "(Ljava/lang/String;)Z", Flags: classfile.AccPublic | classfile.AccStatic, Native: natBooleanParse},
		}))

	mustDefine(k, wrapperDef("java/lang/Short", "S",
		[]MethodDef{{Name: "parseShort", Desc: "(Ljava/lang/String;)S", Flags: classfile.AccPublic | classfile.AccStatic, Native: natWrapperParse("java/lang/Short")}}))
	mustDefine(k, wrapperDef("java/lang/Byte", "B",
		[]MethodDef{{Name: "parseByte", Desc: "(Ljava/lang/String;)B", Flags: classfile.AccPublic | classfile.AccStatic, Native: natWrapperParse("java/lang/Byte")}}))
	mustDefine(k, wrapperDef("java/lang/Character", "C", nil))

	// Note: Long's Super is already set to Object via wrapperDef.
	// Number hierarchy is a v2 concern.
}

func bootstrapCollectionsP7(k *Kernel) {
	// Dependency-ordered: interfaces before implementors.
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Iterable",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:   "java/util/Collection",
		Ifaces: []string{"java/lang/Iterable"},
		Flags:  classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/Iterator",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:   "java/util/List",
		Ifaces: []string{"java/util/Collection"},
		Flags:  classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:   "java/util/Set",
		Ifaces: []string{"java/util/Collection"},
		Flags:  classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/Map",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:   "java/util/HashMap",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/util/Map"},
		Flags:  classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natHashMapInit},
			{Name: "get", Desc: "(Ljava/lang/Object;)Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natHashMapGet},
			{Name: "put", Desc: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natHashMapPut},
			{Name: "remove", Desc: "(Ljava/lang/Object;)Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natHashMapRemove},
			{Name: "containsKey", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natHashMapContainsKey},
			{Name: "size", Desc: "()I", Flags: classfile.AccPublic, Native: natHashMapSize},
			{Name: "isEmpty", Desc: "()Z", Flags: classfile.AccPublic, Native: natHashMapIsEmpty},
			{Name: "clear", Desc: "()V", Flags: classfile.AccPublic, Native: natHashMapClear},
			{Name: "iterator", Desc: "()Ljava/util/Iterator;", Flags: classfile.AccPublic, Native: natArrayListIterator},
		},
	})
	mustDefine(k, &ClassDef{
		Name:   "java/util/HashSet",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/util/Set"},
		Flags:  classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natHashSetInit},
			{Name: "add", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natHashSetAdd},
			{Name: "contains", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natHashSetContains},
			{Name: "remove", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natHashSetRemove},
			{Name: "size", Desc: "()I", Flags: classfile.AccPublic, Native: natHashSetSize},
			{Name: "isEmpty", Desc: "()Z", Flags: classfile.AccPublic, Native: natHashSetIsEmpty},
		},
	})
	mustDefine(k, &ClassDef{
		Name:   "java/util/ArrayList",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/util/List"},
		Flags:  classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natArrayListInit},
			{Name: "add", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natArrayListAdd},
			{Name: "get", Desc: "(I)Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natArrayListGet},
			{Name: "set", Desc: "(ILjava/lang/Object;)Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natArrayListSet},
			{Name: "size", Desc: "()I", Flags: classfile.AccPublic, Native: natArrayListSize},
			{Name: "isEmpty", Desc: "()Z", Flags: classfile.AccPublic, Native: natArrayListIsEmpty},
			{Name: "contains", Desc: "(Ljava/lang/Object;)Z", Flags: classfile.AccPublic, Native: natArrayListContains},
			{Name: "iterator", Desc: "()Ljava/util/Iterator;", Flags: classfile.AccPublic, Native: natArrayListIterator},
		},
	})
}

func bootstrapStringP7(k *Kernel) {
	c := k.lookupClass("java/lang/String")
	if c == nil {
		return
	}
	extra := []MethodDef{
		{Name: "substring", Desc: "(I)Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringSubstring},
		{Name: "substring", Desc: "(II)Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringSubstring2},
		{Name: "trim", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringTrim},
		{Name: "toLowerCase", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringToLowerCase},
		{Name: "toUpperCase", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringToUpperCase},
		{Name: "contains", Desc: "(Ljava/lang/CharSequence;)Z", Flags: classfile.AccPublic, Native: natStringContains},
		{Name: "isEmpty", Desc: "()Z", Flags: classfile.AccPublic, Native: natStringIsEmpty},
		{Name: "split", Desc: "(Ljava/lang/String;)[Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringSplit},
		{Name: "replace", Desc: "(CC)Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringReplace},
		{Name: "startsWith", Desc: "(Ljava/lang/String;)Z", Flags: classfile.AccPublic, Native: natStringStartsWith},
		{Name: "endsWith", Desc: "(Ljava/lang/String;)Z", Flags: classfile.AccPublic, Native: natStringEndsWith},
	}
	for _, m := range extra {
		addMethodToClass(k, c, m)
	}

	mustDefine(k, &ClassDef{
		Name:  "java/lang/CharSequence",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
}

// addMethodToClass appends a method to an already-defined kernel Class.
func addMethodToClass(k *Kernel, c *Class, def MethodDef) {
	m := &Method{
		Holder: c,
		Name:   def.Name,
		Desc:   def.Desc,
		Flags:  def.Flags,
		Native: def.Native,
	}
	key := memberKey(m.Name, m.Desc)
	c.methodsByKey[key] = m
	c.Methods = append(c.Methods, m)
}

// bootstrapFileIO adds the Route-A file surface: File metadata and the
// line-oriented FileReader used by CLI-style fixtures.
func bootstrapFileIO(k *Kernel) {
	mustDefine(k, &ClassDef{
		Name:  "java/io/File",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natObjectInit},
			{Name: "exists", Desc: "()Z", Flags: classfile.AccPublic, Native: natFileExists},
			{Name: "length", Desc: "()J", Flags: classfile.AccPublic, Native: natFileLength},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/Reader",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccAbstract,
		Methods: []MethodDef{
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic | classfile.AccAbstract},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/InputStreamReader",
		Super: "java/io/Reader",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/io/InputStream;)V", Flags: classfile.AccPublic, Native: natObjectInit},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natFileClose},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/FileReader",
		Super: "java/io/InputStreamReader",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natFileReaderInit},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natFileClose},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/BufferedReader",
		Super: "java/io/Reader",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/io/Reader;)V", Flags: classfile.AccPublic, Native: natBufferedReaderInit},
			{Name: "readLine", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natBufferedReadLine},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natFileClose},
		},
	})
}
