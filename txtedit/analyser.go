package txtedit

import (
	"fmt"
)

/*
Text analyser analyses input text character by character, breaks down the whole document into smaller
pieces that are easier for further analysis and reproduction of document text.
*/
type Analyser struct {
	config   *AnalyserConfig  // document style specification and more configuration
	debug    AnalyzerDebugger // output from analyser's progress, and output of debug information
	rootNode *DocumentNode    // the root node of the broken down document

	textInput          string        // the original input text
	lastBranchPosition int           // the character index where previous text entity/node was created
	here               int           // index of the current character where analyser has progressed
	thisNode           *DocumentNode // reference to the current document node

	ignoreNewStatementOnce bool       // do not create the next new statement caused by statement continuation marker
	textContext            *Text      // reference to the current text entity
	commentContext         *Comment   // reference to the current comment entity
	statementContext       *Statement // reference to the current statement
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
func (an *Analyser) createSiblingNode() {
	if an.thisNode.Obj == nil {
		an.debug.Printf("createSiblingNode: does nothing when this node %p is still empty", an.thisNode)
		return
	}
	if parent := an.thisNode.Parent; parent == nil {
		/*
			In case this node is the root node, create a new root node and make both original
			root node and new sibling node leaves.
		*/
		originalRoot := an.rootNode
		newRoot := &DocumentNode{Parent: nil, Leaves: make([]*DocumentNode, 0, 8)}
		newRoot.Leaves = append(newRoot.Leaves, originalRoot)
		newLeaf := &DocumentNode{Parent: newRoot, Leaves: make([]*DocumentNode, 0, 8)}
		newRoot.Leaves = append(newRoot.Leaves, newLeaf)
		an.rootNode = newRoot
		an.thisNode = newLeaf
		an.debug.Printf("createSiblingNode: new root is %p, original root %p is now a leaf, new sibling is %p",
			an.rootNode, originalRoot, newLeaf)
	} else {
		newLeaf := &DocumentNode{Parent: parent, Leaves: make([]*DocumentNode, 0, 8)}
		parent.Leaves = append(parent.Leaves, newLeaf)
		an.thisNode = newLeaf
		an.debug.Printf("createSiblingNode: new sibling is %p", newLeaf)
	}
}

/*
If the current node already holds an object, then create a new sibling.
Save an object into the current node.
*/
func (an *Analyser) saveNodeAndCreateSibling(saveObj interface{}) {
	if saveObj == nil {
		an.debug.Printf("saveNodeAndCreateSibling: does nothing when object to save is nil")
		return
	}
	if an.thisNode.Obj == nil {
		fmt.Printf("assign %+v to %p\n", saveObj, an.thisNode)
		an.thisNode.Obj = saveObj
	} else {
		fmt.Printf("create sibling and assign %p to %p\n", saveObj, an.thisNode)
		// Must not overwrite the object in this node
		an.createSiblingNode()
		an.thisNode.Obj = saveObj
	}
}

// Create a new leaf node and move reference to the new leaf.
func (an *Analyser) createNewLeaf() {
	newLeaf := &DocumentNode{Parent: an.thisNode, Leaves: make([]*DocumentNode, 0, 8)}
	an.thisNode.Leaves = append(an.thisNode.Leaves, newLeaf)
	an.debug.Printf("createNewLeaf: %p now has a new leaf %p", an.thisNode, newLeaf)
	an.thisNode = newLeaf
}

// If comment context is nil, assign the context a new comment entity.
func (an *Analyser) newComment(commentStyle string) {
	if an.commentContext == nil {
		an.commentContext = new(Comment)
		an.commentContext.CommentStyle = commentStyle
		an.debug.Printf("newComment: context comment is assigned to %p", an.commentContext)
	}
}

// If comment context is not nil, move the comment into statement context and clear comment context.
func (an *Analyser) endComment() {
	if an.commentContext == nil {
		return
	}
	an.savePendingTextOrComment()
	an.newStatement()
	an.statementContext.Pieces = append(an.statementContext.Pieces, an.commentContext)
	an.debug.Printf("endComment: comment %p is now a piece of statement %p", an.commentContext, an.statementContext)
	an.commentContext = nil
}

// If text context is nil, assign the context a new text entity.
func (an *Analyser) newText() {
	if an.textContext == nil {
		an.textContext = new(Text)
		an.debug.Printf("newText: context text is assigned to %p", an.textContext)
	}
}

// If text context is not nil, move the text into statement context and clear text context.
func (an *Analyser) endText() {
	if an.textContext == nil {
		return
	}
	an.savePendingTextOrComment()
	an.newStatement()
	an.statementContext.Pieces = append(an.statementContext.Pieces, an.textContext)
	an.debug.Printf("endText: text %p is now a piece of statement %p", an.textContext, an.statementContext)
	an.textContext = nil
}

// If statement context is nil, assign the context a new statement entity.
func (an *Analyser) newStatement() {
	if an.statementContext == nil {
		an.statementContext = new(Statement)
		an.saveNodeAndCreateSibling(an.statementContext)
		fmt.Println("newStatement: context statement is assigned to %p", an.statementContext)
	}
}

// Move context text and comment into context statement (create new statement if necessary), and clear context statement.
func (an *Analyser) endStatement(ending string) {
	// Organise context objects
	an.savePendingTextOrComment()
	an.endComment()
	an.endText()
	if an.ignoreNewStatementOnce {
		an.debug.Printf("endStatement: not creating new document node when ignoreNewStatementOnce is set")
		an.ignoreNewStatementOnce = false
		return
	}
	if an.statementContext == nil && an.commentContext == nil && an.textContext == nil {
		an.debug.Printf("endStatement: all contexts are nil, nothing to save")
		if ending != "" {
			an.debug.Printf("endStatement: save statement ending in a new statement")
			an.saveNodeAndCreateSibling(&Statement{Ending: ending})
		}
		return
	}
	an.statementContext.Ending = ending
	if an.commentContext != nil {
		an.statementContext.Pieces = append(an.statementContext.Pieces, an.commentContext)
	}
	if an.textContext != nil {
		an.statementContext.Pieces = append(an.statementContext.Pieces, an.textContext)
	}
	// Remember - statement context was placed in the document node tree when it was created
	an.commentContext = nil
	an.textContext = nil
	an.statementContext = nil
}

func (an *Analyser) savePendingTextOrComment() {
	if an.here-an.lastBranchPosition > 0 {
		missedContent := an.textInput[an.lastBranchPosition:an.here]
		if an.commentContext != nil {
			fmt.Println("missed content", missedContent, "will be stored in comment")
			an.commentContext.Content += missedContent
		} else {
			fmt.Println("missed content", missedContent, "will be stored in val")
			an.newText()
			an.textContext.Text += missedContent
		}
		an.lastBranchPosition = an.here
	} else {
		fmt.Println("storeContent does nothing")
	}
}
func (an *Analyser) storeSpaces(spaces string) {
	an.savePendingTextOrComment()
	fmt.Println("About to store space '" + spaces + "'")
	if an.ignoreNewStatementOnce {
		fmt.Println("Spaces are going into new val")
		an.endText()
		an.newText()
		an.textContext.TrailingSpaces += spaces
		an.endText()
	} else if an.commentContext != nil {
		fmt.Println("Spaces are going into context comment")
		an.commentContext.Content += spaces
	} else if an.textContext != nil {
		fmt.Println("Spaces are going into value trailing spaces")
		an.textContext.TrailingSpaces = spaces
		an.endText()
	} else if an.statementContext != nil {
		fmt.Println("Spaces set indent")
		an.statementContext.Indent += spaces
	} else if an.statementContext == nil {
		fmt.Println("Spaces set indent and makes a new statement")
		an.newStatement()
		an.statementContext.Indent += spaces
	} else {
		fmt.Println("Spaces have no where to go")
	}
}

func (an *Analyser) ContinueStmt(style string) {
	if an.textContext != nil && an.textContext.QuoteStyle != "" {
		an.textContext.Text += style
		fmt.Println("Continue statement mark goes into value")
		return
	}
	an.savePendingTextOrComment()
	an.endComment()
	an.endText()
	an.newStatement()
	an.statementContext.Pieces = append(an.statementContext.Pieces, &StatementContinue{Style: style})
	an.ignoreNewStatementOnce = true
	fmt.Println("Continue statement flag is set")
}

func (an *Analyser) newSection() {
	an.endStatement("")
	newSection := new(Section)
	if an.thisNode == an.rootNode {
		an.debug.Printf("newSection: create new section from root node %p", an.thisNode)
		an.createNewLeaf()
		an.thisNode.Obj = new(Section)
	} else {
		an.debug.Printf("newSection: create new section from node %p", an.thisNode)
		an.saveNodeAndCreateSibling(newSection)
	}
	an.createNewLeaf()
}

func (an *Analyser) IsSect() bool {
	if an.thisNode.Parent == nil || an.thisNode.Parent.Obj == nil {
		return false
	} else if _, isSect := an.thisNode.Parent.Obj.(*Section); isSect {
		return true
	} else {
		return false
	}
}

func (an *Analyser) FindThisLeaf() int {
	if an.thisNode.Parent == nil {
		return -1
	}
	for i, leaf := range an.thisNode.Parent.Leaves {
		if leaf == an.thisNode {
			return i
		}
	}
	return -1
}

func (an *Analyser) GetPreviousLeaf() *Statement {
	thisLeaf := an.FindThisLeaf()
	if thisLeaf == -1 {
		return nil
	}
	if thisLeaf == 0 {
		return nil
	}
	prevLeaf := an.thisNode.Parent.Leaves[thisLeaf-1]
	if stmt, ok := prevLeaf.Obj.(*Statement); ok {
		return stmt
	}
	return nil
}

func (an *Analyser) EndSection() {
	if _, sect := an.GetSectionState(); sect == nil {
		fmt.Println("this is not a section but it ends here, why?")
	} else {
		fmt.Println("section ends here, saving the latest statement")
		an.endStatement("")
		an.thisNode = an.thisNode.Parent
		// an.this is now the parent - section object
		// Remove the last leaf if it holds no object
		if an.thisNode.Leaves[len(an.thisNode.Leaves)-1].Obj == nil {
			an.thisNode.Leaves = an.thisNode.Leaves[:len(an.thisNode.Leaves)-1]
		}
		minNumLeaves := 0
		if an.config.BeginSectionWithAStatement {
			if len(an.config.SectionBeginningSuffixes) == 0 {
				sect.BeginningStatement = an.GetPreviousLeaf()
			} else {
				firstLeaf := an.thisNode.Leaves[0]
				if stmt, ok := firstLeaf.Obj.(*Statement); ok {
					fmt.Println("successfully set sect.begin")
					an.thisNode.Leaves = an.thisNode.Leaves[1:]
					sect.BeginningStatement = stmt
					minNumLeaves++
				}
			}
		}
		if an.config.EndSectionWithAStatement {
			if len(an.config.SectionEndingSuffixes) > 0 {
				if len(an.thisNode.Leaves) > minNumLeaves {
					lastLeaf := an.thisNode.Leaves[len(an.thisNode.Leaves)-1]
					if stmt, ok := lastLeaf.Obj.(*Statement); ok {
						fmt.Println("successfully set sect.end")
						an.thisNode.Leaves = an.thisNode.Leaves[0 : len(an.thisNode.Leaves)-1]
						sect.EndingStatement = stmt
					}
				}
			}
		}
		fmt.Println("Section has ended")
		an.createSiblingNode()
	}
	// Remember - section object was placed in the document node tree when it was created
}

const (
	SECT_STATE_NONE         = 0
	SECT_STATE_BEFORE_BEGIN = 1
	SECT_STATE_BEGIN_PREFIX = 2
	SECT_STATE_BEGIN_SUFFIX = 3
	SECT_STATE_END_PREFIX   = 4
	SECT_STATE_END_NOW      = 5
)

type SectionState int

func (an *Analyser) GetSectionState() (SectionState, *Section) {
	if !an.IsSect() {
		fmt.Println("not in section!!")
		return SECT_STATE_NONE, nil
	}
	sect := an.thisNode.Parent.Obj.(*Section)
	switch an.config.SectionMatchMechanism {
	case SECTION_MATCH_FLAT_SINGLE_ANCHOR:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	case SECTION_MATCH_FLAT_DOUBLE_ANCHOR:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else if sect.BeginSuffix == "" {
			return SECT_STATE_BEGIN_PREFIX, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	case SECTION_MATCH_NESTED_DOUBLE_ANCHOR:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else if sect.EndPrefix == "" {
			return SECT_STATE_BEGIN_SUFFIX, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	case SECTION_MATCH_NESTED_QUAD_ANCHOR:
		if sect.BeginPrefix == "" {
			return SECT_STATE_BEFORE_BEGIN, sect
		} else if sect.BeginSuffix == "" {
			return SECT_STATE_BEGIN_PREFIX, sect
		} else if sect.EndPrefix == "" {
			return SECT_STATE_BEGIN_SUFFIX, sect
		} else if sect.EndSuffix == "" {
			return SECT_STATE_END_PREFIX, sect
		} else {
			return SECT_STATE_END_NOW, sect
		}
	default:
		return SECT_STATE_NONE, sect
	}
}

func (an *Analyser) BeginSectionSetPrefix(style string) {
	an.endStatement("")
	state, sect := an.GetSectionState()
	if state == SECT_STATE_NONE {
		fmt.Println("BeginSectionSetPrefix: create new first-level section")
		an.newSection()
		fmt.Println(an.thisNode.Parent, an.thisNode)
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = style
		an.savePendingTextOrComment()
	} else if state == SECT_STATE_END_NOW {
		fmt.Println("BeginSectionSetPrefix: End now")
		sect.BeginPrefix = style
		an.savePendingTextOrComment()
		an.EndSection()
	} else {
		fmt.Println("BeginSectionSetPrefix: create new sub section")
		an.newSection()
		fmt.Println(an.thisNode.Parent, an.thisNode)
		an.thisNode.Parent.Obj.(*Section).BeginPrefix = style
		an.savePendingTextOrComment()
	}
}

func (an *Analyser) BeginSectionSetSuffix(style string) {
	an.endStatement("")
	if state, sect := an.GetSectionState(); state == SECT_STATE_END_NOW {
		fmt.Println("BeginSectionSetSuffix: set style and end")
		sect.BeginSuffix = style
		an.savePendingTextOrComment()
		an.EndSection()
	} else if state < SECT_STATE_BEGIN_PREFIX || state > SECT_STATE_BEGIN_SUFFIX {
		fmt.Println("BeginSectionSetSuffix: no match, store content")
		an.savePendingTextOrComment()
	} else {
		fmt.Println("BeginSectionSetSuffix: only set style")
		sect.BeginSuffix = style
		an.savePendingTextOrComment()
	}
}

func (an *Analyser) EndSectionSetPrefix(style string) {
	an.endStatement("")
	if state, sect := an.GetSectionState(); state == SECT_STATE_END_NOW {
		fmt.Println("EndSectionSetPrefix: set style and end")
		sect.EndPrefix = style
		an.savePendingTextOrComment()
		an.EndSection()
	} else if state < SECT_STATE_BEGIN_SUFFIX || state > SECT_STATE_END_PREFIX {
		fmt.Println("EndSectionSetPrefix: no match, store content")
		an.savePendingTextOrComment()
	} else {
		fmt.Println("EndSectionSetPrefix: only set style")
		sect.EndPrefix = style
		an.savePendingTextOrComment()
	}
}

func (an *Analyser) EndSectionSetSuffix(style string) {
	an.endStatement("")
	if state, sect := an.GetSectionState(); state >= SECT_STATE_END_PREFIX {
		fmt.Println("EndSectionSetSuffix: set style and end")
		sect.EndSuffix = style
		an.savePendingTextOrComment()
		an.EndSection()
	} else if state < SECT_STATE_END_PREFIX && an.config.AmbiguousSectionSuffix {
		fmt.Println("EndSectionSetSuffix: ambiguous suffix")
		an.BeginSectionSetSuffix(style)
	} else {
		fmt.Println("EndSectionSetSuffix: only set style")
		sect.EndSuffix = style
		an.savePendingTextOrComment()
	}
}
