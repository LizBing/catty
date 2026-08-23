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
	mustDefine(k, &ClassDef{
		Name:  "java/util/Map",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/Iterator",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/Collection",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:   "java/util/Set",
		Ifaces: []string{"java/util/Collection"},
		Flags:  classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
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
