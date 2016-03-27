package definedfmt

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit/analyser"

var Sysconfig = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{"="},
	SectionStyle:                 analyser.SectionStyle{},
}

var Sysctl = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{"="},
	SectionStyle:                 analyser.SectionStyle{},
}

var Systemd = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{"="},
	SectionStyle: analyser.SectionStyle{
		OpeningPrefix: "[", OpeningSuffix: "]",
		ClosingPrefix: "", ClosingSuffix: "",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var CronAllowDeny = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 analyser.SectionStyle{},
}

var Cron = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 analyser.SectionStyle{},
}

var Hosts = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 analyser.SectionStyle{},
}

var Login = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 analyser.SectionStyle{},
}

var Nsswitch = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 analyser.SectionStyle{},
}

var Httpd = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\"", "'"},
	TokenBreakMarkers:            []string{":"},
	SectionStyle: analyser.SectionStyle{
		OpeningPrefix: "<", OpeningSuffix: ">",
		ClosingPrefix: "</", ClosingSuffix: ">",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: true,
	},
}

var Named = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles: []analyser.CommentStyle{
		analyser.CommentStyle{Opening: "/*", Closing: "*/"},
		analyser.CommentStyle{Opening: "//", Closing: "\n"},
		analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:    []string{"\"", "'"},
	TokenBreakMarkers: []string{},
	SectionStyle: analyser.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "{",
		ClosingPrefix: "", ClosingSuffix: "};",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var NamedZone = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: ";", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	TokenBreakMarkers:            []string{},
	SectionStyle: analyser.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "(",
		ClosingPrefix: "", ClosingSuffix: ");",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var Dhcpd = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle: analyser.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "{",
		ClosingPrefix: "", ClosingSuffix: "}",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var Ntpd = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 analyser.SectionStyle{},
}

var Limits = analyser.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []analyser.CommentStyle{analyser.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	TokenBreakMarkers:            []string{},
	SectionStyle:                 analyser.SectionStyle{},
}
