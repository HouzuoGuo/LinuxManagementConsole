package txtedit
import (
	"fmt"
	"strings"
)

type DebugPrint interface{
	Debug() string
}

type Val struct {
	QuoteStyle     string
	Text           string
	TrailingSpaces string
}

func (val *Val) Debug() string {
	return fmt.Sprintf("q%s '%s%s'", val.QuoteStyle, val.Text, val.TrailingSpaces)
}

type Comment struct {
	CommentStyle string
	Content      string
}

func (comment *Comment) Debug() string {
	return fmt.Sprintf("c%s '%s'", comment.CommentStyle, comment.Content)
}

type StmtContinue struct {
	Style string
}

func (cont *StmtContinue) Debug() string {
	return fmt.Sprintf("continue '%s'", cont.Style)
}

type Stmt struct {
	Indent string
	// Value or Comment or StmtContinue
	Pieces []interface{}
	End    string
}

func (stmt *Stmt) Debug() string {
	var out string
	for _, piece := range stmt.Pieces {
		out += "[" + piece.(DebugPrint).Debug() + "]"
	}
	return fmt.Sprintf("s'%s' %s e%v", stmt.Indent, out, []byte(stmt.End))
}

type Sect struct {
	Begin                    *Stmt
	BeginPrefix, BeginSuffix string
	EndPrefix, EndSuffix     string
	// Stmt or Sect
	Pieces                   []interface{}
	End                      *Stmt
}

type DocNode struct {
	Parent *DocNode
	// Stmt or Sect
	Obj    interface{}
	Leaves []*DocNode
}

type AnalyserStyle struct {
	StmtContinue                       []string
	StmtEnd                            []string
	CommentBegin                       []string
	Quote                              []string
	BeginSectWithStmt, EndSectWithStmt bool

	SectBeginPrefix                    []string
	SectBeginSuffix                    []string
	SectEndPrefix                      []string
	SectEndSuffix                      []string
}

type Analyser struct {
	Style            *AnalyserStyle
	Root             *DocNode

	text             string
	lastBranch, here int
	this             *DocNode

	valCtx           *Val
	commentCtx       *Comment
	stmtCtx          *Stmt
}

func NewAnalyser(style *AnalyserStyle, input string) (ret *Analyser) {
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
	if node.Obj == nil {
		fmt.Print(prefix + "Node - nil")
	} else {
		fmt.Print(prefix + "Node - ", node.Obj.(DebugPrint).Debug())
	}

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
	fmt.Println(root.Obj.(DebugPrint).Debug())
	fmt.Println("-->")
	for _, leaf := range root.Leaves {
		Print(leaf)
	}
}

func (an *Analyser) newSiblingIfNotNil() {
	if an.this.Obj == nil {
		return
	}
	if parent := an.this.Parent; parent == nil {
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

func (an *Analyser) newLeafIfNotNil() {
	if an.this.Obj == nil {
		return
	}
	newLeaf := &DocNode{Parent: an.this, Leaves: make([]*DocNode, 0, 8)}
	an.this.Leaves = append(an.this.Leaves, newLeaf)
	an.this = newLeaf
}

func (an *Analyser) NewComment(style string) {
	if an.commentCtx == nil {
		an.commentCtx = new(Comment)
		an.commentCtx.CommentStyle = style
	}
}

func (an *Analyser) EnterComment(style string) {
	if an.commentCtx == nil {
		an.NewComment(style)
	} else {
		return
	}
}

func (an *Analyser) EndComment() {
	if an.commentCtx == nil {
		return
	} else {
		fmt.Println("endcomment will store content")
		an.storeContent()
		an.EnterStmt()
		an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.commentCtx)
	}
	an.commentCtx = nil
}

func (an *Analyser) NewVal() {
	if an.valCtx == nil {
		an.valCtx = new(Val)
	}
}

func (an *Analyser) EnterVal() {
	if an.valCtx == nil {
		an.NewVal()
	} else {
		return
	}
}

func (an *Analyser) EndVal() {
	if an.valCtx == nil {
		return
	} else {
		fmt.Println("endval will store content")
		an.storeContent()
		an.EnterStmt()
		an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.valCtx)
	}
	an.valCtx = nil
}

func (an *Analyser) NewStmt() {
	fmt.Println("newstmt called")
	if an.stmtCtx == nil {
		an.newSiblingIfNotNil()
		an.stmtCtx = new(Stmt)
		an.this.Obj = an.stmtCtx
	}
}

func (an *Analyser) EnterStmt() {
	fmt.Println("enterstmt called")
	if an.stmtCtx == nil {
		an.NewStmt()
	} else {
		fmt.Println("enterstmt will do nothing")
		return
	}
}

func (an *Analyser) EndStmt() {
	an.storeContent()
	an.EndComment()
	an.EndVal()
	if an.stmtCtx == nil {
		return
	} else {
		an.storeContent()
		an.this.Obj = an.stmtCtx
		an.newSiblingIfNotNil()
	}
	fmt.Println("endstmt will set to nil")
	an.stmtCtx = nil
}

func (an *Analyser) storeContent() {
	if an.here - an.lastBranch > 1 {
		missedContent := an.text[an.lastBranch:an.here]
		if an.commentCtx != nil {
			an.commentCtx.Content += missedContent
		} else {
			an.EnterVal()
			an.valCtx.Text += missedContent
		}
		an.lastBranch = an.here
	}
}
func (an *Analyser) storeSpaces(spaces string) {
	an.storeContent()
	fmt.Println("About to store space '" +spaces + "'")
	if an.commentCtx != nil {
		an.commentCtx.Content += spaces
	} else if an.valCtx != nil {
		an.valCtx.TrailingSpaces = spaces
		an.EndVal()
	} else if an.stmtCtx != nil {
		fmt.Println("store space going to set indent")
		an.stmtCtx.Indent += spaces
	} else if an.stmtCtx == nil{
		fmt.Println("store space going to enter stmt")
		an.EnterStmt()
		an.stmtCtx.Indent += spaces
	} else {
		fmt.Println("Spaces have no where to go")
	}
}
