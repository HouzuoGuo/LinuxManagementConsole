package definedfmt

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit"

var Sysconfig = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{"\""},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   false,
	EndSectionWithAStatement:     false}

var Sysctl = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   false,
	EndSectionWithAStatement:     false}

var Systemd = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{"\""},
	SectionBeginningPrefixes:     []string{"["},
	SectionBeginningSuffixes:     []string{"]"},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   true,
	EndSectionWithAStatement:     false}

var CronAllowDeny = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{},
	TextQuoteStyle:               []string{},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   false,
	EndSectionWithAStatement:     false}

var Cron = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   false,
	EndSectionWithAStatement:     false}

var Hosts = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   false,
	EndSectionWithAStatement:     false}

var LoginDefs = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   false,
	EndSectionWithAStatement:     false}

var NsswitchConf = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{},
	BeginSectionWithAStatement:   false,
	EndSectionWithAStatement:     false}

var Httpd = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{"\"", "'"},
	SectionBeginningPrefixes:     []string{"<"},
	SectionBeginningSuffixes:     []string{">"},
	SectionEndingPrefixes:        []string{"</"},
	SectionEndingSuffixes:        []string{">"},
	BeginSectionWithAStatement:   true,
	EndSectionWithAStatement:     true}

var Named = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{"\\"},
	StatementEndingMarkers:       []string{";"},
	CommentBeginningMarkers:      []string{"#"},
	TextQuoteStyle:               []string{"\"", "'"},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{"{"},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{"};"},
	BeginSectionWithAStatement:   true,
	EndSectionWithAStatement:     false}

var NamedZone = txtedit.AnalyserConfig{
	StatementContinuationMarkers: []string{},
	StatementEndingMarkers:       []string{"\n"},
	CommentBeginningMarkers:      []string{";"},
	TextQuoteStyle:               []string{},
	SectionBeginningPrefixes:     []string{},
	SectionBeginningSuffixes:     []string{"("},
	SectionEndingPrefixes:        []string{},
	SectionEndingSuffixes:        []string{");"},
	BeginSectionWithAStatement:   true,
	EndSectionWithAStatement:     false}
