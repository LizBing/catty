package verify

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"catty/internal/classfile"
)

func TestDumpTargets(t *testing.T) {
	data, _ := os.ReadFile("../../testdata/cp/CollectionsDemo.class")
	cf, _ := classfile.Parse(data)
	m := cf.FindMethod("main", "([Ljava/lang/String;)V")
	tg, err := BranchTargets(m.Code.Code)
	if err != nil {
		t.Fatal(err)
	}
	keys := []int{}
	for k := range tg {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	fmt.Println("targets:", keys)
	fmt.Println("SM offsets:")
	for _, f := range m.Code.StackMaps {
		fmt.Println(" ", f.Offset)
	}
}
