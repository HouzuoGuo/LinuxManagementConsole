package txtedit

import "testing"

var input = `
<SA>
	<SB>
		123
	</SB>
</SA>
456
`

var input2 = `
abc "def" 'ghi' #jkl
   mno

`

func TestAnalyser(t *testing.T) {
	an := NewAnalyser(&AnalyserStyle{
		StmtContinue:      []string{"\\"},
		StmtEnd:           []string{"\n"},
		CommentBegin:      []string{"#"},
		Quote:             []string{"\"", "'"},
		SectBeginPrefix:   []string{"<"},
		SectBeginSuffix:   []string{">"},
		SectEndPrefix:     []string{"</"},
		SectEndSuffix:     []string{">"},
		BeginSectWithStmt: true,
		EndSectWithStmt:   true},
		input)

	an.Analyse()
	DebugNode(an.Root, 0)
}

func TestAnalyser2(t *testing.T) {
	an := NewAnalyser(&AnalyserStyle{
		StmtContinue:      []string{"\\"},
		StmtEnd:           []string{"\n"},
		CommentBegin:      []string{"#"},
		Quote:             []string{"\"", "'"},
		SectBeginPrefix:   []string{"["},
		SectBeginSuffix:   []string{"]"},
		SectEndPrefix:     []string{""},
		SectEndSuffix:     []string{""},
		BeginSectWithStmt: true,
		EndSectWithStmt:   true},
		input2)

	an.Analyse()
	DebugNode(an.Root, 0)
}
