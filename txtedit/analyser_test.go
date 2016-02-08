package txtedit

import (
	"fmt"
	"testing"
)

var input = `<a>
b
</c>`

func TestAnalyser(t *testing.T) {
	an := NewAnalyser(input, &AnalyserConfig{
		StatementContinuationMarkers: []string{"\\"},
		StatementEndingMarkers:       []string{"\n"},
		CommentStyles:                []CommentStyle{CommentStyle{Opening: "#", Closing: "\n"}},
		TextQuoteStyle:               []string{"\"", "'"},
		SectionStyle: SectionStyle{
			OpeningPrefix: "<", OpeningSuffix: ">",
			ClosingPrefix: "</", ClosingSuffix: ">",
			OpenSectionWithAStatement: true, CloseSectionWithAStatement: true,
		},
	},
		&PrintDebugger{})

	an.Run()
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	fmt.Println([]byte(input))
	fmt.Println([]byte(an.rootNode.TextString()))
	if an.rootNode.TextString() != input {
		t.Fatal("no match")
	}
}

var input2 = `a {
b;
};`

func TestAnalyser2(t *testing.T) {
	an := NewAnalyser(input2,
		&AnalyserConfig{
			StatementContinuationMarkers: []string{"\\"},
			StatementEndingMarkers:       []string{";"},
			CommentStyles: []CommentStyle{
				CommentStyle{Opening: "/*", Closing: "*/"},
				CommentStyle{Opening: "//", Closing: "\n"},
				CommentStyle{Opening: "#", Closing: "\n"}},
			TextQuoteStyle: []string{"\"", "'"},
			SectionStyle: SectionStyle{
				OpeningPrefix: "", OpeningSuffix: "{",
				ClosingPrefix: "", ClosingSuffix: "};",
				OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
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

var input3 = `[a]
b
c
[d]
e
f
[g]
h
[i]
[j]`

func TestAnalyser3(t *testing.T) {
	an := NewAnalyser(input3,
		&AnalyserConfig{
			StatementContinuationMarkers: []string{},
			StatementEndingMarkers:       []string{"\n"},
			CommentStyles:                []CommentStyle{CommentStyle{Opening: "#", Closing: "\n"}},
			TextQuoteStyle:               []string{"\""},
			SectionStyle: SectionStyle{
				OpeningPrefix: "[", OpeningSuffix: "]",
				ClosingPrefix: "", ClosingSuffix: "",
				OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
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

var input4 = `
#
# If defined, this command is run after removing a user.
# It should rebuild any NIS database etc. to remove the
# account from it.
#
USERDEL_POSTCMD	/usr/sbin/userdel-post.local`

func TestAnalyser4(t *testing.T) {
	an := NewAnalyser(input4,
		&AnalyserConfig{
			StatementContinuationMarkers: []string{},
			StatementEndingMarkers:       []string{"\n"},
			CommentStyles:                []CommentStyle{CommentStyle{Opening: "#", Closing: "\n"}},
			TextQuoteStyle:               []string{},
			SectionStyle:                 SectionStyle{},
		},
		&PrintDebugger{})

	an.Run()
	fmt.Println(DebugNode(an.rootNode, 0))
	fmt.Println("Reproduced:")
	fmt.Println(an.rootNode.TextString())
	if an.rootNode.TextString() != input4 {
		t.Fatal("no match")
	}
}
