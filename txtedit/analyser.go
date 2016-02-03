package txtedit

/*
Text analyser analyses input text character by character, breaks down the whole document into smaller
pieces that are easier for further analysis and reproduction of document text.
*/
type Analyser struct {
	textInput string           // the original input text
	config    *AnalyserConfig  // document style specification and more configuration
	debug     AnalyzerDebugger // output from analyser's progress, and output of debug information

	previousMarkerPosition int           // the character index where previous marker was encountered
	herePosition           int           // index of the current character where analyser has progressed
	rootNode               *DocumentNode // the root node of the broken down document
	thisNode               *DocumentNode // reference to the current document node

	ignoreNewStatementOnce bool       // do not create the next new statement caused by statement continuation marker
	contextText            *Text      // reference to the current text entity
	contextComment         *Comment   // reference to the current comment entity
	contextStatement       *Statement // reference to the current statement
}

// Initialise a new text analyser.
func NewAnalyser(textInput string, config *AnalyserConfig, debugger AnalyzerDebugger) (ret *Analyser) {
	ret = &Analyser{textInput: textInput, config: config, debug: debugger}
	ret.thisNode = &DocumentNode{Parent: nil, Obj: nil, Leaves: make([]*DocumentNode, 0, 8)}
	ret.rootNode = ret.thisNode
	ret.config.SectionStyle.SetSectionMatchMechanism()
	ret.debug.Printf("NewAnalyser: initialised with section match mechanism being %v", ret.config.SectionStyle.SectionMatchMechanism)
	return
}

// Create a new sibling node if the current node is already holding an object. Move reference to the new sibling.
func (an *Analyser) createSiblingNodeIfNotNil() {
	if an.thisNode.Obj == nil {
		an.debug.Printf("createSiblingNodeIfNotNil: does nothing when this node %p is still empty", an.thisNode)
		return
	}
	if parent := an.thisNode.Parent; parent == nil {
		/*
			In case thisNode is the root node, create a new root node and make both original
			root node and new sibling node leaves.
		*/
		originalRoot := an.rootNode
		newRoot := &DocumentNode{Parent: nil, Leaves: make([]*DocumentNode, 0, 8)}
		newRoot.Leaves = append(newRoot.Leaves, originalRoot)
		newLeaf := &DocumentNode{Parent: newRoot, Leaves: make([]*DocumentNode, 0, 8)}
		newRoot.Leaves = append(newRoot.Leaves, newLeaf)
		an.rootNode = newRoot
		an.thisNode = newLeaf
		an.debug.Printf("createSiblingNodeIfNotNil: new root is %p, original root %p is now a leaf, new sibling is %p",
			an.rootNode, originalRoot, newLeaf)
	} else {
		newLeaf := &DocumentNode{Parent: parent, Leaves: make([]*DocumentNode, 0, 8)}
		parent.Leaves = append(parent.Leaves, newLeaf)
		an.thisNode = newLeaf
		an.debug.Printf("createSiblingNodeIfNotNil: new sibling is %p", newLeaf)
	}
}

/*
If the current node already holds an object, then create a new sibling.
Save an object in the current node.
*/
func (an *Analyser) createDocumentSiblingNode(nodeContent interface{}) {
	if nodeContent == nil {
		an.debug.Printf("createDocumentSiblingNode: does nothing when node content is nil")
		return
	}
	if an.thisNode.Obj == nil {
		an.thisNode.Obj = nodeContent
	} else {
		// Must not overwrite the object in thisNode
		an.createSiblingNodeIfNotNil()
		an.debug.Printf("createDocumentSiblingNode: store object %p in node %p", nodeContent, an.thisNode)
		an.thisNode.Obj = nodeContent
	}
}

// Create a new leaf node and move reference to the new leaf.
func (an *Analyser) createLeaf() {
	newLeaf := &DocumentNode{Parent: an.thisNode, Leaves: make([]*DocumentNode, 0, 8)}
	an.thisNode.Leaves = append(an.thisNode.Leaves, newLeaf)
	an.debug.Printf("createLeaf: %p now has a new leaf %p", an.thisNode, newLeaf)
	an.thisNode = newLeaf
}

// If comment context is nil, assign the context a new comment entity.
func (an *Analyser) createCommentIfNil(commentStyle CommentStyle) {
	if an.contextComment == nil {
		an.contextComment = new(Comment)
		an.contextComment.CommentStyle = commentStyle
		an.debug.Printf("createCommentIfNil: context comment is assigned to %p", an.contextComment)
	} else {
		an.debug.Printf("createCommentIfNil: comment style goes into %p", an.contextComment)
		an.saveMissedCharacters()
		//		an.contextComment.Content += commentStyle
	}
}

// If comment context is not nil, move the comment into statement context and clear comment context.
func (an *Analyser) endComment(closed bool) {
	if an.contextComment == nil {
		return
	}
	an.contextComment.Closed = closed
	an.saveMissedCharacters()
	an.createStatementIfNil()
	an.contextStatement.Pieces = append(an.contextStatement.Pieces, an.contextComment)
	an.debug.Printf("endComment: comment %p is now a piece of statement %p", an.contextComment, an.contextStatement)
	an.contextComment = nil
}

// If text context is nil, assign the context a new text entity.
func (an *Analyser) createTextIfNil() {
	if an.contextText == nil {
		an.contextText = new(Text)
		an.debug.Printf("createTextIfNil: context text is assigned to %p", an.contextText)
	}
}

// If text context is not nil, move the text into statement context and clear text context.
func (an *Analyser) endText() {
	if an.contextText == nil {
		return
	}
	an.saveMissedCharacters()
	an.createStatementIfNil()
	an.contextStatement.Pieces = append(an.contextStatement.Pieces, an.contextText)
	an.debug.Printf("endText: text %p is now a piece of statement %p", an.contextText, an.contextStatement)
	an.contextText = nil
}

// If statement context is nil, create a new statement as document node, and assign it to the context.
func (an *Analyser) createStatementIfNil() {
	if an.contextStatement == nil {
		an.contextStatement = new(Statement)
		an.createDocumentSiblingNode(an.contextStatement)
		an.debug.Printf("createStatementIfNil: context statement is assigned to %p", an.contextStatement)
	}
}

// Move context text and comment into context statement (create new statement if necessary), and clear context statement.
func (an *Analyser) endStatement(ending string) {
	an.debug.Printf("endStatement: trying to end with %v", []byte(ending))
	if an.contextText != nil && an.contextText.QuoteStyle != "" {
		an.saveMissedCharacters()
		an.debug.Printf("endStatement: the statement ending goes into context text %p", an.contextText)
		an.contextText.Text += ending
		return
	}
	// Organise context objects
	an.saveMissedCharacters()
	an.endComment(false)
	an.endText()
	if an.ignoreNewStatementOnce {
		an.debug.Printf("endStatement: not creating new document node when ignoreNewStatementOnce is set")
		an.ignoreNewStatementOnce = false
		return
	}
	if an.contextStatement == nil && an.contextComment == nil && an.contextText == nil {
		an.debug.Printf("endStatement: context comment and text are nil, nothing to save")
		if ending != "" {
			an.debug.Printf("endStatement: save statement ending in a new statement")
			// The following branch is a workaround for input text "[]\n", assuming [] denotes a section.
			if state, _ := an.getSectionState(); state == SECTION_STATE_END_NOW {
				an.endSection()
			}
			an.createDocumentSiblingNode(&Statement{Ending: ending})
		}
		return
	}
	// Save the remaining text/comment piece
	an.contextStatement.Ending = ending
	if an.contextComment != nil {
		an.contextStatement.Pieces = append(an.contextStatement.Pieces, an.contextComment)
	}
	if an.contextText != nil {
		an.contextStatement.Pieces = append(an.contextStatement.Pieces, an.contextText)
	}
	// Remember - statement context was placed in the document node tree when it was created
	an.contextComment = nil
	an.contextText = nil
	an.contextStatement = nil
}

/*
In the context comment or text, save the characters that have not yet been placed in any entity.
Return true only if the missed characters has been saved.
*/
func (an *Analyser) saveMissedCharacters() bool {
	if an.herePosition-an.previousMarkerPosition <= 0 {
		return false // nothing missed
	}
	missedContent := an.textInput[an.previousMarkerPosition:an.herePosition]
	if an.contextComment != nil {
		an.debug.Printf("saveMissedText: missed content '%s' is stored in comment %p",
			missedContent, an.contextComment)
		an.contextComment.Content += missedContent
	} else {
		an.createTextIfNil()
		an.debug.Printf("saveMissedText: missed content '%s' is stored in text %p",
			missedContent, an.contextText)
		an.contextText.Text += missedContent
	}
	an.previousMarkerPosition = an.herePosition
	return true
}

// Place the space characters inside statement indentation or text entity's trailing spaces.
func (an *Analyser) saveSpaces(spaces string) {
	an.saveMissedCharacters()
	length := len(spaces)
	if an.ignoreNewStatementOnce {
		an.endText()
		an.debug.Printf("saveSpaces: ignoreNewStatementOnce is true, %d spaces go into text %p", length, an.contextText)
		an.createTextIfNil()
		an.contextText.TrailingSpaces += spaces
		an.endText()
	} else if an.contextComment != nil {
		an.debug.Printf("saveSpaces: %d spaces go into comment %p", length, an.contextComment)
		an.contextComment.Content += spaces
	} else if an.contextText != nil {
		an.debug.Printf("saveSpaces: %d spaces go into text %p", length, an.contextText)
		an.contextText.TrailingSpaces = spaces
		an.endText()
	} else if an.contextStatement != nil && len(an.contextStatement.Pieces) > 0 {
		lastPiece := an.contextStatement.Pieces[len(an.contextStatement.Pieces)-1]
		switch t := lastPiece.(type) {
		case *Text:
			t.TrailingSpaces += spaces
			an.debug.Printf("saveSpaces: %d spaces go into last text piece %p", t)
		case *Comment:
			t.Content += spaces
			an.debug.Printf("saveSpaces: %d spaces go into last comment piece %p", t)
		case *StatementContinue:
			an.createTextIfNil()
			an.debug.Printf("saveSpaces: %d spaces go into new text piece %p", an.contextText)
			an.contextText.TrailingSpaces += spaces
			an.endText()
		}
	} else if an.contextStatement != nil {
		an.debug.Printf("saveSpaces: %d spaces go into indentation of context statement %p",
			length, an.contextStatement)
		an.contextStatement.Indent += spaces
	} else if an.contextStatement == nil {
		an.createStatementIfNil()
		an.debug.Printf("saveSpaces: %d spaces go into indentation of a new statement %p",
			length, an.contextStatement)
		an.contextStatement.Indent += spaces
	} else {
		an.debug.Printf("saveSpaces: %d spaces have nowhere to go")
	}
}

// In the context comment or quoted text, save the characters. Return true only if such context is found.
func (an *Analyser) saveQuoteOrCommentCharacters(str string) bool {
	if an.contextComment != nil {
		an.saveMissedCharacters()
		an.debug.Printf("saveQuoteOrCommentCharacters: save '%s' in context comment %p", str, an.contextComment)
		an.contextComment.Content += str
		return true
	}
	if an.contextText != nil && an.contextText.QuoteStyle != "" {
		an.saveMissedCharacters()
		an.debug.Printf("saveQuoteOrCommentCharacters: save '%s' in context text %p", str, an.contextText)
		an.contextText.Text += str
		return true
	}
	return false
}

// Save missed text and prevent the next new statement from being created.
func (an *Analyser) continueStatement(marker string) {
	if an.saveQuoteOrCommentCharacters(marker) {
		an.debug.Printf("continueStatement: the marker went to saveQuoteOrCommentCharacters")
		return
	}
	an.saveMissedCharacters()
	an.endComment(false)
	an.endText()
	an.createStatementIfNil()
	an.contextStatement.Pieces = append(an.contextStatement.Pieces, &StatementContinue{Style: marker})
	an.ignoreNewStatementOnce = true
	an.debug.Printf("continueStatement: ignoreNewStatementOnce is set to true")
}

// Create a new section as document node, then shift thisNode to a leaf.
func (an *Analyser) createSection() {
	an.endStatement("")
	newSection := new(Section)
	if an.thisNode == an.rootNode {
		an.debug.Printf("createSection: root node %p has the new section %p", an.thisNode, newSection)
		an.createLeaf()
		an.thisNode.Obj = newSection
	} else {
		an.debug.Printf("newSection: node %p has the new section %p", an.thisNode, newSection)
		an.createDocumentSiblingNode(newSection)
	}
	an.createLeaf()
}

/*
Look for a Statement in the previous sibling and remove the sibling node.
Return the Statement if it is found, return nil if not found.
*/
func (an *Analyser) removePreviousSiblingStatement() *Statement {
	index := an.thisNode.GetMyLeafIndex()
	if index == -1 {
		an.debug.Printf("removePreviousSiblingStatement: cannot find node %p's leaf index",
			an.thisNode)
		return nil
	} else if index == 0 {
		an.debug.Printf("removePreviousSiblingStatement: this node %p does not have a previous sibling",
			an.thisNode)
		return nil
	}
	// Look for a statement in the previous sibling
	previousSibling := an.thisNode.Parent.Leaves[index-1]
	if obj := previousSibling.Obj; obj == nil {
		an.debug.Printf("removePreviousSiblingStatement: this node %p's previous leaf %p is empty",
			an.thisNode, previousSibling)
		return nil
	} else if stmt, ok := obj.(*Statement); !ok {
		an.debug.Printf("removePreviousSiblingStatement: this node %p's previous leaf %p does not hold a statement",
			an.thisNode, previousSibling)
		return nil
	} else {
		// Remove the sibling and return the statement
		an.debug.Printf("removePreviousSiblingStatement: this node %p's previous leaf %p holds a statement %p",
			an.thisNode, previousSibling, stmt)
		leaves := an.thisNode.Parent.Leaves
		an.thisNode.Parent.Leaves = append(leaves[0:index-1], leaves[index:]...)
		return stmt
	}
}

// Assign the section its first and final statements if necessary, then move thisNode to a new sibling of its parent's.
func (an *Analyser) endSection() {
	if _, sect := an.getSectionState(); sect == nil {
		an.debug.Printf("endSection: this node %p is not in a section", an.thisNode)
	} else {
		an.endStatement("")
		// Move thisNode to its parent, the one holding *Section
		an.thisNode = an.thisNode.Parent
		an.debug.Printf("endSection: trying to finish section in node %p", an.thisNode)
		// Calculate the first and final statements
		minNumLeaves := 0
		if an.config.SectionStyle.BeginSectionWithAStatement {
			/*
				If prefix styles are empty, the text should look like:
				sectionA {
				content
				}
				Therefore the section's first statement is the section node's previous sibling.

				If prefix styles are present, the text looks like:
				<sectionA>
				content
				</sectionA>
				Then the first statement is the section node's first leaf.
			*/
			if an.config.SectionStyle.OpeningPrefix == "" {
				sect.FirstStatement = an.removePreviousSiblingStatement()
				an.debug.Printf("endSection: first statement is the previous sibling %p", sect.FirstStatement)
			} else {
				sectionFirstLeaf := an.thisNode.Leaves[0]
				if leafObj := sectionFirstLeaf.Obj; leafObj == nil {
					an.debug.Printf("endSection: first statement should be the first leaf %p but it holds nothing",
						sectionFirstLeaf)
				} else if stmt, ok := sectionFirstLeaf.Obj.(*Statement); ok {
					an.thisNode.Leaves = an.thisNode.Leaves[1:]
					sect.FirstStatement = stmt
					minNumLeaves++
					an.debug.Printf("endSection: first statement is the first leaf %p's content, statement %p",
						sectionFirstLeaf, stmt)
				} else {
					an.debug.Printf("endSection: first statement should be the first leaf %p but it does not hold a statement",
						sectionFirstLeaf)
				}
			}
		}
		if an.config.SectionStyle.EndSectionWithAStatement {
			/*
				Rather than two scenarios supported by BeginSectionWithAStatement, there is only one scenario
				to deal with here, the section ending must use both prefixes and suffixes, like this:
				<sectionA>
				content
				</sectionA>     <=== "sectionA" is the ending statement
			*/
			if an.config.SectionStyle.ClosingPrefix != "" && an.config.SectionStyle.ClosingSuffix != "" {
				// minNumLeaves is 0 if section does not begin with statement that is also the section's leaf
				// minNumLeaves is 1 if section begins with a statement that is the section's leaf
				if len(an.thisNode.Leaves) > minNumLeaves {
					lastLeaf := an.thisNode.Leaves[len(an.thisNode.Leaves)-1]
					if lastLeaf == nil {
						an.debug.Printf("endSection: cannot assign final statement, the last leaf is nil.")
					} else if stmt, ok := lastLeaf.Obj.(*Statement); ok {
						an.thisNode.Leaves = an.thisNode.Leaves[0 : len(an.thisNode.Leaves)-1]
						sect.FinalStatement = stmt
						an.debug.Printf("endSection: final statement is the last leaf %p's content, statement %p",
							lastLeaf, stmt)
					} else {
						an.debug.Printf("endSection: final statement should be the last leaf %p but it does not hold a statement",
							lastLeaf)
					}
				} else {
					an.debug.Printf("endSection: cannot end section with a statement, there are not enough leaves.")
				}
			} else {
				an.debug.Printf("endSection: the config should specify both prefix and suffix in order to end a section with a statement")
			}
		}
		an.createSiblingNodeIfNotNil()
	}
	// Remember - section object was placed in the document node tree when it was created
}

const (
	SECTION_STATE_BEFORE_BEGIN     = 0
	SECTION_STATE_HAS_BEGIN_PREFIX = 1000
	SECTION_STATE_HAS_BEGIN_SUFFIX = 1100
	SECTION_STATE_HAS_END_PREFIX   = 1110
	SECTION_STATE_END_NOW          = 1111
)

type SectionState int

func (an *Analyser) getSectionState() (SectionState, *Section) {
	// If thisNode is in a section, the parent node should hold a *Section.
	if an.thisNode.Parent == nil {
		an.debug.Printf("getSectionState: node %p's parent is nil", an.thisNode)
		return SECTION_STATE_BEFORE_BEGIN, nil
	} else if an.thisNode.Parent.Obj == nil {
		an.debug.Printf("getSectionState: node %p's parent %p is empty", an.thisNode, an.thisNode.Parent)
		return SECTION_STATE_BEFORE_BEGIN, nil
	}
	section, isSect := an.thisNode.Parent.Obj.(*Section)
	if !isSect {
		an.debug.Printf("getSectionState: node %p's parent %p holds a %+v, which is not a *Section",
			an.thisNode, an.thisNode.Parent, an.thisNode.Parent.Obj)
		return SECTION_STATE_BEFORE_BEGIN, nil
	}
	var state SectionState
	switch an.config.SectionStyle.SectionMatchMechanism {
	case SECTION_MATCH_FLAT_SINGLE_ANCHOR:
		if section.BeginPrefix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else {
			state = SECTION_STATE_END_NOW
		}
	case SECTION_MATCH_FLAT_DOUBLE_ANCHOR:
		if section.BeginPrefix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else if section.BeginSuffix == "" {
			state = SECTION_STATE_HAS_BEGIN_PREFIX
		} else {
			state = SECTION_STATE_END_NOW
		}
	case SECTION_MATCH_NESTED_DOUBLE_ANCHOR:
		if section.BeginSuffix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else if section.EndSuffix == "" {
			state = SECTION_STATE_HAS_END_PREFIX
		} else {
			state = SECTION_STATE_END_NOW
		}
	case SECTION_MATCH_NESTED_QUAD_ANCHOR:
		if section.BeginPrefix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else if section.BeginSuffix == "" {
			state = SECTION_STATE_HAS_BEGIN_PREFIX
		} else if section.EndPrefix == "" {
			state = SECTION_STATE_HAS_BEGIN_SUFFIX
		} else if section.EndSuffix == "" {
			state = SECTION_STATE_HAS_END_PREFIX
		} else {
			state = SECTION_STATE_END_NOW
		}
	default:
		an.debug.Printf("getSectionState: unknown SectionMatchMechanism")
	}
	an.debug.Printf("getSectionState: state is %v, section is %p", state, section)
	return state, section
}

// Save the section beginning's prefix marking, create a new section if there is not one.
func (an *Analyser) setSectionBeginPrefix(prefix string) {
	if an.saveQuoteOrCommentCharacters(prefix) {
		return
	}
	an.endStatement("")
	state, sect := an.getSectionState()
	if state == SECTION_STATE_BEFORE_BEGIN {
		// Create a new section while not being in a section
		an.debug.Printf("setSectionBeginPrefix: create a new section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = prefix
	} else if an.config.SectionStyle.SectionMatchMechanism == SECTION_MATCH_FLAT_DOUBLE_ANCHOR ||
		an.config.SectionStyle.SectionMatchMechanism == SECTION_MATCH_FLAT_DOUBLE_ANCHOR {
		// Already in a section but document does not allow nested section
		an.debug.Printf("setSectionBeginPrefix: end section of node %p and create a new section", an.thisNode.Parent)
		an.endSection()
		an.createSection()
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = prefix
	} else if state == SECTION_STATE_END_NOW {
		// Marker matches but section should end now
		an.debug.Printf("setSectionBeginPrefix: end section right now")
		sect.BeginPrefix = prefix
		an.endSection()
	} else {
		// Create a nested section
		an.debug.Printf("setSectionBeginPrefix: create a nested section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = prefix
	}
}

// Save the section beginning's suffix marking.
func (an *Analyser) setSectionBeginSuffix(suffix string) {
	if an.saveQuoteOrCommentCharacters(suffix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state == SECTION_STATE_END_NOW {
		// Marker matches but section should end now
		an.debug.Printf("setSectionBeginSuffix: end section right now")
		sect.BeginSuffix = suffix
		an.endSection()
	} else if an.config.SectionStyle.SectionMatchMechanism == SECTION_MATCH_NESTED_DOUBLE_ANCHOR {
		// Create a section or nested section
		an.debug.Printf("setSectionBeginSuffix: create a new section/nested section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Obj.(*Section).BeginSuffix = suffix
	} else if state < SECTION_STATE_HAS_BEGIN_PREFIX || state > SECTION_STATE_HAS_BEGIN_SUFFIX {
		// State is not right so the marker must have been text
		an.debug.Printf("setSectionBeginSuffix: state is not right so only store the characters")
		an.saveMissedCharacters()
	} else {
		// Set suffix if state is right
		an.debug.Printf("setSectionBeginSuffix: set suffix")
		sect.BeginSuffix = suffix
	}
}

// Save the section beginning's prefix marking.
func (an *Analyser) setSectionEndPrefix(prefix string) {
	if an.saveQuoteOrCommentCharacters(prefix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state == SECTION_STATE_END_NOW {
		an.debug.Printf("setSectionEndPrefix: end section right now")
		sect.EndPrefix = prefix
		an.endSection()
	} else if state < SECTION_STATE_HAS_BEGIN_SUFFIX || state > SECTION_STATE_HAS_END_PREFIX {
		an.debug.Printf("setSectionEndPrefix: state is not right so only store the characters")
		an.saveMissedCharacters()
	} else {
		an.debug.Printf("setSectionEndPrefix: set prefix")
		sect.EndPrefix = prefix
	}
}

// Save the section ending's suffix marking.
func (an *Analyser) setSectionEndSuffix(suffix string) {
	if an.saveQuoteOrCommentCharacters(suffix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state >= SECTION_STATE_HAS_END_PREFIX {
		an.debug.Printf("setSectionEndSuffix: end section right now")
		sect.EndSuffix = suffix
		an.endSection()
	} else if state < SECTION_STATE_HAS_END_PREFIX && an.config.SectionStyle.AmbiguousSectionSuffix {
		an.debug.Printf("setSectionEndSuffix: call setSectionSetBeginSuffix due to ambiguous suffix choice")
		an.setSectionBeginSuffix(suffix)
	} else {
		an.debug.Printf("EndSectionSetSuffix: set suffix")
		sect.EndSuffix = suffix
	}
}

// Look for the string from position here. Return the matching string and length of the match.
func (an *Analyser) lookFor(match string) (string, int) {
	if len(match) == 1 {
		// Match single character
		if an.textInput[an.herePosition] == match[0] {
			return match, 1
		}
		return "", 0
	} else {
		// Match string more than two characters long
		if an.herePosition+len(match) > len(an.textInput) {
			return "", 0
		} else if string(an.textInput[an.herePosition:an.herePosition+len(match)]) != match {
			return "", 0
		}
		return match, len(match)
	}
}

// Look for any string among the match list, from position here. Return the matching string and length of the match.
func (an *Analyser) lookForAnyOf(matches []string) (string, int) {
	for _, match := range matches {
		if len(match) == 1 {
			// Match single character
			if an.textInput[an.herePosition] == match[0] {
				return match, 1
			} else {
				continue
			}
		} else {
			// Match string more than two characters long
			if an.herePosition+len(match) > len(an.textInput) {
				continue
			} else if string(an.textInput[an.herePosition:an.herePosition+len(match)]) != match {
				continue
			} else {
				return match, len(match)
			}
		}
	}
	return "", 0
}

/*
Look for consecutive spaces from here position. Return the string of consecutive spaces and its length.
Space characters are ' ' and '\t'.
*/
func (an *Analyser) lookForSpaces() (string, int) {
	pos := an.herePosition
	for ; pos < len(an.textInput); pos++ {
		if an.textInput[pos] != ' ' && an.textInput[pos] != '\t' {
			break
		}
	}
	return an.textInput[an.herePosition:pos], pos - an.herePosition
}

// Toggle text quoting in the analyser' context.
func (an *Analyser) setQuote(quoteStyle string) {
	if an.contextComment != nil {
		an.saveMissedCharacters()
		an.debug.Printf("setQuote: quote '%s' goes into context comment", quoteStyle)
		an.contextComment.Content += quoteStyle
		return
	}
	// Save missed text that is not being quoted
	if (an.contextText == nil || an.contextText.QuoteStyle == "") && an.saveMissedCharacters() {
		// Let the quote mark go into a new text entity
		an.endText()
	}
	an.createTextIfNil()
	if an.contextText.QuoteStyle == "" {
		an.debug.Printf("setQuote: begin quoting in text %p", an.contextText)
		an.contextText.QuoteStyle = quoteStyle
	} else {
		if an.contextText.QuoteStyle == quoteStyle {
			an.debug.Printf("setQuote: finish quoting in text %p", an.contextText)
			an.endText()
		} else {
			an.debug.Printf("setQuote: quote '%s' goes into context text %p", an.contextText)
			an.saveMissedCharacters()
			an.contextText.Text += quoteStyle
		}
	}
}

// Tell the analyser to open a comment if the text at position here matches any comment opening style.
func (an *Analyser) isOpeningComment() int {
	if an.contextComment != nil {
		// A comment is already open, so it is not possible to open another comment.
		return 0
	}
	for _, style := range an.config.CommentStyles {
		if match, advance := an.lookFor(style.Opening); advance > 0 {
			an.debug.Printf("Comment opening: %s", match)
			an.createCommentIfNil(style)
			return advance
		}
	}
	return 0
}

// Tell the analyser to close a comment if the text at position here matches any comment closing style.
func (an *Analyser) isClosingComment() int {
	if an.contextComment == nil {
		// Comment has not been opened, so it is not possible to close a comment.
		return 0
	}
	for _, style := range an.config.CommentStyles {
		if match, advance := an.lookFor(style.Closing); advance > 0 {
			if match == an.contextComment.CommentStyle.Closing {
				an.debug.Printf("Comment closing: %s", match)
				an.endComment(true)
				return advance
			}
		}
	}
	return 0
}

// Break down input text according to analyser's configuration. Return the root document node.
func (an *Analyser) Run() *DocumentNode {
	/*
		The loop visits the input text character by character, which sets "advance" to 1; unless it meets
		a marker, which can be longer than one character, and "advance" will be the marker string's length.
		A marker triggers special processing logic, such as toggling text quote, starting a section, etc.
		The previousMarkerPosition is updated with the every marker along the way.
	*/
	var advance int // how many characters to advance for the next iteration
	for an.herePosition = 0; an.herePosition < len(an.textInput); an.herePosition += advance {
		var match string  // the marker string immediate ahead
		var spaces string // number of consecutive spaces immediate ahead
		if advance = an.isOpeningComment(); advance > 0 {
			an.previousMarkerPosition = an.herePosition + advance
		} else if advance = an.isClosingComment(); advance > 0 {
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookForAnyOf(an.config.TextQuoteStyle); advance > 0 {
			an.debug.Printf("Quote: %s", match)
			an.setQuote(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if spaces, advance = an.lookForSpaces(); advance > 0 {
			an.debug.Printf("Spaces: length %d", advance)
			an.saveSpaces(spaces)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookForAnyOf(an.config.StatementContinuationMarkers); advance > 0 {
			an.debug.Printf("Statement continuation: %s", match)
			an.continueStatement(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookForAnyOf(an.config.StatementEndingMarkers); advance > 0 {
			an.debug.Printf("Statement ending: %v", []byte(match))
			an.endStatement(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.ClosingSuffix); advance > 0 {
			an.debug.Printf("Section closing suffix: %s", match)
			an.setSectionEndSuffix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.ClosingPrefix); advance > 0 {
			an.debug.Printf("Section closing prefix: %s", match)
			an.setSectionEndPrefix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.OpeningSuffix); advance > 0 {
			an.debug.Printf("Section opening suffix: %s", match)
			an.setSectionBeginSuffix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.OpeningPrefix); advance > 0 {
			an.debug.Printf("Section opening prefix: %s", match)
			an.setSectionBeginPrefix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else {
			advance = 1
		}
	}
	an.debug.Printf("Run: end statement for the last time")
	an.endStatement("")
	an.debug.Printf("Run: end all open sections for the last time")
	// End all sections
	for recursionGuard := 0; recursionGuard < 100 && an.thisNode.Parent != nil && an.thisNode.Parent != an.rootNode; recursionGuard++ {
		an.endSection()
	}
	return an.rootNode
}
