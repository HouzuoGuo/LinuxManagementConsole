package txtedit

const (
	SECTION_MATCH_FLAT_SINGLE_ANCHOR   = 11 // For example ==Foobar
	SECTION_MATCH_FLAT_DOUBLE_ANCHOR   = 12 // For example [Foobar]
	SECTION_MATCH_NESTED_DOUBLE_ANCHOR = 22 // For example Foo{bar}
	SECTION_MATCH_NESTED_QUAD_ANCHOR   = 24 // For example <Foo>bar</Foo>
)

type SectionMatchMechanism int // Influence how section beginning/ending are detected

// Describe how sections are opened/closed.
type SectionStyle struct {
	OpeningPrefix, OpeningSuffix string
	ClosingPrefix, ClosingSuffix string
	OpenSectionWithAStatement    bool
	CloseSectionWithAStatement   bool
	AmbiguousSectionSuffix       bool
	SectionMatchMechanism        SectionMatchMechanism
}

// Determine the mechanism for detecting section's opening and closing.
func (style *SectionStyle) SetSectionMatchMechanism() {
	if style.OpeningPrefix != "" && style.OpeningSuffix != "" &&
		style.ClosingPrefix != "" && style.ClosingSuffix != "" {
		// All markers are present, sections can be nested. For example <Foo>bar</Foo>
		style.SectionMatchMechanism = SECTION_MATCH_NESTED_QUAD_ANCHOR
	} else if style.OpeningSuffix != "" && style.ClosingSuffix != "" {
		// Two markers surround the section, sections can be nested. For example Foo{bar}
		style.SectionMatchMechanism = SECTION_MATCH_NESTED_DOUBLE_ANCHOR
	} else if style.OpeningPrefix != "" && style.OpeningSuffix != "" {
		// Two markers surround the section title, sections do not nest. For example [Foobar]
		style.SectionMatchMechanism = SECTION_MATCH_FLAT_DOUBLE_ANCHOR
	} else if style.OpeningPrefix != "" {
		// Single marker marks beginning of a section, sections do not nest. For example ==Foobar
		style.SectionMatchMechanism = SECTION_MATCH_FLAT_SINGLE_ANCHOR
	}
	/*
		Ambiguous section suffixes require special treatment in the analyser according to the bool flag.
		This is an example of using ambiguous suffixes: <Foo>bar</Foo>
		This is an example of non-ambiguous suffixes: <Foo>bar</Foo}
	*/
	if style.OpeningSuffix == style.ClosingSuffix {
		style.AmbiguousSectionSuffix = true
	}
}

// Describe how comments are written.
type CommentStyle struct {
	Opening, Closing string
}

// Describe the writing style of the document so that analyser can break it down correctly.
type AnalyserConfig struct {
	StatementContinuationMarkers []string
	StatementEndingMarkers       []string
	CommentStyles                []CommentStyle
	TextQuoteStyle               []string

	SectionStyle SectionStyle
}
