package txtedit
import "testing"

var input = `
# this is comment
abc 123
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