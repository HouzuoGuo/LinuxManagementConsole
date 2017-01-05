package lexer

/*
The lexer analyses input text character by character, breaks down the whole document into smaller
pieces that are easier for further analysis and reproduction of document text.
*/
type Lexer struct {
	textInput string       // the original input text
	config    *LexerConfig // document style specification and more configuration
	debug     LexerDebug   // handle output from lexer's progress and debug information

	previousMarkerPosition int           // the character index where previous marker was encountered
	herePosition           int           // index of the current character where lexer has progressed
	rootNode               *DocumentNode // the root node of the broken down document
	thisNode               *DocumentNode // reference to the current document node

	ignoreNewStatementOnce bool       // do not create the next new statement caused by statement continuation marker
	contextText            *Text      // reference to the current text entity
	contextComment         *Comment   // reference to the current comment entity
	contextStatement       *Statement // reference to the current statement

	statementCounter int // total number of statements that have been ended
}

// Initialise a new text lexer.
func NewLexer(textInput string, config *LexerConfig, debugger LexerDebug) (ret *Lexer) {
	ret = &Lexer{textInput: textInput, config: config, debug: debugger}
	ret.thisNode = &DocumentNode{Parent: nil, Entity: nil, Leaves: make([]*DocumentNode, 0, 8)}
	ret.rootNode = ret.thisNode
	ret.config.SectionStyle.SetSectionMatchMechanism()
	ret.debug.Printfln("NewLexer: initialised with section match mechanism being %v", ret.config.SectionStyle.SectionMatchMechanism)
	return
}

// Create a new sibling node if the current node is already holding an object. Move reference to the new sibling.
func (an *Lexer) createSiblingNodeIfNotNil() {
	if an.thisNode.Entity == nil {
		an.debug.Printfln("createSiblingNodeIfNotNil: does nothing when this node %p is still empty", an.thisNode)
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
		an.debug.Printfln("createSiblingNodeIfNotNil: new root is %p, original root %p is now a leaf, new sibling is %p",
			an.rootNode, originalRoot, newLeaf)
	} else {
		newLeaf := &DocumentNode{Parent: parent, Leaves: make([]*DocumentNode, 0, 8)}
		parent.Leaves = append(parent.Leaves, newLeaf)
		an.thisNode = newLeaf
		an.debug.Printfln("createSiblingNodeIfNotNil: new sibling is %p", newLeaf)
	}
}

/*
If the current node already holds an object, then create a new sibling.
Save an object in the current node.
*/
func (an *Lexer) createDocumentSiblingNode(nodeContent interface{}) {
	if nodeContent == nil {
		an.debug.Printfln("createDocumentSiblingNode: does nothing when node content is nil")
		return
	}
	if an.thisNode.Entity == nil {
		an.thisNode.Entity = nodeContent
	} else {
		// Must not overwrite the object in thisNode
		an.createSiblingNodeIfNotNil()
		an.debug.Printfln("createDocumentSiblingNode: store object %p in node %p", nodeContent, an.thisNode)
		an.thisNode.Entity = nodeContent
	}
}

// Create a new leaf node and move reference to the new leaf.
func (an *Lexer) createLeaf() {
	newLeaf := &DocumentNode{Parent: an.thisNode, Leaves: make([]*DocumentNode, 0, 8)}
	an.thisNode.Leaves = append(an.thisNode.Leaves, newLeaf)
	an.debug.Printfln("createLeaf: %p now has a new leaf %p", an.thisNode, newLeaf)
	an.thisNode = newLeaf
}

// If comment context is nil, assign the context a new comment entity.
func (an *Lexer) createCommentIfNil(commentStyle CommentStyle) {
	if an.contextComment == nil {
		an.contextComment = new(Comment)
		an.contextComment.CommentStyle = commentStyle
		an.debug.Printfln("createCommentIfNil: context comment is assigned to %p", an.contextComment)
	} else {
		an.debug.Printfln("createCommentIfNil: comment style goes into %p", an.contextComment)
		//		an.contextComment.Content += commentStyle
	}
}

// If comment context is not nil, move the comment into statement context and clear comment context.
func (an *Lexer) endComment(marker string, closed bool) {
	if an.contextComment == nil {
		return
	}
	an.contextComment.Closed = closed
	an.saveMissedCharacters()
	an.createStatementIfNil()
	an.contextStatement.Pieces = append(an.contextStatement.Pieces, an.contextComment)
	an.debug.Printfln("endComment: comment %p is now a piece of statement %p", an.contextComment, an.contextStatement)

	oldComment := an.contextComment
	an.contextComment = nil
	// Test if comment is also ending the statement
	if marker != "" {
		for _, stmtEndingMarker := range an.config.StatementEndingMarkers {
			if marker == stmtEndingMarker {
				/*
					The "closed" flag must be unset on the comment, otherwise when reproducing the original text,
					the ending marker will be reproduced twice - once by the comment, and once more by the statement.
				*/
				oldComment.Closed = false
				an.endStatement(marker)
				break
			}
		}
	}
}

// If text context is nil, assign the context a new text entity.
func (an *Lexer) createTextIfNil() {
	if an.contextText == nil {
		an.contextText = new(Text)
		an.debug.Printfln("createTextIfNil: context text is assigned to %p", an.contextText)
	}
}

// If text context is not nil, move the text into statement context and clear text context.
func (an *Lexer) endText() {
	if an.contextText == nil {
		return
	}
	an.saveMissedCharacters()
	an.createStatementIfNil()
	an.contextStatement.Pieces = append(an.contextStatement.Pieces, an.contextText)
	an.debug.Printfln("endText: text %p is now a piece of statement %p", an.contextText, an.contextStatement)
	an.contextText = nil
}

// If statement context is nil, create a new statement as document node, and assign it to the context.
func (an *Lexer) createStatementIfNil() {
	if an.contextStatement == nil {
		an.contextStatement = new(Statement)
		an.createDocumentSiblingNode(an.contextStatement)
		an.debug.Printfln("createStatementIfNil: context statement is assigned to %p", an.contextStatement)
	}
}

// Move context text and comment into context statement (create new statement if necessary), and clear context statement.
func (an *Lexer) endStatement(ending string) {
	an.debug.Printfln("endStatement: trying to end with %v", []byte(ending))
	if an.contextComment != nil && // if there is still a comment ...
		!an.contextComment.Closed && // that has not been closed ...
		an.contextComment.CommentStyle.Closing != ending && // and the statement ending does not close the comment
		ending != "" { // and this is not "ending statement no matter what" situation
		// If there is still a comment and the ending marker does not close the comment, only save the ending marker in the comment.
		an.saveMissedCharacters()
		an.debug.Printfln("endStatement: the statement ending goes into open context comment %p", an.contextComment)
		an.contextComment.Content += ending
		return
	}
	if an.contextText != nil && an.contextText.QuoteStyle != "" {
		an.saveMissedCharacters()
		an.debug.Printfln("endStatement: the statement ending goes into context text %p", an.contextText)
		an.contextText.Text += ending
		return
	}
	// Organise context objects
	an.saveMissedCharacters()
	an.endComment("", false)
	an.endText()
	if an.ignoreNewStatementOnce {
		an.debug.Printfln("endStatement: not creating new document node when ignoreNewStatementOnce is set")
		an.ignoreNewStatementOnce = false
		return
	}
	if an.contextStatement == nil && an.contextComment == nil && an.contextText == nil {
		an.debug.Printfln("endStatement: context comment and text are nil, nothing to save")
		if ending != "" {
			an.debug.Printfln("endStatement: save statement ending in a new statement")
			// The following branch is a workaround for input text "[]\n", assuming [] denotes a section.
			//			if state, _ := an.getSectionState(); state == SECTION_STATE_END_NOW {
			//				an.endSection()
			//			}
			// ^^^^^^^^^^ Why did I write that? ^^^^^^^^^^
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
	an.statementCounter++
}

/*
In the context comment or text, save the characters that have not yet been placed in any entity.
Return true only if the missed characters has been saved.
*/
func (an *Lexer) saveMissedCharacters() bool {
	if an.herePosition-an.previousMarkerPosition <= 0 {
		return false // nothing missed
	}
	missedContent := an.textInput[an.previousMarkerPosition:an.herePosition]
	if an.contextComment != nil {
		an.debug.Printfln("saveMissedText: missed content '%s' is stored in comment %p",
			missedContent, an.contextComment)
		an.contextComment.Content += missedContent
	} else {
		an.createTextIfNil()
		an.debug.Printfln("saveMissedText: missed content '%s' is stored in text %p",
			missedContent, an.contextText)
		an.contextText.Text += missedContent
	}
	an.previousMarkerPosition = an.herePosition
	return true
}

// Place the space characters inside statement indentation or text entity's trailing spaces.
func (an *Lexer) saveSpaces(spaces string) {
	an.saveMissedCharacters()
	length := len(spaces)
	if an.ignoreNewStatementOnce {
		an.endText()
		an.debug.Printfln("saveSpaces: ignoreNewStatementOnce is true, %d spaces go into text %p", length, an.contextText)
		an.createTextIfNil()
		an.contextText.TrailingSpaces += spaces
		an.endText()
	} else if an.contextComment != nil {
		an.debug.Printfln("saveSpaces: %d spaces go into comment %p", length, an.contextComment)
		an.contextComment.Content += spaces
	} else if an.contextText != nil {
		an.debug.Printfln("saveSpaces: %d spaces go into text %p", length, an.contextText)
		an.contextText.TrailingSpaces = spaces
		an.endText()
	} else if an.contextStatement != nil && len(an.contextStatement.Pieces) > 0 {
		lastPiece := an.contextStatement.Pieces[len(an.contextStatement.Pieces)-1]
		switch t := lastPiece.(type) {
		case *Text:
			t.TrailingSpaces += spaces
			an.debug.Printfln("saveSpaces: %d spaces go into last text piece %p", length, t)
		case *StatementContinue:
			an.createTextIfNil()
			an.debug.Printfln("saveSpaces: %d spaces go into new text piece %p", an.contextText)
			an.contextText.TrailingSpaces += spaces
			an.endText()
		case *Comment:
			if t.Closed {
				stmtClosedWithComment := false
				for _, stmtEndingMarker := range an.config.StatementEndingMarkers {
					if t.CommentStyle.Closing == stmtEndingMarker {
						/*
							In case the comment is closed along with the statement, the spaces should
							indent the next statement.
						*/
						stmtClosedWithComment = true
						an.endStatement("")
						an.createStatementIfNil()
						an.debug.Printfln("saveSpaces: %d spaces go into indentation of a new statement %p",
							length, an.contextStatement)
						an.contextStatement.Indent += spaces
						break
					}
				}
				if !stmtClosedWithComment {
					/*
						In case the comment is closed but the statement is not, the spaces go into
						a new text piece, and the text piece only holds the spaces.
					*/
					an.createTextIfNil()
					an.debug.Printfln("saveSpaces: %d spaces go into a new text piece %p",
						length, an.contextText)
					an.contextText.TrailingSpaces = spaces
					an.endText()
				}
			} else {
				t.Content += spaces
				an.debug.Printfln("saveSpaces: %d spaces go into last comment piece %p", length, t)
			}
		}
	} else if an.contextStatement != nil {
		an.debug.Printfln("saveSpaces: %d spaces go into indentation of context statement %p",
			length, an.contextStatement)
		an.contextStatement.Indent += spaces
	} else if an.contextStatement == nil {
		an.createStatementIfNil()
		an.debug.Printfln("saveSpaces: %d spaces go into indentation of a new statement %p",
			length, an.contextStatement)
		an.contextStatement.Indent += spaces
	} else {
		an.debug.Printfln("saveSpaces: %d spaces have nowhere to go")
	}
}

// In the context comment or quoted text, save the characters. Return true only if such context is found.
func (an *Lexer) saveQuoteOrCommentCharacters(str string) bool {
	if an.contextComment != nil {
		an.saveMissedCharacters()
		an.debug.Printfln("saveQuoteOrCommentCharacters: save '%s' in context comment %p", str, an.contextComment)
		an.contextComment.Content += str
		return true
	}
	if an.contextText != nil && an.contextText.QuoteStyle != "" {
		an.saveMissedCharacters()
		an.debug.Printfln("saveQuoteOrCommentCharacters: save '%s' in context text %p", str, an.contextText)
		an.contextText.Text += str
		return true
	}
	return false
}

// Immediately end and finish the current text, then save the marker into a new text piece.
func (an *Lexer) breakText(marker string) {
	if an.saveQuoteOrCommentCharacters(marker) {
		an.debug.Printfln("breakToken: the marker went to saveQuoteOrCommentCharacters")
		return
	}
	an.saveMissedCharacters()
	// Finish the current text
	an.endText()
	// Save the marker into its own text piece
	an.createTextIfNil()
	an.contextText.Text = marker
	an.endText()
}

// Save missed text and prevent the next new statement from being created.
func (an *Lexer) continueStatement(marker string) {
	if an.saveQuoteOrCommentCharacters(marker) {
		an.debug.Printfln("continueStatement: the marker went to saveQuoteOrCommentCharacters")
		return
	}
	an.saveMissedCharacters()
	an.endComment("", false)
	an.endText()
	an.createStatementIfNil()
	an.contextStatement.Pieces = append(an.contextStatement.Pieces, &StatementContinue{Style: marker})
	an.ignoreNewStatementOnce = true
	an.debug.Printfln("continueStatement: ignoreNewStatementOnce is set to true")
}

// Create a new section as document node, then shift thisNode to a leaf.
func (an *Lexer) createSection() {
	an.endStatement("")
	newSection := new(Section)
	newSection.StatementCounterAtOpening = an.statementCounter
	if an.thisNode == an.rootNode {
		an.debug.Printfln("createSection: root node %p has the new section %p", an.thisNode, newSection)
		an.createLeaf()
		an.thisNode.Entity = newSection
	} else {
		an.debug.Printfln("newSection: node %p has the new section %p", an.thisNode, newSection)
		an.createDocumentSiblingNode(newSection)
	}
	an.createLeaf()
}

/*
Look for a Statement in the previous sibling and remove the sibling node.
Return the Statement if it is found, return nil if not found.
*/
func (an *Lexer) removePreviousSiblingStatement() *Statement {
	index := an.thisNode.GetMyLeafIndex()
	if index == -1 {
		an.debug.Printfln("removePreviousSiblingStatement: cannot find node %p's leaf index",
			an.thisNode)
		return nil
	} else if index == 0 {
		an.debug.Printfln("removePreviousSiblingStatement: this node %p does not have a previous sibling",
			an.thisNode)
		return nil
	}
	// Look for a statement in the previous sibling
	previousSibling := an.thisNode.Parent.Leaves[index-1]
	if obj := previousSibling.Entity; obj == nil {
		an.debug.Printfln("removePreviousSiblingStatement: this node %p's previous leaf %p is empty",
			an.thisNode, previousSibling)
		return nil
	} else if stmt, ok := obj.(*Statement); !ok {
		an.debug.Printfln("removePreviousSiblingStatement: this node %p's previous leaf %p does not hold a statement",
			an.thisNode, previousSibling)
		return nil
	} else {
		// Remove the sibling and return the statement
		an.debug.Printfln("removePreviousSiblingStatement: this node %p's previous leaf %p holds a statement %p",
			an.thisNode, previousSibling, stmt)
		leaves := an.thisNode.Parent.Leaves
		an.thisNode.Parent.Leaves = append(leaves[0:index-1], leaves[index:]...)
		return stmt
	}
}

// Assign the section its first and final statements if necessary, then move thisNode to a new sibling of its parent's.
func (an *Lexer) endSection() {
	if _, sect := an.getSectionState(); sect == nil {
		an.debug.Printfln("endSection: this node %p is not in a section", an.thisNode)
	} else {
		an.endStatement("")
		// Move thisNode to its parent, the one holding *Section
		an.thisNode = an.thisNode.Parent
		an.debug.Printfln("endSection: trying to finish section in node %p", an.thisNode)
		// Calculate the first and final statements
		minNumLeaves := 0
		if an.config.SectionStyle.OpenSectionWithAStatement {
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
				an.debug.Printfln("endSection: first statement is the previous sibling %p", sect.FirstStatement)
			} else {
				sectionFirstLeaf := an.thisNode.Leaves[0]
				if sect.MissingOpeningStatement {
					an.debug.Printfln("endSection: the section is missing its opening statement")
				} else if leafObj := sectionFirstLeaf.Entity; leafObj == nil {
					an.debug.Printfln("endSection: first statement should be the first leaf %p but it holds nothing",
						sectionFirstLeaf)
				} else if stmt, ok := sectionFirstLeaf.Entity.(*Statement); ok {
					an.thisNode.Leaves = an.thisNode.Leaves[1:]
					sect.FirstStatement = stmt
					minNumLeaves++
					an.debug.Printfln("endSection: first statement is the first leaf %p's content, statement %p",
						sectionFirstLeaf, stmt)
				} else {
					an.debug.Printfln("endSection: first statement should be the first leaf %p but it does not hold a statement",
						sectionFirstLeaf)
				}
			}
		}
		if an.config.SectionStyle.CloseSectionWithAStatement {
			/*
				Rather than two scenarios supported by CloseSectionWithAStatement, there is only one scenario
				to deal with here, the section ending must use both prefixes and suffixes, like this:
				<sectionA>
				content
				</sectionA>     <=== "sectionA" is the ending statement
			*/
			if sect.MissingClosingStatement {
				an.debug.Printfln("endSection: although the section requires closing statement, there is none.")
			} else if an.config.SectionStyle.ClosingPrefix != "" && an.config.SectionStyle.ClosingSuffix != "" {
				// minNumLeaves is 0 if section does not begin with statement that is also the section's leaf
				// minNumLeaves is 1 if section begins with a statement that is the section's leaf
				if len(an.thisNode.Leaves) > minNumLeaves {
					lastLeaf := an.thisNode.Leaves[len(an.thisNode.Leaves)-1]
					if lastLeaf == nil {
						an.debug.Printfln("endSection: cannot assign final statement, the last leaf is nil.")
					} else if stmt, ok := lastLeaf.Entity.(*Statement); ok {
						an.thisNode.Leaves = an.thisNode.Leaves[0 : len(an.thisNode.Leaves)-1]
						sect.FinalStatement = stmt
						an.debug.Printfln("endSection: final statement is the last leaf %p's content, statement %p",
							lastLeaf, stmt)
					} else {
						an.debug.Printfln("endSection: final statement should be the last leaf %p but it does not hold a statement",
							lastLeaf)
					}
				} else {
					an.debug.Printfln("endSection: cannot end section with a statement, there are not enough leaves.")
				}
			} else {
				an.debug.Printfln("endSection: the config should specify both prefix and suffix in order to end a section with a statement")
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

func (an *Lexer) getSectionState() (SectionState, *Section) {
	// If thisNode is in a section, the parent node should hold a *Section.
	if an.thisNode.Parent == nil {
		an.debug.Printfln("getSectionState: node %p's parent is nil", an.thisNode)
		return SECTION_STATE_BEFORE_BEGIN, nil
	} else if an.thisNode.Parent.Entity == nil {
		an.debug.Printfln("getSectionState: node %p's parent %p is empty", an.thisNode, an.thisNode.Parent)
		return SECTION_STATE_BEFORE_BEGIN, nil
	}
	section, isSect := an.thisNode.Parent.Entity.(*Section)
	if !isSect {
		an.debug.Printfln("getSectionState: node %p's parent %p holds a %+v, which is not a *Section",
			an.thisNode, an.thisNode.Parent, an.thisNode.Parent.Entity)
		return SECTION_STATE_BEFORE_BEGIN, nil
	}
	var state SectionState
	switch an.config.SectionStyle.SectionMatchMechanism {
	case SECTION_MATCH_FLAT_SINGLE_ANCHOR:
		if section.OpeningPrefix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else {
			state = SECTION_STATE_END_NOW
		}
	case SECTION_MATCH_FLAT_DOUBLE_ANCHOR:
		if section.OpeningPrefix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else if section.OpeningSuffix == "" {
			state = SECTION_STATE_HAS_BEGIN_PREFIX
		} else {
			state = SECTION_STATE_END_NOW
		}
	case SECTION_MATCH_NESTED_DOUBLE_ANCHOR:
		if section.OpeningSuffix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else if section.ClosingSuffix == "" {
			state = SECTION_STATE_HAS_END_PREFIX
		} else {
			state = SECTION_STATE_END_NOW
		}
	case SECTION_MATCH_NESTED_QUAD_ANCHOR:
		if section.OpeningPrefix == "" {
			state = SECTION_STATE_BEFORE_BEGIN
		} else if section.OpeningSuffix == "" {
			state = SECTION_STATE_HAS_BEGIN_PREFIX
		} else if section.ClosingPrefix == "" {
			state = SECTION_STATE_HAS_BEGIN_SUFFIX
		} else if section.ClosingSuffix == "" {
			state = SECTION_STATE_HAS_END_PREFIX
		} else {
			state = SECTION_STATE_END_NOW
		}
	default:
		an.debug.Printfln("getSectionState: unknown SectionMatchMechanism")
	}
	an.debug.Printfln("getSectionState: state is %v, section is %p", state, section)
	return state, section
}

// Save the section opening's prefix marking, create a new section if there is not one.
func (an *Lexer) setSectionOpeningPrefix(prefix string) {
	if an.saveQuoteOrCommentCharacters(prefix) {
		return
	}
	an.endStatement("")
	state, sect := an.getSectionState()
	if state == SECTION_STATE_BEFORE_BEGIN {
		// Create a new section while not being in a section
		an.debug.Printfln("setSectionOpeningPrefix: create a new section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Entity.(*Section).OpeningPrefix = prefix
	} else if an.config.SectionStyle.SectionMatchMechanism == SECTION_MATCH_FLAT_DOUBLE_ANCHOR {
		// Already in a section but document does not allow nested section
		an.debug.Printfln("setSectionOpeningPrefix: end section of node %p and create a new section", an.thisNode.Parent)
		an.endSection()
		an.createSection()
		an.thisNode.Parent.Entity.(*Section).OpeningPrefix = prefix
	} else if state == SECTION_STATE_END_NOW {
		// Marker matches but section should end now
		an.debug.Printfln("setSectionOpeningPrefix: end section right now")
		sect.OpeningPrefix = prefix
		an.endSection()
	} else {
		// Create a nested section
		an.debug.Printfln("setSectionOpeningPrefix: create a nested section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Entity.(*Section).OpeningPrefix = prefix
	}
}

// Save the section opening's suffix marking.
func (an *Lexer) setSectionOpeningSuffix(suffix string) {
	if an.saveQuoteOrCommentCharacters(suffix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state == SECTION_STATE_END_NOW {
		// Marker matches but section should end now
		sect.OpeningSuffix = suffix
		// If statement counter has not increased, then the opening statement does not exist.
		if sect.StatementCounterAtOpening == an.statementCounter {
			sect.MissingOpeningStatement = true
		}
		an.debug.Printfln("setSectionOpeningSuffix: end section right now (missing opening stmt? %v)", sect.MissingOpeningStatement)
		an.endSection()
	} else if an.config.SectionStyle.SectionMatchMechanism == SECTION_MATCH_NESTED_DOUBLE_ANCHOR {
		// Create a section or nested section
		an.debug.Printfln("setSectionOpeningSuffix: create a new section/nested section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Entity.(*Section).OpeningSuffix = suffix
	} else if state < SECTION_STATE_HAS_BEGIN_PREFIX || state > SECTION_STATE_HAS_BEGIN_SUFFIX {
		// State is not right so the marker must have been text
		an.debug.Printfln("setSectionOpeningSuffix: state is not right so only store the characters")
		an.saveMissedCharacters()
	} else {
		// Set suffix if state is right
		sect.OpeningSuffix = suffix
		// If statement counter has not increased, then the opening statement does not exist.
		if sect.StatementCounterAtOpening == an.statementCounter {
			sect.MissingOpeningStatement = true
		}
		an.debug.Printfln("setSectionOpeningSuffix: set suffix (missing opening stmt? %v)", sect.MissingOpeningStatement)
	}
}

// Save the section opening's prefix marking.
func (an *Lexer) setSectionClosingPrefix(prefix string) {
	if an.saveQuoteOrCommentCharacters(prefix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state == SECTION_STATE_END_NOW {
		an.debug.Printfln("setSectionClosingPrefix: end section right now")
		sect.ClosingPrefix = prefix
		sect.StatementCounterAtClosing = an.statementCounter
		an.endSection()
	} else if state < SECTION_STATE_HAS_BEGIN_SUFFIX || state > SECTION_STATE_HAS_END_PREFIX {
		an.debug.Printfln("setSectionClosingPrefix: state is not right so only store the characters")
		an.saveMissedCharacters()
	} else {
		an.debug.Printfln("setSectionClosingPrefix: set prefix")
		sect.ClosingPrefix = prefix
		sect.StatementCounterAtClosing = an.statementCounter
	}
}

// Save the section ending's suffix marking.
func (an *Lexer) setSectionClosingSuffix(suffix string) {
	if an.saveQuoteOrCommentCharacters(suffix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state >= SECTION_STATE_HAS_END_PREFIX {
		an.debug.Printfln("setSectionClosingSuffix: end section right now")
		sect.ClosingSuffix = suffix
		// If statement counter has not increased, then the opening statement does not exist.
		if sect.StatementCounterAtClosing == an.statementCounter {
			sect.MissingClosingStatement = true
		}
		an.endSection()
	} else if state < SECTION_STATE_HAS_END_PREFIX && an.config.SectionStyle.AmbiguousSectionSuffix {
		an.debug.Printfln("setSectionClosingSuffix: call setSectionSetBeginSuffix due to ambiguous suffix choice")
		an.setSectionOpeningSuffix(suffix)
	} else {
		an.debug.Printfln("setSectionClosingSuffix: set suffix")
		sect.ClosingSuffix = suffix
		// If statement counter has not increased, then the opening statement does not exist.
		if sect.StatementCounterAtClosing == an.statementCounter {
			sect.MissingClosingStatement = true
		}
	}
}

// Look for the string from position here. Return the matching string and length of the match.
func (an *Lexer) lookFor(match string) (string, int) {
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
func (an *Lexer) lookForAnyOf(matches []string) (string, int) {
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
func (an *Lexer) lookForSpaces() (string, int) {
	pos := an.herePosition
	for ; pos < len(an.textInput); pos++ {
		if an.textInput[pos] != ' ' && an.textInput[pos] != '\t' {
			break
		}
	}
	return an.textInput[an.herePosition:pos], pos - an.herePosition
}

// Toggle text quoting in the lexer' context.
func (an *Lexer) setQuote(quoteStyle string) {
	if an.contextComment != nil {
		an.saveMissedCharacters()
		an.debug.Printfln("setQuote: quote '%s' goes into context comment", quoteStyle)
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
		an.debug.Printfln("setQuote: begin quoting in text %p", an.contextText)
		an.contextText.QuoteStyle = quoteStyle
	} else {
		if an.contextText.QuoteStyle == quoteStyle {
			an.debug.Printfln("setQuote: finish quoting in text %p", an.contextText)
			an.endText()
		} else {
			an.debug.Printfln("setQuote: quote '%s' goes into context text %p", an.contextText)
			an.saveMissedCharacters()
			an.contextText.Text += quoteStyle
		}
	}
}

// Tell the lexer to open a comment if the text at position here matches any comment opening style.
func (an *Lexer) isOpeningComment() int {
	if an.contextComment != nil {
		// A comment is already open, so it is not possible to open another comment.
		return 0
	}
	for _, style := range an.config.CommentStyles {
		if match, advance := an.lookFor(style.Opening); advance > 0 {
			an.debug.Printfln("Comment opening: %s", match)
			an.saveMissedCharacters()
			an.endText()
			an.createCommentIfNil(style)
			return advance
		}
	}
	return 0
}

// Tell the lexer to close a comment if the text at position here matches any comment closing style.
func (an *Lexer) isClosingComment() int {
	if an.contextComment == nil {
		// Comment has not been opened, so it is not possible to close a comment.
		return 0
	}
	for _, style := range an.config.CommentStyles {
		if match, advance := an.lookFor(style.Closing); advance > 0 {
			if match == an.contextComment.CommentStyle.Closing {
				an.debug.Printfln("Comment closing: %s", match)
				an.endComment(match, true)
				return advance
			}
		}
	}
	return 0
}

// Break down input text according to lexer's configuration. Return the root document node.
func (an *Lexer) Run() *DocumentNode {
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
			an.debug.Printfln("Quote: %s", match)
			an.setQuote(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if spaces, advance = an.lookForSpaces(); advance > 0 {
			an.debug.Printfln("Spaces: length %d", advance)
			an.saveSpaces(spaces)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookForAnyOf(an.config.TokenBreakMarkers); advance > 0 {
			an.debug.Printfln("Breaks: %s", match)
			an.breakText(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookForAnyOf(an.config.StatementContinuationMarkers); advance > 0 {
			an.debug.Printfln("Statement continuation: %s", match)
			an.continueStatement(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookForAnyOf(an.config.StatementEndingMarkers); advance > 0 {
			an.debug.Printfln("Statement ending: %v", []byte(match))
			an.endStatement(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.ClosingSuffix); advance > 0 {
			an.debug.Printfln("Section closing suffix: %s", match)
			an.setSectionClosingSuffix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.ClosingPrefix); advance > 0 {
			an.debug.Printfln("Section closing prefix: %s", match)
			an.setSectionClosingPrefix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.OpeningSuffix); advance > 0 {
			an.debug.Printfln("Section opening suffix: %s", match)
			an.setSectionOpeningSuffix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else if match, advance = an.lookFor(an.config.SectionStyle.OpeningPrefix); advance > 0 {
			an.debug.Printfln("Section opening prefix: %s", match)
			an.setSectionOpeningPrefix(match)
			an.previousMarkerPosition = an.herePosition + advance
		} else {
			advance = 1
		}
	}
	an.debug.Printfln("Run: end statement for the last time")
	an.endStatement("")
	an.debug.Printfln("Run: end all open sections for the last time")
	// End all sections
	for recursionGuard := 0; recursionGuard < 100 && an.thisNode.Parent != nil && an.thisNode.Parent != an.rootNode; recursionGuard++ {
		an.endSection()
	}
	return an.rootNode
}
