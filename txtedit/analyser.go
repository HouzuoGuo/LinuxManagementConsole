package txtedit

import (
	"fmt"
	"strings"
)

type DebugPrint interface {
	Debug() string
}

type Val struct {
	QuoteStyle     string
	Text           string
	TrailingSpaces string
}

func (val *Val) Debug() string {
	return fmt.Sprintf("q%s %s%s", val.QuoteStyle, val.Text, val.TrailingSpaces)
}

type Comment struct {
	CommentStyle string
	Content      string
}

func (comment *Comment) Debug() string {
	return fmt.Sprintf("c%s %s", comment.CommentStyle, comment.Content)
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
	// Pieces inside section are DocNodes
	End                      *Stmt
}

func (sect *Sect) Debug() string {
	beginStmtStr := ""
	if sect.Begin != nil {
		beginStmtStr = sect.Begin.Debug()
	}
	endStmtStr := ""
	if sect.End != nil {
		endStmtStr = sect.End.Debug()
	}
	return fmt.Sprintf("Section %s %s %s, ends with %s %s %s",
		sect.BeginPrefix, beginStmtStr, sect.BeginSuffix,
		sect.EndPrefix, endStmtStr, sect.EndSuffix)
}

type DocNode struct {
	Parent *DocNode
	// Stmt or Sect
	Obj    interface{}
	Leaves []*DocNode
}

const (
	SECT_MATCH_BEGIN_PREFIX = 0
	SECT_MATCH_BEGIN_PREFIX_SUFFIX = 1
	SECT_MATCH_BEGIN_PREFIX_END_PREFIX = 2
	SECT_MATCH_ALL = 3
)

type SectionMatchStyle int

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

	SectMatchStyle                     SectionMatchStyle
}

func (style *AnalyserStyle) SetMatchStyle() {
	if len(style.SectEndSuffix) > 0 && len(style.SectEndPrefix) > 0 {
		style.SectMatchStyle = SECT_MATCH_ALL
	}else if len(style.SectEndPrefix) > 0 && len(style.SectBeginPrefix) > 0 {
		style.SectMatchStyle = SECT_MATCH_BEGIN_PREFIX_END_PREFIX
	}else if len(style.SectBeginSuffix) > 0 && len(style.SectBeginPrefix) > 0 {
		style.SectMatchStyle = SECT_MATCH_BEGIN_PREFIX_SUFFIX
	}else {
		style.SectMatchStyle = SECT_MATCH_BEGIN_PREFIX
	}
}

type Analyser struct {
	Style             *AnalyserStyle
	Root              *DocNode

	text              string
	lastBranch, here  int
	this              *DocNode
	ignoreNewStmtOnce bool

	valCtx            *Val
	commentCtx        *Comment
	stmtCtx           *Stmt
}

func NewAnalyser(style *AnalyserStyle, input string) (ret *Analyser) {
	ret = &Analyser{Style: style, text: input}
	ret.this = &DocNode{Parent: nil, Obj: nil, Leaves: make([]*DocNode, 0, 8)}
	ret.Root = ret.this
	ret.Style.SetMatchStyle()
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
		// root carries nothing
		// create a leaf that was the root
		rootAsLeaf := an.Root
		an.Root = &DocNode{Parent: nil, Leaves: make([]*DocNode, 0, 8)}
		an.Root.Leaves = append(an.Root.Leaves, rootAsLeaf)
		newLeaf := &DocNode{Parent: an.Root, Leaves: make([]*DocNode, 0, 8)}
		an.Root.Leaves = append(an.Root.Leaves, newLeaf)
		an.this = newLeaf
	} else {
		newLeaf := &DocNode{Parent: parent, Leaves: make([]*DocNode, 0, 8)}
		parent.Leaves = append(parent.Leaves, newLeaf)
		an.this = newLeaf
	}
}

func (an *Analyser) newLeaf() {
	newLeaf := &DocNode{Parent: an.this, Leaves: make([]*DocNode, 0, 8)}
	an.this.Leaves = append(an.this.Leaves, newLeaf)
	fmt.Println("New leaf is ", newLeaf, "and this is ", an.this, "and parent will be ", newLeaf.Parent)
	an.this = newLeaf
}

func (an *Analyser) BeginComment(style string) {
	if an.commentCtx == nil {
		an.commentCtx = new(Comment)
		an.commentCtx.CommentStyle = style
	}
}

func (an *Analyser) EnterComment(style string) {
	if an.commentCtx == nil {
		an.BeginComment(style)
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
	an.EndComment()
	an.EndVal()
	if an.ignoreNewStmtOnce {
		an.ignoreNewStmtOnce = false
		return
	}
	if an.stmtCtx == nil {
		return
	} else {
		an.storeContent()
		an.this.Obj = an.stmtCtx
		an.newSiblingIfNotNil()
	}
	an.stmtCtx = nil
}

func (an *Analyser) storeContent() {
	if an.here - an.lastBranch > 0 {
		missedContent := an.text[an.lastBranch:an.here]
		if an.commentCtx != nil {
			fmt.Println("missed content ", missedContent, "will be stored in comment")
			an.commentCtx.Content += missedContent
		} else {
			fmt.Println("missed content ", missedContent, " will be stored in val")
			an.EnterVal()
			an.valCtx.Text += missedContent
		}
		an.lastBranch = an.here
	}
}
func (an *Analyser) storeSpaces(spaces string) {
	an.storeContent()
	fmt.Println("About to store space '" + spaces + "'")
	if an.ignoreNewStmtOnce {
		fmt.Println("spaces are going to new val")
		an.EndVal()
		an.NewVal()
		an.valCtx.TrailingSpaces += spaces
		an.EndVal()
	} else if an.commentCtx != nil {
		an.commentCtx.Content += spaces
	} else if an.valCtx != nil {
		an.valCtx.TrailingSpaces = spaces
		an.EndVal()
	} else if an.stmtCtx != nil {
		fmt.Println("store space going to set indent")
		an.stmtCtx.Indent += spaces
	} else if an.stmtCtx == nil {
		fmt.Println("store space going to enter stmt")
		an.EnterStmt()
		an.stmtCtx.Indent += spaces
	} else {
		fmt.Println("Spaces have no where to go")
	}
}

func (an *Analyser) ContinueStmt(style string) {
	if an.valCtx != nil && an.valCtx.QuoteStyle != "" {
		an.valCtx.Text += style
		return
	}
	an.storeContent()
	an.EndComment()
	an.EndVal()
	an.EnterStmt()
	an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, &StmtContinue{Style:style})
	an.ignoreNewStmtOnce = true
}

func (an *Analyser) NewSection() {
	an.EndStmt()

	an.this.Obj = new(Sect)
	fmt.Println("new section creates new leaf,", an.this.Obj)
	an.newLeaf()
}

func (an *Analyser) IsSect() bool {
	fmt.Println("issect, this is ", an.this)
	if an.this.Parent == nil || an.this.Parent.Obj == nil {
		return false
	} else if _, isSect := an.this.Parent.Obj.(*Sect); isSect {
		return true
	} else {
		return false
	}
}

func (an *Analyser) FindThisLeaf() int {
	if an.this.Parent == nil {
		return -1
	}
	for i, leaf := range an.this.Parent.Leaves {
		if leaf == an.this {
			return i
		}
	}
	return -1
}

func (an *Analyser) GetPreviousLeaf() *Stmt {
	thisLeaf := an.FindThisLeaf()
	if thisLeaf == -1 {
		return nil
	}
	if thisLeaf == 0 {
		return nil
	}
	prevLeaf := an.this.Parent.Leaves[thisLeaf - 1]
	if stmt, ok := prevLeaf.Obj.(*Stmt); ok {
		return stmt
	}
	return nil
}

func (an *Analyser) EndSection() {
	if sect, isSect := an.this.Obj.(*Sect); isSect {
		fmt.Println("section ends here")
		minNumLeaves := 0
		if an.Style.BeginSectWithStmt {
			if len(an.Style.SectBeginSuffix) == 0 {
				sect.Begin = an.GetPreviousLeaf()
			} else {
				firstLeaf := an.this.Leaves[0]
				if stmt, ok := firstLeaf.Obj.(*Stmt); ok {
					fmt.Println("successfully set sect.begin")
					an.this.Leaves = an.this.Leaves[1:]
					sect.Begin = stmt
					minNumLeaves++
				}
			}
		}
		if an.Style.EndSectWithStmt {
			if len(an.Style.SectEndSuffix) > 0 {
				if len(an.this.Leaves) > minNumLeaves {
					lastLeaf := an.this.Leaves[len(an.this.Leaves) - 1]
					if stmt, ok := lastLeaf.Obj.(*Stmt); ok {
						fmt.Println("successfully set sect.end")
						an.this.Leaves = an.this.Leaves[0:len(an.this.Leaves) - 1]
						sect.End = stmt
					}
				}
			}
		}
	}else {
		fmt.Println("this is not a section but it ends here, why?")
	}
	an.this = an.this.Parent

	fmt.Println("section ends here")
}

const (
	SECT_STATE_NONE = 0
	SECT_STATE_BEFORE_BEGIN = 1
	SECT_STATE_BEGIN_PREFIX = 2
	SECT_STATE_BEGIN_SUFFIX = 3
	SECT_STATE_END_PREFIX = 4
	SECT_STATE_END_NOW = 5
)

type SectionState int

func (an *Analyser) GetSectionState() (SectionState, *Sect) {
	if !an.IsSect() {
		fmt.Println("not in section!!")
		return SECT_STATE_NONE, nil
	}
	sect := an.this.Parent.Obj.(*Sect)
	switch an.Style.SectMatchStyle {
	case SECT_MATCH_BEGIN_PREFIX:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	case SECT_MATCH_BEGIN_PREFIX_SUFFIX:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else if sect.BeginSuffix == "" {
			return SECT_STATE_BEGIN_PREFIX, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	case SECT_MATCH_BEGIN_PREFIX_END_PREFIX:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else if sect.EndPrefix == "" {
			return SECT_STATE_BEGIN_SUFFIX, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	case SECT_MATCH_ALL:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else if sect.BeginSuffix == "" {
			return SECT_STATE_BEGIN_PREFIX, sect
		} else if sect.EndPrefix == "" {
			return SECT_STATE_BEGIN_SUFFIX, sect
		} else if sect.EndSuffix == "" {
			return SECT_STATE_END_PREFIX, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	default:
		return SECT_STATE_NONE, sect
	}
}

func (an *Analyser) BeginSectionSetPrefix(style string) {
	state, sect := an.GetSectionState()
	fmt.Println("prefix state is ", state)
	if  state == SECT_STATE_END_NOW {
		sect.BeginPrefix = style
		an.storeContent()
		an.EndSection()
	} else if state > SECT_STATE_BEGIN_PREFIX {
		an.storeContent()
	} else {
		fmt.Println("New section is going to be created")
		an.NewSection()
		an.this.Parent.Obj.(*Sect).BeginPrefix = style
		an.storeContent()
	}
}
func (an *Analyser) BeginSectionSetSuffix(style string) {
	if state, sect := an.GetSectionState(); state == SECT_STATE_END_NOW {
		sect.BeginSuffix = style
		an.storeContent()
		an.EndSection()
	}else if state < SECT_STATE_BEGIN_PREFIX || state > SECT_STATE_BEGIN_SUFFIX {
		an.storeContent()
	} else {
		sect.BeginSuffix = style
		an.storeContent()
	}
}
func (an *Analyser) EndSectionSetPrefix(style string) {
	if state, sect := an.GetSectionState(); state == SECT_STATE_END_NOW {
		sect.EndPrefix = style
		an.storeContent()
		an.EndSection()
	}else if state < SECT_STATE_BEGIN_SUFFIX || state > SECT_STATE_END_PREFIX {
		an.storeContent()
	} else {
		sect.EndPrefix = style
		an.storeContent()
	}

}
func (an *Analyser) EndSectionSetSuffix(style string) {
	if state, sect := an.GetSectionState(); state == SECT_STATE_END_NOW {
		sect.EndSuffix = style
		an.storeContent()
		an.EndSection()
	}else if state < SECT_STATE_END_PREFIX {
		an.storeContent()
	} else {
		sect.EndPrefix = style
		an.storeContent()
	}
}