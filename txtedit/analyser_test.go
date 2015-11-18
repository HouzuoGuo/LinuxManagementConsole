package txtedit

import "testing"

var input = `
<SA>
	Abc
</SA>
`

func TestAnalyser(t *testing.T) {
	an := NewAnalyser(&AnalyserStyle{
		StmtContinue:      []string{"\\"},
		StmtEnd:           []string{"\n"},
		CommentBegin:      []string{"#"},
		Quote:             []string{"\"", "'"},
		SectBeginPrefix: []string{"<"},
		SectBeginSuffix:[]string{">"},
		SectEndPrefix: []string{"</"},
		SectEndSuffix:[]string{">"},
		BeginSectWithStmt: true,
		EndSectWithStmt:   true},
		input)

	an.Analyse()
	DebugNode(an.Root, 0)
}
