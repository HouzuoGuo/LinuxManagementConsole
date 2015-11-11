package txtedit

import "testing"

var input = "here 123 \\   \n\n  there"

func TestAnalyser(t *testing.T) {
	an := NewAnalyser(&AnalyserStyle{
		StmtContinue:      []string{"\\"},
		StmtEnd:           []string{"\n"},
		CommentBegin:      []string{"#"},
		Quote:             []string{"\"", "'"},
		BeginSectWithStmt: true,
		EndSectWithStmt:   true},
		input)

	an.Analyse()
	DebugNode(an.Root, 0)
}
