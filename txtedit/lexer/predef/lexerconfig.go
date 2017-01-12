package predef

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit/lexer"

var Sysconfig = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{"="},
	SectionStyle:                 lexer.SectionStyle{},
}

var SysctlConf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{"="},
	SectionStyle:                 lexer.SectionStyle{},
}

var SystemdConf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{"="},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "[", OpeningSuffix: "]",
		ClosingPrefix: "", ClosingSuffix: "",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var CronAllow = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Crontab = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Hosts = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var LoginDefs = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Nsswitch = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var HttpdConf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\"", "'"},
	TokenBreakMarkers:            []string{":"},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "<", OpeningSuffix: ">",
		ClosingPrefix: "</", ClosingSuffix: ">",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: true,
	},
}

var NamedConf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles: []lexer.CommentStyle{
		{Opening: "/*", Closing: "*/"},
		{Opening: "//", Closing: "\n"},
		{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:    []string{"\"", "'"},
	TokenBreakMarkers: []string{},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "{",
		ClosingPrefix: "", ClosingSuffix: "};",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var NamedZone = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: ";", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "(",
		ClosingPrefix: "", ClosingSuffix: ");",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var DhcpdConf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "{",
		ClosingPrefix: "", ClosingSuffix: "}",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var NtpConf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var LimitsConf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var PostfixMainCf = lexer.LexerConfig{
	StatementContinuationMarkers: []string{"\n "},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{"="},
	SectionStyle:                 lexer.SectionStyle{},
}
