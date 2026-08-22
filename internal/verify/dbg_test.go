package verify

import (
	"fmt"
	"os"
	"testing"

	"catty/internal/classfile"
)

func TestDumpEchoMain(t *testing.T) {
	data, _ := os.ReadFile("../../testdata/cp/HttpEcho.class")
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m := cf.FindMethod("main", "([Ljava/lang/String;)V")
	code := m.Code.Code
	fmt.Println("codelen:", len(code))
	fmt.Printf("tail: % x\n", code[70:])
	_, err = branchTargets(code)
	fmt.Println("scan err:", err)
}
