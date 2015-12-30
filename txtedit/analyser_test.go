package txtedit

import (
	"fmt"
	"testing"
)

var input = `DirectoryIndex index.html index.html.var
<Files ~ "^\.ht">
    <IfModule mod_access_compat.c>
        Order allow,deny
        Deny from all
    </IfModule>
</Files>
`

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
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input {
		t.Fatal("no match")
	}
}

var input2 = `zone "." in {
    type hint;
    file "root.hint";
    forwarders { 192.0.2.1; 192.0.2.2; };
};`

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
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input2 {
		t.Fatal("no match")
	}
}

var input3 = `# See systemd-system.conf(5) for details.

[Manager]
#LogLevel=info
#LogTarget=journal-or-kmsg
[Journald]
haha
`

func TestAnalyser3(t *testing.T) {
	an := NewAnalyser(input3,
		&AnalyserConfig{
			StatementContinuationMarkers: []string{},
			StatementEndingMarkers:       []string{"\n"},
			CommentBeginningMarkers:      []string{"#"},
			TextQuoteStyle:               []string{"\""},
			SectionBeginningPrefixes:     []string{"["},
			SectionBeginningSuffixes:     []string{"]"},
			SectionEndingPrefixes:        []string{},
			SectionEndingSuffixes:        []string{},
			BeginSectionWithAStatement:   true,
			EndSectionWithAStatement:     false},
		&PrintDebugger{})

	an.Run()
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input3 {
		t.Fatal("no match")
	}
}
