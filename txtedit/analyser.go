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
	return fmt.Sprintf("Quote[%s] Text[%s] Trailing[%s]", val.QuoteStyle, val.Text, val.TrailingSpaces)
}

type Comment struct {
	CommentStyle string
	Content      string
}

func (comment *Comment) Debug() string {
	return fmt.Sprintf("Comment[%s] Content[%s]", comment.CommentStyle, comment.Content)
}

type StmtContinue struct {
	Style string
}

func (cont *StmtContinue) Debug() string {
	return fmt.Sprintf("Continue[%s]", cont.Style)
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
	return fmt.Sprintf("Indent[%s] Pieces[%s] End[%v]", stmt.Indent, out, []byte(stmt.End))
}

type Sect struct {
	Begin                    *Stmt
	BeginPrefix, BeginSuffix string
	EndPrefix, EndSuffix     string
	// Pieces inside section are DocNodes
	End *Stmt
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
	return fmt.Sprintf("Section[%s][%s][%s] End[%s][%s][%s]",
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
	SECT_MATCH_BEGIN_PREFIX            = 0
	SECT_MATCH_BEGIN_PREFIX_SUFFIX     = 1
	SECT_MATCH_BEGIN_PREFIX_END_PREFIX = 2
	SECT_MATCH_ALL                     = 3
)

type SectionMatchStyle int

type AnalyserStyle struct {
	StmtContinue                       []string
	StmtEnd                            []string
	CommentBegin                       []string
	Quote                              []string
	BeginSectWithStmt, EndSectWithStmt bool

	SectBeginPrefix []string
	SectBeginSuffix []string
	SectEndPrefix   []string
	SectEndSuffix   []string

	SectMatchStyle         SectionMatchStyle
	AmbiguousSectionSuffix bool
}

func (style *AnalyserStyle) SetMatchStyle() {
	if len(style.SectEndSuffix) > 0 && len(style.SectEndPrefix) > 0 {
		style.SectMatchStyle = SECT_MATCH_ALL
	} else if len(style.SectEndPrefix) > 0 && len(style.SectBeginPrefix) > 0 {
		style.SectMatchStyle = SECT_MATCH_BEGIN_PREFIX_END_PREFIX
	} else if len(style.SectBeginSuffix) > 0 && len(style.SectBeginPrefix) > 0 {
		style.SectMatchStyle = SECT_MATCH_BEGIN_PREFIX_SUFFIX
	} else {
		style.SectMatchStyle = SECT_MATCH_BEGIN_PREFIX
	}
	for _, style1 := range style.SectBeginSuffix {
		for _, style2 := range style.SectEndSuffix {
			if style1 == style2 {
				style.AmbiguousSectionSuffix = true
			}
		}
	}
	fmt.Println("Analyser match style is ", style.SectMatchStyle)
}

type Analyser struct {
	Style *AnalyserStyle
	Root  *DocNode

	text              string
	lastBranch, here  int
	this              *DocNode
	ignoreNewStmtOnce bool

	valCtx     *Val
	commentCtx *Comment
	stmtCtx    *Stmt
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
		fmt.Print(prefix+"Node - ", node.Obj.(DebugPrint).Debug())
	}

	if len(node.Leaves) > 0 {
		fmt.Println(" -->")
		for _, leaf := range node.Leaves {
			DebugNode(leaf, indent+2)
		}
	} else {
		fmt.Println()
	}
}

func Print(root *DocNode) {
	fmt.Println(root.Obj.(DebugPrint).Debug())
	fmt.Println(">>>")
	for _, leaf := range root.Leaves {
		Print(leaf)
	}
}

func (an *Analyser) newSiblingIfNotNil() {
	if an.this.Obj == nil {
		fmt.Println("No new sibling because obj is nil")
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
		fmt.Println("New sibling for root has been created")
	} else {
		newLeaf := &DocNode{Parent: parent, Leaves: make([]*DocNode, 0, 8)}
		parent.Leaves = append(parent.Leaves, newLeaf)
		an.this = newLeaf
		fmt.Println("New sibling for non-root node")
	}
}

func (an *Analyser) newLeaf() {
	newLeaf := &DocNode{Parent: an.this, Leaves: make([]*DocNode, 0, 8)}
	an.this.Leaves = append(an.this.Leaves, newLeaf)
	fmt.Println("New leaf is created, this is ", an.this, ", new leaf is ", newLeaf)
	an.this = newLeaf
}

func (an *Analyser) BeginComment(style string) {
	if an.commentCtx == nil {
		an.commentCtx = new(Comment)
		an.commentCtx.CommentStyle = style
		fmt.Println("New comment is created")
	}
}

func (an *Analyser) EnterComment(style string) {
	an.EnterStmt()
	if an.commentCtx == nil {
		an.BeginComment(style)
	} else {
		fmt.Println("EnterComment does nothing")
	}
}

func (an *Analyser) EndComment() {
	if an.commentCtx == nil {
		fmt.Println("EndComment does nothing")
		return
	} else {
		fmt.Println("EndComment will store content")
		an.storeContent()
		an.EnterStmt()
		an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.commentCtx)
	}
	an.commentCtx = nil
}

func (an *Analyser) NewVal() {
	if an.valCtx == nil {
		fmt.Println("New val is created")
		an.valCtx = new(Val)
	}
}

func (an *Analyser) EnterVal() {
	an.EnterStmt()
	if an.valCtx == nil {
		an.NewVal()
	} else {
		fmt.Println("EnterVal does nothing")
	}
}

func (an *Analyser) EndVal() {
	if an.valCtx == nil {
		fmt.Println("EndVal does nothing")
		return
	} else {
		fmt.Println("EndVal will store content")
		an.storeContent()
		an.EnterStmt()
		an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.valCtx)
	}
	an.valCtx = nil
}

func (an *Analyser) NewStmt() {
	if an.stmtCtx == nil {
		fmt.Println("New stmt will be created")
		an.newSiblingIfNotNil()
		an.stmtCtx = new(Stmt)
		an.this.Obj = an.stmtCtx
		fmt.Println("New stmt is created")
	}
}

func (an *Analyser) EnterStmt() {
	if an.stmtCtx == nil {
		an.NewStmt()
	} else {
		fmt.Println("EnterStmt does nothing")
		return
	}
}

func (an *Analyser) EndStmt() {
	an.storeContent()
	an.EndComment()
	an.EndVal()
	if an.ignoreNewStmtOnce {
		an.ignoreNewStmtOnce = false
		fmt.Println("EndStmt does nothing because flag is true")
		return
	}
	if an.stmtCtx == nil {
		fmt.Println("EndStmt does nothing because stmt is nil")
		return
	} else {
		if an.commentCtx != nil {
			an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.commentCtx)
		}
		if an.valCtx != nil {
			an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.valCtx)
		}
		an.this.Obj = an.stmtCtx
		an.newSiblingIfNotNil()
	}
	an.commentCtx = nil
	an.valCtx = nil
	an.stmtCtx = nil
}

func (an *Analyser) storeContent() {
	if an.here-an.lastBranch > 0 {
		missedContent := an.text[an.lastBranch:an.here]
		if an.commentCtx != nil {
			fmt.Println("missed content", missedContent, "will be stored in comment")
			an.commentCtx.Content += missedContent
		} else {
			fmt.Println("missed content", missedContent, "will be stored in val")
			an.EnterVal()
			an.valCtx.Text += missedContent
		}
		an.lastBranch = an.here
	} else {
		fmt.Println("storeContent does nothing")
	}
}
func (an *Analyser) storeSpaces(spaces string) {
	an.storeContent()
	fmt.Println("About to store space '" + spaces + "'")
	if an.ignoreNewStmtOnce {
		fmt.Println("Spaces are going into new val")
		an.EndVal()
		an.NewVal()
		an.valCtx.TrailingSpaces += spaces
		an.EndVal()
	} else if an.commentCtx != nil {
		fmt.Println("Spaces are going into context comment")
		an.commentCtx.Content += spaces
	} else if an.valCtx != nil {
		fmt.Println("Spaces are going into value trailing spaces")
		an.valCtx.TrailingSpaces = spaces
		an.EndVal()
	} else if an.stmtCtx != nil {
		fmt.Println("Spaces set indent")
		an.stmtCtx.Indent += spaces
	} else if an.stmtCtx == nil {
		fmt.Println("Spaces set indent and makes a new statement")
		an.EnterStmt()
		an.stmtCtx.Indent += spaces
	} else {
		fmt.Println("Spaces have no where to go")
	}
}

func (an *Analyser) ContinueStmt(style string) {
	if an.valCtx != nil && an.valCtx.QuoteStyle != "" {
		an.valCtx.Text += style
		fmt.Println("Continue statement mark goes into value")
		return
	}
	an.storeContent()
	an.EndComment()
	an.EndVal()
	an.EnterStmt()
	an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, &StmtContinue{Style: style})
	an.ignoreNewStmtOnce = true
	fmt.Println("Continue statement flag is set")
}

func (an *Analyser) NewSection() {
	an.EndStmt()
	fmt.Println("NewSection from here:", an.this)
	if an.this == an.Root {
		an.newLeaf()
	} else {
		an.newSiblingIfNotNil()
	}
	an.this.Obj = new(Sect)
	an.newLeaf()
}

func (an *Analyser) IsSect() bool {
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
	prevLeaf := an.this.Parent.Leaves[thisLeaf-1]
	if stmt, ok := prevLeaf.Obj.(*Stmt); ok {
		return stmt
	}
	return nil
}

func (an *Analyser) EndSection() {
	if _, sect := an.GetSectionState(); sect == nil {
		fmt.Println("this is not a section but it ends here, why?")
	} else {
		fmt.Println("section ends here, saving the latest statement")
		an.EndStmt()
		an.this = an.this.Parent
		// an.this is now the parent - section object
		// Remove the last leaf if it holds no object
		if an.this.Leaves[len(an.this.Leaves)-1].Obj == nil {
			an.this.Leaves = an.this.Leaves[:len(an.this.Leaves)-1]
		}
		fmt.Println("section leaves:")
		for _, leaf := range an.this.Leaves {
			fmt.Println(leaf.Obj.(DebugPrint).Debug())
		}
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
					lastLeaf := an.this.Leaves[len(an.this.Leaves)-1]
					if stmt, ok := lastLeaf.Obj.(*Stmt); ok {
						fmt.Println("successfully set sect.end")
						an.this.Leaves = an.this.Leaves[0 : len(an.this.Leaves)-1]
						sect.End = stmt
					}
				}
			}
		}
		fmt.Println("Section has ended")
		fmt.Println(an.this.Parent, an.Root)
		fmt.Println(an.this, an.Root.Leaves[0])
		fmt.Println(an.this.Obj, an.Root.Leaves[0].Obj)
		an.newSiblingIfNotNil()
	}
}

const (
	SECT_STATE_NONE         = 0
	SECT_STATE_BEFORE_BEGIN = 1
	SECT_STATE_BEGIN_PREFIX = 2
	SECT_STATE_BEGIN_SUFFIX = 3
	SECT_STATE_END_PREFIX   = 4
	SECT_STATE_END_NOW      = 5
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
	if state == SECT_STATE_END_NOW {
		sect.BeginPrefix = style
		an.storeContent()
		an.EndSection()
	} else if state > SECT_STATE_BEGIN_PREFIX {
		an.storeContent()
	} else {
		fmt.Println("New section is going to be created")
		an.NewSection()
		fmt.Println(an.this.Parent, an.this)
		an.this.Parent.Obj.(*Sect).BeginPrefix = style
		an.storeContent()
	}
}

func (an *Analyser) BeginSectionSetSuffix(style string) {
	if state, sect := an.GetSectionState(); state == SECT_STATE_END_NOW {
		fmt.Println("BeginSectionSetSuffix: set style and end")
		sect.BeginSuffix = style
		an.storeContent()
		an.EndSection()
	} else if state < SECT_STATE_BEGIN_PREFIX || state > SECT_STATE_BEGIN_SUFFIX {
		fmt.Println("BeginSectionSetSuffix: no match, store content")
		an.storeContent()
	} else {
		fmt.Println("BeginSectionSetSuffix: only set style")
		sect.BeginSuffix = style
		an.storeContent()
	}
}

func (an *Analyser) EndSectionSetPrefix(style string) {
	if state, sect := an.GetSectionState(); state == SECT_STATE_END_NOW {
		fmt.Println("EndSectionSetPrefix: set style and end")
		sect.EndPrefix = style
		an.storeContent()
		an.EndSection()
	} else if state < SECT_STATE_BEGIN_SUFFIX || state > SECT_STATE_END_PREFIX {
		fmt.Println("EndSectionSetPrefix: no match, store content")
		an.storeContent()
	} else {
		fmt.Println("EndSectionSetPrefix: only set style")
		sect.EndPrefix = style
		an.storeContent()
	}
}

func (an *Analyser) EndSectionSetSuffix(style string) {
	if state, sect := an.GetSectionState(); state >= SECT_STATE_END_PREFIX {
		fmt.Println("EndSectionSetSuffix: set style and end")
		sect.EndSuffix = style
		an.storeContent()
		an.EndSection()
	} else if state < SECT_STATE_END_PREFIX && an.Style.AmbiguousSectionSuffix {
		fmt.Println("EndSectionSetSuffix: ambiguous suffix")
		an.BeginSectionSetSuffix(style)
	} else {
		fmt.Println("EndSectionSetSuffix: only set style")
		sect.EndSuffix = style
		an.storeContent()
	}
}
