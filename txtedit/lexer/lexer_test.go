package lexer

import (
	"fmt"
	"testing"
)

var input = `<a>
b
<c>
d
</c>
</a>`

func TestLexer(t *testing.T) {
	an := NewLexer(input, &LexerConfig{
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

var input2 = `a{
#b
};
#abc`

func TestLexer2(t *testing.T) {
	an := NewLexer(input2,
		&LexerConfig{
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

func TestLexer3(t *testing.T) {
	an := NewLexer(input3,
		&LexerConfig{
			StatementContinuationMarkers: []string{},
			StatementEndingMarkers:       []string{"\n"},
			CommentStyles:                []CommentStyle{CommentStyle{Opening: "#", Closing: "\n"}},
			TextQuoteStyle:               []string{"\""},
			TokenBreakMarkers:            []string{"="},
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

func TestLexer4(t *testing.T) {
	an := NewLexer(input4,
		&LexerConfig{
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
