package verify

import (
	"fmt"
	"os"
	"testing"

	"catty/internal/classfile"
)

func TestDumpSMOffsets(t *testing.T) {
	data, _ := os.ReadFile("../../testdata/cp/CollectionsDemo.class")
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m := cf.FindMethod("main", "([Ljava/lang/String;)V")
	fmt.Println("frames:", len(m.Code.StackMaps))
	for _, f := range m.Code.StackMaps {
		fmt.Println(" off:", f.Offset)
	}
}
