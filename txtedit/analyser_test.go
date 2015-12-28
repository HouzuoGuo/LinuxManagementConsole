package txtedit

import (
	"fmt"
	"testing"
)

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

var input0 = `<A>
</A>`

/*
var input2 = `zone "." in {
    type hint;
    file "root.hint";
    forwarders { 192.0.2.1; 192.0.2.2; };
};
`
*/
var input2 = `a {
};`

func TestAnalyser(t *testing.T) {
	an := NewAnalyser(input, &AnalyserConfig{
		StatementContinuationMarkers: []string{"\\"},
		StatementEndingMarkers:       []string{"\n"},
		CommentBeginningMarkers:      []string{"#"},
		TextQuoteStyle:               []string{"\"", "'"},
		SectionBeginningPrefixes:     []string{"<"},
		SectionBeginningSuffixes:     []string{">"},
		SectionEndingPrefixes:        []string{"</"},
		SectionEndingSuffixes:        []string{">"},
		BeginSectionWithAStatement:   true,
		EndSectionWithAStatement:     true},
		&PrintDebugger{})

	an.Run()
	DebugNode(an.rootNode, 0)
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input {
		t.Fatal("no match")
	}
}

func TestAnalyser2(t *testing.T) {
	an := NewAnalyser(input2,
		&AnalyserConfig{
			StatementContinuationMarkers: []string{"\\"},
			StatementEndingMarkers:       []string{";"},
			CommentBeginningMarkers:      []string{"#"},
			TextQuoteStyle:               []string{"\"", "'"},
			SectionBeginningPrefixes:     []string{},
			SectionBeginningSuffixes:     []string{"{"},
			SectionEndingPrefixes:        []string{},
			SectionEndingSuffixes:        []string{"};"},
			BeginSectionWithAStatement:   true,
			EndSectionWithAStatement:     false},
		&PrintDebugger{})

	an.Run()
	DebugNode(an.rootNode, 0)
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
}
