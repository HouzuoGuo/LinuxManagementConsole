package txtedit

import (
	"fmt"
)

/*
Text analyser analyses input text character by character, breaks down the whole document into smaller
pieces that are easier for further analysis and reproduction of document text.
*/
type Analyser struct {
	textInput string           // the original input text
	config    *AnalyserConfig  // document style specification and more configuration
	debug     AnalyzerDebugger // output from analyser's progress, and output of debug information

	positionLastBranch int           // the character index where previous text entity/node was created
	positionHere       int           // index of the current character where analyser has progressed
	rootNode           *DocumentNode // the root node of the broken down document
	thisNode           *DocumentNode // reference to the current document node

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
	ret.config.DetectSectionMatchMechanism()
	ret.debug.Printf("NewAnalyser: initialised with section match mechanism being %v", ret.config.SectionMatchMechanism)
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
func (an *Analyser) createCommentIfNil(commentStyle string) {
	if an.contextComment == nil {
		an.contextComment = new(Comment)
		an.contextComment.CommentStyle = commentStyle
		an.debug.Printf("createCommentIfNil: context comment is assigned to %p", an.contextComment)
	}
}

// If comment context is not nil, move the comment into statement context and clear comment context.
func (an *Analyser) endComment() {
	if an.contextComment == nil {
		return
	}
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
	if an.contextText != nil && an.contextText.QuoteStyle != "" {
		an.contextText.Text += ending
		an.debug.Printf("endStatement: the statement ending goes into context text %p", an.contextText)
		return
	}
	// Organise context objects
	an.saveMissedCharacters()
	an.endComment()
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

// In the context comment or text, save the characters that have not yet been placed in any entity.
func (an *Analyser) saveMissedCharacters() {
	if an.positionHere-an.positionLastBranch <= 0 {
		return // nothing missed
	}
	missedContent := an.textInput[an.positionLastBranch:an.positionHere]
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
	an.positionLastBranch = an.positionHere
}

// Place the space characters inside statement indentation or text entity's trailing spaces.
func (an *Analyser) saveSpaces(spaces string) {
	an.saveMissedCharacters()
	length := len(spaces)
	if an.ignoreNewStatementOnce {
		an.endText()
		an.createTextIfNil()
		an.debug.Printf("saveSpaces: ignoreNewStatementOnce is true, %d spaces go into text %p", length, an.contextText)
		an.contextText.TrailingSpaces += spaces
		an.endText()
	} else if an.contextComment != nil {
		an.debug.Printf("saveSpaces: %d spaces go into comment %p", length, an.contextComment)
		an.contextComment.Content += spaces
	} else if an.contextText != nil {
		an.debug.Printf("saveSpaces: %d spaces go into text %p", length, an.contextText)
		an.contextText.TrailingSpaces = spaces
		an.endText()
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
		an.debug.Printf("saveQuoteOrCommentCharacters: save '%s' in context comment %p", str, an.contextComment)
		an.contextComment.Content += str
		return true
	}
	if an.contextText != nil && an.contextText.QuoteStyle != "" {
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
	an.endComment()
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

// Look for a Statement in the previous sibling, return it if found, return nil if not found.
func (an *Analyser) getPreviousSiblingStatement() *Statement {
	// Find the index of thisNode among its parent's leaves
	parent := an.thisNode.Parent
	index := -1
	if an.thisNode.Parent == nil {
		an.debug.Printf("getPreviousSiblingStatement: this node %p does not have a parent", an.thisNode)
		return nil
	}
	for i, leaf := range parent.Leaves {
		if leaf == an.thisNode {
			index = i
			break
		}
	}
	if index == -1 {
		an.debug.Printf("getPreviousSiblingStatement: cannot find this node %p among parent %p's leaves", an.thisNode, parent)
		return nil
	} else if index == 0 {
		an.debug.Printf("getPreviousSiblingStatement: this node %p does not have a previous sibling", an.thisNode)
		return nil
	}
	// Look for a statement in the sibling
	previousSibling := an.thisNode.Parent.Leaves[index-1]
	if obj := previousSibling.Obj; obj == nil {
		an.debug.Printf("getPreviousSiblingStatement: this node %p's previs leaf %p is empty", an.thisNode, previousSibling)
		return nil
	} else if stmt, ok := obj.(*Statement); ok {
		an.debug.Printf("getPreviousSiblingStatement: this node %p's previous leaf %p holds a statement %p",
			an.thisNode, previousSibling, stmt)
		return stmt
	} else {
		an.debug.Printf("getPreviousSiblingStatement: this node %p's previous leaf %p does not hold a statement",
			an.thisNode, previousSibling)
		return nil
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
		if an.config.BeginSectionWithAStatement {
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
			if len(an.config.SectionBeginningPrefixes) == 0 {
				sect.FirstStatement = an.getPreviousSiblingStatement()
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
		if an.config.EndSectionWithAStatement {
			/*
				Rather than two scenarios supported by BeginSectionWithAStatement, there is only one scenario
				to deal with here, the section ending must use both prefixes and suffixes, like this:
				<sectionA>
				content
				</sectionA>     <=== "sectionA" is the ending statement
			*/
			if len(an.config.SectionEndingSuffixes) > 0 {
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
	switch an.config.SectionMatchMechanism {
	case SECTION_MATCH_NO_SECTION:
		state = SECTION_STATE_END_NOW
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

// Assign the prefix marking to the section, create a new section if there is not one.
func (an *Analyser) setSectionBeginPrefix(prefix string) {
	if an.saveQuoteOrCommentCharacters(prefix) {
		return
	}
	an.endStatement("")
	state, sect := an.getSectionState()
	if state == SECTION_STATE_BEFORE_BEGIN {
		an.debug.Printf("setSectionPrefix: create a new section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = prefix
	} else if state == SECTION_STATE_END_NOW {
		an.debug.Printf("setSectionPrefix: end now")
		sect.BeginPrefix = prefix
		an.endSection()
	} else if an.config.SectionMatchMechanism == SECTION_MATCH_FLAT_DOUBLE_ANCHOR ||
		an.config.SectionMatchMechanism == SECTION_MATCH_FLAT_DOUBLE_ANCHOR {
		an.debug.Printf("setSectionPrefix: end section of node %p and create a new section", an.thisNode.Parent)
		an.endSection()
		an.createSection()
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = prefix
	} else {
		an.debug.Printf("setSectionPrefix: create a nested section from node %p", an.thisNode)
		an.createSection()
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = prefix
	}
}

func (an *Analyser) setSectionBeginSuffix(suffix string) {
	if an.saveQuoteOrCommentCharacters(suffix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state == SECTION_STATE_END_NOW {
		fmt.Println("BeginSectionSetSuffix: set style and end")
		sect.BeginSuffix = suffix
		an.saveMissedCharacters()
		an.endSection()
	} else if state < SECTION_STATE_HAS_BEGIN_PREFIX || state > SECTION_STATE_HAS_BEGIN_SUFFIX {
		fmt.Println("BeginSectionSetSuffix: no match, store content")
		an.saveMissedCharacters()
	} else {
		fmt.Println("BeginSectionSetSuffix: only set style")
		sect.BeginSuffix = suffix
		an.saveMissedCharacters()
	}
}

func (an *Analyser) setSectionEndPrefix(prefix string) {
	if an.saveQuoteOrCommentCharacters(prefix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state == SECTION_STATE_END_NOW {
		fmt.Println("EndSectionSetPrefix: set style and end")
		sect.EndPrefix = prefix
		an.saveMissedCharacters()
		an.endSection()
	} else if state < SECTION_STATE_HAS_BEGIN_SUFFIX || state > SECTION_STATE_HAS_END_PREFIX {
		fmt.Println("EndSectionSetPrefix: no match, store content")
		an.saveMissedCharacters()
	} else {
		fmt.Println("EndSectionSetPrefix: only set style")
		sect.EndPrefix = prefix
		an.saveMissedCharacters()
	}
}

func (an *Analyser) setSectionEndSuffix(suffix string) {
	if an.saveQuoteOrCommentCharacters(suffix) {
		return
	}
	an.endStatement("")
	if state, sect := an.getSectionState(); state >= SECTION_STATE_HAS_END_PREFIX {
		fmt.Println("EndSectionSetSuffix: set style and end")
		sect.EndSuffix = suffix
		an.saveMissedCharacters()
		an.endSection()
	} else if state < SECTION_STATE_HAS_END_PREFIX && an.config.AmbiguousSectionSuffix {
		fmt.Println("EndSectionSetSuffix: ambiguous suffix")
		an.setSectionBeginSuffix(suffix)
	} else {
		fmt.Println("EndSectionSetSuffix: only set style")
		sect.EndSuffix = suffix
		an.saveMissedCharacters()
	}
}
