package lexer

import (
	"fmt"
	"testing"
)

func TestManipulatingParentNode(t *testing.T) {
	base := &DocumentNode{}
	node1 := &DocumentNode{Obj: &Text{Text: "node1"}}

	// Wrong leaf anchor
	if base.InsertBefore(node1, node1) {
		t.Fatal("not false")
	}
	if base.InsertAfter(node1, node1) {
		t.Fatal("not false")
	}

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
	if len(base.Leaves) != 3 {
		t.Fatal(base.Leaves)
	}
	fmt.Println(DebugNode(base, 0))
}

func TestManipulatingNode(t *testing.T) {
	base := &DocumentNode{}
	node1 := &DocumentNode{Obj: &Text{Text: "node1"}}

	// Wrong leaf
	if base.InsertBeforeSelf(node1) {
		t.Fatal("not false")
	}
	if base.InsertAfterSelf(node1) {
		t.Fatal("not false")
	}

	if !base.InsertAfter(nil, node1) {
		t.Fatal("not done")
	}
	node0 := &DocumentNode{Obj: &Text{Text: "node0"}}
	if !node1.InsertBeforeSelf(node0) {
		t.Fatal("not done")
	}
	node2 := &DocumentNode{Obj: &Text{Text: "node2"}}
	if !node1.InsertAfterSelf(node2) {
		t.Fatal("not done")
	}
	if len(base.Leaves) != 3 {
		t.Fatal(base.Leaves)
	}
	if node0.Parent != base {
		t.Fatal("wrong base")
	}
	fmt.Println(DebugNode(base, 0))
	node1.DeleteSelf()
	node2.DeleteSelf()
	node0.DeleteSelf()
	if len(base.Leaves) != 0 {
		t.Fatal(base.Leaves)
	}
	fmt.Println(DebugNode(base, 0))
}
