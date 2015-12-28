package txtedit

const (
	SECTION_MATCH_NO_SECTION           = 0  // The document does not use any section marker
	SECTION_MATCH_FLAT_SINGLE_ANCHOR   = 11 // For example ==Foobar
	SECTION_MATCH_FLAT_DOUBLE_ANCHOR   = 12 // For example [Foobar]
	SECTION_MATCH_NESTED_DOUBLE_ANCHOR = 22 // For example Foo{bar}
	SECTION_MATCH_NESTED_QUAD_ANCHOR   = 24 // For example <Foo>bar</Foo>
)

type SectionMatchMechanism int // Influence how section beginning/ending are detected

// Describe the writing style of the document so that analyser can break it down correctly.
type AnalyserConfig struct {
	StatementContinuationMarkers []string
	StatementEndingMarkers       []string
	CommentBeginningMarkers      []string
	TextQuoteStyle               []string
	BeginSectionWithAStatement   bool
	EndSectionWithAStatement     bool

	SectionBeginningPrefixes []string
	SectionBeginningSuffixes []string
	SectionEndingPrefixes    []string
	SectionEndingSuffixes    []string
	AmbiguousSectionSuffix   bool
	SectionMatchMechanism    SectionMatchMechanism
}

// Determine the mechanism for detecting section's beginning/ending.
func (cfg *AnalyserConfig) DetectSectionMatchMechanism() {
	if len(cfg.SectionBeginningPrefixes) > 0 && len(cfg.SectionBeginningSuffixes) > 0 &&
		len(cfg.SectionEndingSuffixes) > 0 && len(cfg.SectionEndingPrefixes) > 0 {
		// All markers are present, sections can be nested. For example <Foo>bar</Foo>
		cfg.SectionMatchMechanism = SECTION_MATCH_NESTED_QUAD_ANCHOR
	} else if len(cfg.SectionBeginningSuffixes) > 0 && len(cfg.SectionEndingPrefixes) > 0 {
		// Two markers surround the section, sections can be nested. For example Foo{bar}
		cfg.SectionMatchMechanism = SECTION_MATCH_NESTED_DOUBLE_ANCHOR
	} else if len(cfg.SectionBeginningPrefixes) > 0 && len(cfg.SectionBeginningSuffixes) > 0 {
		// Two markers surround the section title, sections do not nest. For example [Foobar]
		cfg.SectionMatchMechanism = SECTION_MATCH_FLAT_DOUBLE_ANCHOR
	} else if len(cfg.SectionBeginningPrefixes) > 0 {
		// Single marker marks beginning of a section, sections do not nest. For example ==Foobar
		cfg.SectionMatchMechanism = SECTION_MATCH_FLAT_SINGLE_ANCHOR
	} else {
		// Document does not use any section marker
		cfg.SectionMatchMechanism = SECTION_MATCH_NO_SECTION
	}
	/*
		Ambiguous section suffixes require special treatment in the analyser according to the bool flag.
		This is an example of using ambiguous suffixes: <Foo>bar</Foo>
		This is an example of non-ambiguous suffixes: <Foo>bar</Foo}
	*/
	for _, style1 := range cfg.SectionBeginningSuffixes {
		for _, style2 := range cfg.SectionEndingSuffixes {
			if style1 == style2 {
				cfg.AmbiguousSectionSuffix = true
			}
		}
	}
}
