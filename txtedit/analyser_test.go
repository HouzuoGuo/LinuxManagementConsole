package txtedit
import "testing"

var input = `abc def #ghi
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