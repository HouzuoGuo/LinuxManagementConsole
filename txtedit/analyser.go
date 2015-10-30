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
	lastBranch, here int
	this *DocNode

	valCtx *Val
	commentCtx *Comment
	stmtCtx *Stmt
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

func (an *Analyser) NewLeaf() {
	if parent := an.this.Parent ; parent== nil {
		// this is root
		newLeaf := &DocNode{Parent:an.Root, Leaves: make([]*DocNode, 0, 8)}
		an.this.Leaves = append(an.this.Leaves, newLeaf)
		an.this = newLeaf
	} else {
		newLeaf := &DocNode{Parent: parent, Leaves:make([]*DocNode, 0, 8)}
		parent.Leaves = append(parent.Leaves, newLeaf)
		an.this = newLeaf
	}
}

func (an *Analyser) NewStmt() {
	an.NewLeaf()
	an.this.Obj = &Stmt{Indent:"", Pieces:make([]interface{}, 0, 8), End: ""}
}

func (an *Analyser) EndVal() {
	if an.valCtx  == nil {
		fmt.Println("Ending a val without starting one")
		return
	}
	if an.stmtCtx == nil {
		fmt.Println("Ending a val with a new stmt")
		an.NewStmt()
	}
	an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.valCtx)
}

func (an *Analyser) StoreContent() {
	if an.here - an.lastBranch > 1 {
		missedContent := an.text[an.lastBranch:an.here]
		if an.commentCtx != nil {
			an.commentCtx.Content += missedContent
		} else if an.valCtx != nil{
			an.valCtx.Text += missedContent
		} else {
			fmt.Println("There is text nowhere to go: ", missedContent)
		}
	}
}

func (an *Analyser) Spaces(spaces string) {
	an.StoreContent()
	if an.commentCtx != nil {
		an.commentCtx.Content += spaces
	} else if an.valCtx != nil {
		an.valCtx.TrailingSpaces = spaces
		an.EndVal()
	} else if an.stmtCtx != nil{
		an.stmtCtx.Indent += spaces
	} else {
		fmt.Println("Dunno what to do with spaces")
	}
}

func (an *Analyser) EnterCommentCtx() {
}

func (an *Analyser) NewComment(style string) {
	an.EnterCommentCtx()
	an.commentCtx.CommentStyle = style
}