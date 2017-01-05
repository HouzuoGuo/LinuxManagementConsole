package lexer

import (
	"bytes"
	"fmt"
	"strings"
)

/*
Document entities (such as Text, Comment, StatementContinue) as well as certain types of document
nodes (such as Statement) store characters from the original document. The interface provides
DebugInfo() for inspecting the attributes of such objects in human-readable form, and RecoverText()
for reproducing exact characters that were written in the original document.
*/
type ContainVerbatimText interface {
	DebugInfo() string
	VerbatimText() string
}

// Text is optionally surrounded by quotation marks and trailing spaces.
type Text struct {
	QuoteStyle     string
	Text           string
	TrailingSpaces string
}

func (txt *Text) DebugInfo() string {
	return fmt.Sprintf("Quote[%s] Text[%s] Trailing[%s]", txt.QuoteStyle, txt.Text, txt.TrailingSpaces)
}
func (txt *Text) VerbatimText() string {
	return fmt.Sprintf("%s%s%s%s", txt.QuoteStyle, txt.Text, txt.QuoteStyle, txt.TrailingSpaces)
}

// Comment is led by a single marker, or surrounded by a pair of markers, depending on the style.
type Comment struct {
	CommentStyle CommentStyle
	Closed       bool // style carries comment anchors, this is true if the comment has a closing anchor.
	Content      string
}

func (comment *Comment) DebugInfo() string {
	return fmt.Sprintf("Comment[%s] Content[%s]", comment.CommentStyle.Opening, comment.Content)
}
func (comment *Comment) VerbatimText() string {
	closing := comment.CommentStyle.Closing
	if !comment.Closed {
		closing = ""
	}
	return fmt.Sprintf("%s%s%s", comment.CommentStyle.Opening, comment.Content, closing)
}

// Continuation marker leads to the merge of pieces from both current and the next statement.
type StatementContinue struct {
	Style string
}

func (cont *StatementContinue) DebugInfo() string {
	return fmt.Sprintf("Continue[%s]", cont.Style)
}
func (cont *StatementContinue) VerbatimText() string {
	return cont.Style
}

/*
Statement is made of a leading indentation, a suffix ending, pieces of document entities such as
texts, comments, and continuation markers.
Statement is an entity of DocumentNode.
*/
type Statement struct {
	Indent string                // the leading spaces or tabs that indent the statement
	Pieces []ContainVerbatimText // pieces can be anything (e.g. Text, Comment) but Statement.
	Ending string                // the suffix (such as new-line character) that marks end of the statement
}

func (stmt *Statement) DebugInfo() string {
	var out bytes.Buffer
	for _, piece := range stmt.Pieces {
		out.WriteString("[" + piece.DebugInfo() + "]")
	}
	return fmt.Sprintf("Indent[%s] Pieces%s End[%v]", stmt.Indent, out.String(), []byte(stmt.Ending))
}
func (stmt *Statement) VerbatimText() string {
	var out bytes.Buffer
	out.WriteString(stmt.Indent)
	for _, piece := range stmt.Pieces {
		out.WriteString(piece.VerbatimText())
	}
	out.WriteString(stmt.Ending)
	return out.String()
}

/*
Skipping the first N pieces (no matter text or not), then scan through the remaining pieces, looking for the
specified string in the text pieces. Return index of the text piece.
*/
func (stmt *Statement) IndexOfText(skip int, str string, ignoreCase bool) int {
	if ignoreCase {
		str = strings.ToLower(str)
	}
	for i, piece := range stmt.Pieces[skip+1:] {
		switch thing := piece.(type) {
		case *Text:
			if ignoreCase && strings.ToLower(thing.Text) == str || !ignoreCase && thing.Text == str {
				return i
			}
		}
	}
	return -1
}

/*
Skipping the first N pieces (no matter text or not), then scan through the remaining pieces, looking for the
sequence of strings in the text pieces, and return index to each piece.
*/
func (stmt *Statement) IndexOfTextSeq(skip int, seq []string, ignoreCase bool) (indexes []int, sufficient bool) {
	indexes = make([]int, 0, len(seq))
	if ignoreCase {
		for i, str := range seq {
			seq[i] = strings.ToLower(str)
		}
	}
	for seqPiece, stmtPiece := 0, skip+1; seqPiece < len(seq) && stmtPiece < len(stmt.Pieces); {
		switch thing := stmt.Pieces[stmtPiece].(type) {
		case *Text:
			if ignoreCase && strings.ToLower(thing.Text) == seq[seqPiece] || !ignoreCase && thing.Text == seq[seqPiece] {
				indexes = append(indexes, stmtPiece)
				seqPiece++
				continue
			}
		}
		stmtPiece++
	}
	sufficient = len(indexes) == len(seq)
	return
}

/*
Section's opening and closing are determined by markers. Optionally, the markers surround statements.
Section is an entity of DocumentNode. Section content such as statements and nested sections are
stored in leaves of DocumentNode. Verbatim text of a Section cannot be recovered via Section alone,
it is recovered via DocumentNode.VerbatimText().
*/
type Section struct {
	FirstStatement               *Statement
	OpeningPrefix, OpeningSuffix string
	ClosingPrefix, ClosingSuffix string
	FinalStatement               *Statement

	/*
	 The following flags are used when section's opening is marked by both prefix and suffix,
	 or its closing is marked by both prefix and suffix.
	*/

	StatementCounterAtOpening int
	MissingOpeningStatement   bool
	StatementCounterAtClosing int
	MissingClosingStatement   bool
}

func (sect *Section) DebugInfo() string {
	beginStmtStr := ""
	if sect.FirstStatement != nil {
		beginStmtStr = sect.FirstStatement.DebugInfo()
	}
	endStmtStr := ""
	if sect.FinalStatement != nil {
		endStmtStr = sect.FinalStatement.DebugInfo()
	}
	return fmt.Sprintf("Section %s%s%s Closing with %s%s%s",
		sect.OpeningPrefix, beginStmtStr, sect.OpeningSuffix,
		sect.ClosingPrefix, endStmtStr, sect.ClosingSuffix)
}

/*
DocumentNode contains a Section or Statement as its entity.
Section node has leaf section/statement as its entity.
The root DocumentNode can recover verbatim text of input document.
*/
type DocumentNode struct {
	Parent *DocumentNode
	Entity interface{} // pointer to Statement or Section
	Leaves []*DocumentNode
}

// Return the index of this node among its parent's leaves. Return -1 if parent is nil or this leaf is not found.
func (node *DocumentNode) GetMyLeafIndex() int {
	if node.Parent == nil {
		return -1
	}
	for i, leaf := range node.Parent.Leaves {
		if leaf == node {
			return i
		}
	}
	return -1
}

func (node *DocumentNode) VerbatimText() string {
	var out bytes.Buffer
	section, isSection := node.Entity.(*Section)
	if isSection {
		// Write section opening prefix, first statement, and suffix.
		out.WriteString(section.OpeningPrefix)
		if section.FirstStatement != nil {
			out.WriteString(section.FirstStatement.VerbatimText())
		}
		out.WriteString(section.OpeningSuffix)
	} else if node.Entity != nil {
		out.WriteString(node.Entity.(ContainVerbatimText).VerbatimText())
	}
	for _, leaf := range node.Leaves {
		out.WriteString(leaf.VerbatimText())
	}
	if isSection {
		// Write section closing prefix, final statement, and suffix.
		out.WriteString(section.ClosingPrefix)
		if section.FinalStatement != nil {
			out.WriteString(section.FinalStatement.VerbatimText())
		}
		out.WriteString(section.ClosingSuffix)
	}
	return out.String()
}

/*
If this node has leaves and the leaf is among them, return the index of the leaf.
Otherwise, return -1.
*/
func (node *DocumentNode) FindLeafIndex(leaf *DocumentNode) int {
	if node.Leaves == nil || len(node.Leaves) == 0 {
		return -1
	}
	for i, aLeaf := range node.Leaves {
		if aLeaf == leaf {
			return i
		}
	}
	return -1
}

// Remove reference to this node from the parent's leaves. Return true only if the node has been deleted.
func (node *DocumentNode) DeleteSelf() bool {
	if node.Parent == nil {
		return false
	}
	i := node.Parent.FindLeafIndex(node)
	if i == -1 {
		return false
	}
	leaves := node.Parent.Leaves
	node.Parent.Leaves = append(leaves[:i], leaves[i+1:]...)
	return true
}

// Place the new node before this node in the parent's leaves. Return true only if the new node has been placed.
func (node *DocumentNode) InsertBeforeSelf(newNode *DocumentNode) bool {
	if node.Parent == nil {
		return false
	}
	i := node.Parent.FindLeafIndex(node)
	if i == -1 {
		return false // this node is not found among parent's leaves
	}
	// newLeaves = [leaves before i], newNode, [leaves at and after i]
	newLeaves := make([]*DocumentNode, len(node.Parent.Leaves)+1)
	copy(newLeaves, node.Parent.Leaves[:i])
	newLeaves[i] = newNode
	copy(newLeaves[i+1:], node.Parent.Leaves[i:])
	node.Parent.Leaves = newLeaves
	newNode.Parent = node.Parent
	return true

}

// Place the new node after this node in the parent's leaves. Return true only if the new node has been placed.
func (node *DocumentNode) InsertAfterSelf(newNode *DocumentNode) bool {
	if node.Parent == nil {
		return false
	}
	i := node.Parent.FindLeafIndex(node)
	if i == -1 {
		return false // this node is not found among parent's leaves
	}
	newLeaves := make([]*DocumentNode, len(node.Parent.Leaves)+1)
	// newLeaves = [leaves before and at i], newNode, [leaves after i]
	copy(newLeaves, node.Parent.Leaves[:i+1])
	newLeaves[i+1] = newNode
	copy(newLeaves[i+2:], node.Parent.Leaves[i:+1])
	node.Parent.Leaves = newLeaves
	newNode.Parent = node.Parent
	return true
}

/*
If this node has leaves and the leaf node is among them, insert the new node right before the leaf node.
If this node does not have leaves and leaf node is nil, the new node will be made the first leave of this node.
Return true only if the new node has been placed.
*/
func (node *DocumentNode) InsertBefore(leaf *DocumentNode, newNode *DocumentNode) bool {
	if node.Leaves == nil || len(node.Leaves) == 0 {
		if leaf == nil {
			// New node is the first leaf
			node.Leaves = make([]*DocumentNode, 0, 0)
			node.Leaves = append(node.Leaves, newNode)
			newNode.Parent = node
			return true
		}
		return false
	}
	// Insert new node among leaves
	i := node.FindLeafIndex(leaf)
	if i == -1 {
		return false // leaf is not found
	}
	// newLeaves = [leaves before i], newNode, [leaves at and after i]
	newLeaves := make([]*DocumentNode, len(node.Leaves)+1)
	copy(newLeaves, node.Leaves[:i])
	newLeaves[i] = newNode
	copy(newLeaves[i+1:], node.Leaves[i:])
	node.Leaves = newLeaves
	newNode.Parent = node
	return true
}

/*
If this node has leaves and the leaf node is among them, insert the new node right after the leaf node.
If this node does not have leaves and leaf node is nil, the new node will be made the first leave of this node.
Return true only if the new node has been placed.
*/
func (node *DocumentNode) InsertAfter(leaf *DocumentNode, newNode *DocumentNode) bool {
	if node.Leaves == nil || len(node.Leaves) == 0 {
		if leaf == nil {
			// New node is the first leaf
			node.Leaves = make([]*DocumentNode, 0, 0)
			node.Leaves = append(node.Leaves, newNode)
			newNode.Parent = node
			return true
		}
		return false
	}
	// Insert new node among leaves
	i := node.FindLeafIndex(leaf)
	if i == -1 {
		return false // leaf is not found
	}
	newLeaves := make([]*DocumentNode, len(node.Leaves)+1)
	// newLeaves = [leaves before and at i], newNode, [leaves after i]
	copy(newLeaves, node.Leaves[:i+1])
	newLeaves[i+1] = newNode
	copy(newLeaves[i+2:], node.Leaves[i:+1])
	node.Leaves = newLeaves
	newNode.Parent = node
	return true
}

// A criteria to match among nodes that are being traversed and searched.
type MatchCriteria interface {
	Match(*DocumentNode) bool
}

// Match this node against a set of criteria.
func (node *DocumentNode) Match(criteria ...MatchCriteria) bool {
	for _, c := range criteria {
		if !c.Match(node) {
			return false
		}
	}
	return true
}

// Match each leaf against set of criteria, return all matched leaves.
func (node *DocumentNode) SearchLeaves(criteria ...MatchCriteria) (matches []*DocumentNode) {
	matches = make([]*DocumentNode, 0, 0)
	if node.Leaves == nil {
		return
	}
	for _, leaf := range node.Leaves {
		if leaf.Match(criteria...) {
			matches = append(matches, leaf)
		}
	}
	return
}

// Match each leaf against set of criteria, recursively to leaves of the leaf, return all matched leaves.
func (node *DocumentNode) SearchLeavesRecursively(criteria ...MatchCriteria) (matches []*DocumentNode) {
	matches = make([]*DocumentNode, 0, 0)
	if node.Leaves == nil {
		return
	}
	for _, leaf := range node.Leaves {
		if leaf.Match(criteria...) {
			matches = append(matches, leaf)
		} else {
			matches = append(matches, leaf.SearchLeavesRecursively(criteria...)...)
		}
	}
	return
}
