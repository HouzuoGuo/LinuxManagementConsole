package predeflex

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit/lexer"

var Sysconfig = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{"="},
	SectionStyle:                 lexer.SectionStyle{},
}

var Sysctl = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{"="},
	SectionStyle:                 lexer.SectionStyle{},
}

var Systemd = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{"="},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "[", OpeningSuffix: "]",
		ClosingPrefix: "", ClosingSuffix: "",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var CronAllowDeny = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Cron = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Hosts = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Login = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Nsswitch = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Httpd = lexer.LexerConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\"", "'"},
	TokenBreakMarkers:            []string{":"},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "<", OpeningSuffix: ">",
		ClosingPrefix: "</", ClosingSuffix: ">",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: true,
	},
}

var Named = lexer.LexerConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles: []lexer.CommentStyle{
		lexer.CommentStyle{Opening: "/*", Closing: "*/"},
		lexer.CommentStyle{Opening: "//", Closing: "\n"},
		lexer.CommentStyle{Opening: "#", Closing: "\n"}},
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
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: ";", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "(",
		ClosingPrefix: "", ClosingSuffix: ");",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var Dhcpd = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle: lexer.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "{",
		ClosingPrefix: "", ClosingSuffix: "}",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var Ntpd = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}

var Limits = lexer.LexerConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []lexer.CommentStyle{lexer.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 lexer.SectionStyle{},
}
