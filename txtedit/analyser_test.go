package txtedit

import "testing"

var input = `
<SA>
	<SB>
		<SC>
			"123"
			'456'
		</SC>
		789
		012
	</SB>
	345
	<SC>
		678
	</SC>
	901
</SA>
234
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
