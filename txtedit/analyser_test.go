package txtedit

import (
	"fmt"
	"testing"
)

var input = `# http://httpd.apache.org/docs/2.4/mod/core.html#options`

func TestAnalyser(t *testing.T) {
	an := NewAnalyser(input, &AnalyserConfig{
		StatementContinuationMarkers: []string{"\\"},
		StatementEndingMarkers:       []string{"\n"},
		CommentStyle:                 CommentStyle{Opening: "#"},
		TextQuoteStyle:               []string{"\"", "'"},
		SectionStyle: SectionStyle{
			OpeningPrefix: "<", OpeningSuffix: ">",
			ClosingPrefix: "</", ClosingSuffix: ">",
			BeginSectionWithAStatement: true, EndSectionWithAStatement: true,
		},
	},
		&PrintDebugger{})

	an.Run()
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input {
		t.Fatal("no match")
	}
}

var input2 = `#a
b {
c;
};`

func TestAnalyser2(t *testing.T) {
	an := NewAnalyser(input2,
		&AnalyserConfig{
			StatementContinuationMarkers: []string{"\\"},
			StatementEndingMarkers:       []string{";"},
			CommentStyle:                 CommentStyle{Opening: "#"},
			TextQuoteStyle:               []string{"\"", "'"},
			SectionStyle: SectionStyle{
				OpeningPrefix: "", OpeningSuffix: "{",
				ClosingPrefix: "", ClosingSuffix: "};",
				BeginSectionWithAStatement: true, EndSectionWithAStatement: false,
			},
		},
		&PrintDebugger{})

	an.Run()
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input2 {
		t.Fatal("no match")
	}
}

var input3 = `## See systemd-system.conf(5) for details.
# a="a"
[Manager]
#LogLevel=info
#LogTarget=journal-or-kmsg
[Journald]
[]

haha

`

func TestAnalyser3(t *testing.T) {
	an := NewAnalyser(input3,
		&AnalyserConfig{
			StatementContinuationMarkers: []string{},
			StatementEndingMarkers:       []string{"\n"},
			CommentStyle:                 CommentStyle{Opening: "#"},
			TextQuoteStyle:               []string{"\""},
			SectionStyle: SectionStyle{
				OpeningPrefix: "[", OpeningSuffix: "]",
				ClosingPrefix: "", ClosingSuffix: "",
				BeginSectionWithAStatement: true, EndSectionWithAStatement: false,
			},
		},
		&PrintDebugger{})

	an.Run()
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input3 {
		t.Fatal("no match")
	}
}
