package txtedit
import (
	"fmt"
	"strings"
)

type Val struct {
	QuoteStyle string
	Text string
	TrailingSpaces string
}

type Comment struct {
	CommentStyle string
	Content string
}

type StmtContinue struct {
	Style string
}

type Stmt struct {
	Indent string
	// Value or Comment or StmtContinue
	Pieces []interface{}
	End string
}

type Sect struct {
	Begin *Stmt
	BeginPrefix, BeginSuffix string
	EndPrefix, EndSuffix string
	// Stmt or Sect
	Pieces []interface{}
	End *Stmt
}

type DocNode struct {
	Parent *DocNode
	// Stmt or Sect
	Obj interface{}
	Leaves []*DocNode
}

type AnalyserStyle struct {
	StmtContinue []string
	StmtEnd []string
	CommentBegin []string
	Quote []string
	BeginSectWithStmt, EndSectWithStmt bool

	SectBeginPrefix []string
	SectBeginSuffix []string
	SectEndPrefix []string
	SectEndSuffix []string
}

type Analyser struct {
	Style *AnalyserStyle
	Root *DocNode

	text string
	here int
	this *DocNode
}

func NewAnalyser(style *AnalyserStyle, input string) (ret *Analyser ){
	ret = &Analyser{Style: style, text: input}
	ret.this = &DocNode{Parent: nil, Obj: nil, Leaves:make([]*DocNode, 0, 8)}
	ret.Root = ret.this
	return
}

func DebugNode(node *DocNode, indent int) {
	prefix := strings.Repeat(" ", indent)
	if node == nil {
		fmt.Println(prefix + "(nil)")
		return
	}
	fmt.Print(prefix + "Node - ", node.Obj)
	if len(node.Leaves) > 0 {
		fmt.Println(" -->")
		for _, leaf := range node.Leaves {
			DebugNode(leaf, indent + 2)
		}
	} else {
		fmt.Println()
	}
}


func Print(root *DocNode) {
	fmt.Println(root.Obj)
	fmt.Println("-->")
	for _, leaf := range root.Leaves {
		Print(leaf)
	}
}