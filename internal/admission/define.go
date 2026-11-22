package admission

import (
	"fmt"

	"catty/internal/classfile"
	"catty/rt"
)

// define converts one parsed, verified class into immutable runtime metadata.
// Every concrete plugin method executes in the interpreter.
func define(c *classfile.Class) (*rt.ClassDef, error) {
	def := &rt.ClassDef{
		Name:        c.Name,
		Super:       c.Super,
		Interfaces:  append([]string(nil), c.Interfaces...),
		Flags:       c.Flags,
		NestHost:    c.NestHost,
		NestMembers: append([]string(nil), c.NestMembers...),
	}
	for _, f := range c.Fields {
		fd := rt.FieldDef{Name: f.Name, Descriptor: f.Descriptor, Flags: f.Flags, Constant: f.Constant}
		if f.Flags&rt.AccStatic != 0 {
			def.StaticFields = append(def.StaticFields, fd)
		} else {
			def.InstanceFields = append(def.InstanceFields, fd)
		}
	}
	for _, m := range c.Methods {
		engine := "interp"
		if m.Flags&rt.AccAbstract != 0 {
			engine = "abstract"
		}
		md := &rt.MethodDef{
			Owner:      c.Name,
			Name:       m.Name,
			Descriptor: m.Descriptor,
			Flags:      m.Flags,
			Engine:     engine,
			Reason:     "plugin: interpreted",
			MaxStack:   m.MaxStack,
			MaxLocals:  m.MaxLocals,
			CodeSize:   m.CodeLength,
			Refs:       append([]rt.Ref(nil), m.Refs...),
			Insns:      append([]rt.Insn(nil), m.Insns...),
			Handlers:   append([]rt.Handler(nil), m.Handlers...),
			EntryIndex: -1,
		}
		if engine == "abstract" {
			md.CodeSize = 0
			md.Insns = nil
			md.Reason = "plugin: abstract declaration"
		} else if m.CodeLength == 0 || len(m.Insns) == 0 {
			return nil, fmt.Errorf("concrete plugin method %s.%s%s has no code", c.Name, m.Name, m.Descriptor)
		}
		def.Methods = append(def.Methods, md)
	}
	return def, nil
}
