package txtedit

type Val struct {
	QuoteStyle string
	Text string
	TrailingSpaces string
}

type Comment struct {
	CommentStyle string
	Content string
}

type Stmt struct {
	Indent string
	// Value or Comment
	Pieces []interface{}
	End string
}

type Sect struct {
	Begin *Stmt
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
	BeginSectWithStmt bool

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
}

type Context struct {
	// Val, Comment, Stmt, or Sect
	this interface{}

}