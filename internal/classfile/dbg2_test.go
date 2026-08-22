package classfile

import (
	"fmt"
	"os"
	"testing"
)

func TestDumpSMHex(t *testing.T) {
	data, _ := os.ReadFile("../../testdata/CollectionsDemo.class")
	cf, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m := cf.FindMethod("main", "([Ljava/lang/String;)V")
	fmt.Println("frames:", len(m.Code.StackMaps))
	for _, f := range m.Code.StackMaps {
		fmt.Printf("off=%d locals=%d stack=%d\n", f.Offset, len(f.Locals), len(f.Stack))
	}
}
