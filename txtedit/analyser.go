package txtedit
import "fmt"

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

	inQuote bool
	contStmt bool
	valCtx *Val
	commentCtx *Comment
	stmtCtx *Stmt
	nodeCtx *DocNode
}

func (an *Analyser) NewStmt() {
	if an.stmtCtx != nil {
		an.nodeCtx.Leaves = append(an.nodeCtx.Leaves, an.stmtCtx)
	}
	an.stmtCtx = new(Stmt)
}

func (an *Analyser) EndVal() {
	if an.valCtx == nil {
		return // do nothing
	}
	if an.stmtCtx == nil {
		an.NewStmt()
	} else {
		an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.valCtx)
	}
	an.valCtx = nil
	an.inQuote = false
}

func (an *Analyser) NewComment(style string) {
	if an.valCtx != nil {
		an.EndVal()
	}
	if an.commentCtx != nil {
		return // already a comment
	}
	an.commentCtx = new(Comment)
}

func (an *Analyser) NewVal() {
	if an.valCtx != nil {
		if an.stmtCtx == nil {
			an.NewStmt()
		}
		an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, an.valCtx)
	}
	an.valCtx = new(Val)
	an.inQuote = false
}

func (an *Analyser) SetQuote(style string) {
	if an.valCtx == nil {
		an.NewVal()
	} else {
		an.valCtx.QuoteStyle = style
	}
}

func (an *Analyser) SetTrailingSpacesOrIndent(spaces string) {
	if an.valCtx == nil {
		if an.stmtCtx == nil {
			an.NewStmt()
			an.stmtCtx.Indent = spaces
		} else if len(an.stmtCtx.Pieces) > 0{
			lastPiece := an.stmtCtx.Pieces[len(an.stmtCtx.Pieces - 1)]
			switch t := lastPiece.(type) {
			case Comment:
				t.Content += spaces
			case Val:
				t.TrailingSpaces += spaces
			}
		} else {
			an.stmtCtx.Indent = spaces
		}
	} else {
		if !an.inQuote {
			an.valCtx.TrailingSpaces = spaces
		}
	}
}

func (an *Analyser) ContinueStmt(style string) {
	an.EndVal()
	if an.stmtCtx == nil {
		an.NewStmt()
	}
	an.stmtCtx.Pieces = append(an.stmtCtx.Pieces, &StmtContinue{style})
}

func (an *Analyser) NewLeaf() {
	if an.nodeCtx != nil {
		an.nodeCtx.Parent.Leaves = append(an.nodeCtx.Parent.Leaves, an.nodeCtx)
	}
	an.nodeCtx = &DocNode{Parent:an.nodeCtx.Parent}
}

func (an *Analyser) NewBranch() {
	if an.nodeCtx != nil {
		an.nodeCtx.Parent.Leaves = append(an.nodeCtx.Parent.Leaves, an.nodeCtx)
	}
	an.nodeCtx = &DocNode{Parent: an.nodeCtx}
}

func (an *Analyser) EndStmt(style string) {
	if an.stmtCtx == nil {
		fmt.Println("!!ending a statement without starting one!!")
		return
	}
	an.stmtCtx.End = style
	if an.nodeCtx == nil {
		an.NewLeaf()
		an.nodeCtx.Obj = an.stmtCtx
	}
	an.stmtCtx = nil
}

func (an *Analyser) NewSection(style string) {
	if an.nodeCtx.Obj == nil {
		an.NewLeaf()
		an.nodeCtx.Obj = &Sect{BeginPrefix:style}
	} else {
		switch an.nodeCtx.Obj.(type) {
		case *Sect:
			an.NewBranch()
			an.nodeCtx.Obj = &Sect{BeginPrefix:style}
		case *Stmt:
			if an.Style.BeginSectWithStmt {
				an.NewBranch()
			} else {
				an.EndStmt("")
				an.NewLeaf()
			}
			an.nodeCtx.Obj = &Sect{BeginPrefix: style}
		}
	}
}

func (an *Analyser) NewValIfNotSect(content string) (*Sect ){
	sect, isSect := an.nodeCtx.Obj.(*Sect)
	if an.nodeCtx.Obj == nil || !isSect {
		if an.stmtCtx == nil {
			an.NewStmt()
		}
		if an.valCtx == nil {
			an.NewVal()
		}
		an.valCtx.Text = content
		return sect
	}
	return nil
}

func (an *Analyser) SetSectBeginSuffix(style string) {
	sect := an.NewValIfNotSect(style)
	if sect != nil {
		sect.BeginSuffix = style
	}
}

func (an *Analyser) EndSect() {
	???
}

func (an *Analyser) SetSectEndPrefix(style string) {
	sect := an.NewValIfNotSect(style)
	if sect != nil {
		if an.Style.EndSectWithStmt {
			???
		} else {
			???
		}
		sect.EndPrefix = style
	}
}

func (an *Analyser) SetSectEndSuffix(style string) {
	sect := an.NewValIfNotSect(style)
	if sect != nil {
		if an.Style.EndSectWithStmt {
			???
		}
	}
}