package analyser

import (
	"fmt"
	"testing"
)

func TestNodeManipulation(t *testing.T) {
	base := &DocumentNode{}
	node1 := &DocumentNode{Obj: &Text{Text: "node1"}}
	if !base.InsertBefore(nil, node1) {
		t.Fatal("not done")
	}
	if i := base.FindLeafIndex(node1); i != 0 {
		t.Fatal(i)
	}
	node0 := &DocumentNode{Obj: &Text{Text: "node0"}}
	if !base.InsertBefore(node1, node0) {
		t.Fatal("not done")
	}
	if i := base.FindLeafIndex(node0); i != 0 {
		t.Fatal(i)
	}
	if i := base.FindLeafIndex(node1); i != 1 {
		t.Fatal(i)
	}
	node2 := &DocumentNode{Obj: &Text{Text: "node2"}}
	if !base.InsertAfter(node1, node2) {
		t.Fatal("not done")
	}
	if i := base.FindLeafIndex(node2); i != 2 {
		t.Fatal(i)
	}
	fmt.Println(DebugNode(base, 0))
}
