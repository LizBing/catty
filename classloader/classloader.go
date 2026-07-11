package classloader

import (
	"strings"

	"catty/classfile"
	"catty/classpath"
	"catty/native"
	"catty/rtda"
)

// ClassLoader loads, links, and caches classes. It implements rtda.Loader so the
// interpreter can resolve classes at run time without an import cycle.
//
// Loading order for a name:
//  1. cached?  -> return
//  2. array type ("[...")  -> rtda.NewArrayClass
//  3. a natively-implemented core class  -> native.NativeClass
//  4. otherwise  -> read+parse the .class file from the classpath
//
// Initialization (<clinit>) is NOT done here: it is triggered lazily by the
// interpreter at the JVMS §5.5 events (new / getstatic / putstatic /
// invokestatic), so the classloader never depends on the interpreter.
type ClassLoader struct {
	cp    *classpath.Classpath
	cache map[string]*rtda.Class
}

func New(cp *classpath.Classpath) *ClassLoader {
	return &ClassLoader{cp: cp, cache: make(map[string]*rtda.Class)}
}

// LoadClass returns the loaded class for the given internal name, loading it on
// first access. Names use internal slashes ("java/lang/Object").
func (cl *ClassLoader) LoadClass(name string) *rtda.Class {
	if c := cl.cache[name]; c != nil {
		return c
	}
	var c *rtda.Class
	switch {
	case strings.HasPrefix(name, "["):
		c = rtda.NewArrayClass(name, cl)
	default:
		if nc := native.NativeClass(cl, name); nc != nil {
			c = nc
		} else {
			c = cl.loadFileClass(name)
		}
	}
	cl.cache[name] = c
	return c
}

func (cl *ClassLoader) loadFileClass(name string) *rtda.Class {
	data, _, err := cl.cp.ReadClass(strings.ReplaceAll(name, ".", "/"))
	if err != nil {
		panic("catty: class not found: " + name + " (" + err.Error() + ")")
	}
	cf, err := classfile.Parse(data)
	if err != nil {
		panic("catty: parse error for " + name + ": " + err.Error())
	}
	return rtda.NewClass(cf, cl)
}
