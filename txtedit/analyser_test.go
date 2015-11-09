package txtedit
import "testing"

var input = `
first 0
    indent 1
    second 2
back 3
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