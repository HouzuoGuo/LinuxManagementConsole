package txtedit
import "testing"

var input = `
first line # with some comment
	indent test
    second line
back to 0
`

func TestAnalyser(t *testing.T) {
	an := NewAnalyser(&AnalyserStyle{
			StmtContinue:make([]string, 0),
			StmtEnd:[]string{"\n"},
			CommentBegin:[]string{"#"},
			Quote: []string{"\""},
			BeginSectWithStmt:true,
			EndSectWithStmt:true},
	input)

	an.Analyse()
	DebugNode(an.Root, 0)
}