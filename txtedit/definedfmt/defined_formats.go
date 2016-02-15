package definedfmt

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit"

var Sysconfig = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	SectionStyle:                 txtedit.SectionStyle{},
}

var Sysctl = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	SectionStyle:                 txtedit.SectionStyle{},
}

var Systemd = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	SectionStyle: txtedit.SectionStyle{
		OpeningPrefix: "[", OpeningSuffix: "]",
		ClosingPrefix: "", ClosingSuffix: "",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var CronAllowDeny = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{}},
	TextQuoteStyle:               []string{},
	SectionStyle:                 txtedit.SectionStyle{},
}

var Cron = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	SectionStyle:                 txtedit.SectionStyle{},
}

var Hosts = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	SectionStyle:                 txtedit.SectionStyle{},
}

var Login = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	SectionStyle:                 txtedit.SectionStyle{},
}

var Nsswitch = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	SectionStyle:                 txtedit.SectionStyle{},
}

var Httpd = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\"", "'"},
	SectionStyle: txtedit.SectionStyle{
		OpeningPrefix: "<", OpeningSuffix: ">",
		ClosingPrefix: "</", ClosingSuffix: ">",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: true,
	},
}

var Named = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles: []txtedit.CommentStyle{
		txtedit.CommentStyle{Opening: "/*", Closing: "*/"},
		txtedit.CommentStyle{Opening: "//", Closing: "\n"},
		txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle: []string{"\"", "'"},
	SectionStyle: txtedit.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "{",
		ClosingPrefix: "", ClosingSuffix: "};",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var NamedZone = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: ";", Closing: "\n"}},
	TextQuoteStyle:               []string{},
	SectionStyle: txtedit.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "(",
		ClosingPrefix: "", ClosingSuffix: ");",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}

var Dhcpd = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{";\n", ";"},
	CommentStyles:                []txtedit.CommentStyle{txtedit.CommentStyle{Opening: "#", Closing: "\n"}},
	TextQuoteStyle:               []string{"\""},
	SectionStyle: txtedit.SectionStyle{
		OpeningPrefix: "", OpeningSuffix: "{",
		ClosingPrefix: "", ClosingSuffix: "}",
		OpenSectionWithAStatement: true, CloseSectionWithAStatement: false,
	},
}
