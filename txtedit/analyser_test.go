package txtedit
import "testing"

var input = `
first 0 # haha
    indent 1 # haha hoho
    second 2 # ey there
back 3
# un ho
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