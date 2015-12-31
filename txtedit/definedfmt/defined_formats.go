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
