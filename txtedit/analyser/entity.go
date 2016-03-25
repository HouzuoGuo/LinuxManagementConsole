package analyser

import (
	"bytes"
	"fmt"
)

// Get the descriptive debug information from the document entity.
type EntityDebug interface {
	DebugString() string
}

// Recover original string from the document entity.
type EntityToText interface {
	TextString() string
}

// Text is optionally surrounded by quotation marks and trailing spaces.
type Text struct {
	QuoteStyle     string
	Text           string
	TrailingSpaces string
}

func (txt *Text) DebugString() string {
	return fmt.Sprintf("Quote[%s] Text[%s] Trailing[%s]", txt.QuoteStyle, txt.Text, txt.TrailingSpaces)
}
func (txt *Text) TextString() string {
	return fmt.Sprintf("%s%s%s%s", txt.QuoteStyle, txt.Text, txt.QuoteStyle, txt.TrailingSpaces)
}

// Comment is led by a marker.
type Comment struct {
	CommentStyle CommentStyle
	Closed       bool // style carries an opening and closing, this flag is true if the comment is properly closed.
	Content      string
}

func (comment *Comment) DebugString() string {
	return fmt.Sprintf("Comment[%s] Content[%s]", comment.CommentStyle.Opening, comment.Content)
}
func (comment *Comment) TextString() string {
	closing := comment.CommentStyle.Closing
	if !comment.Closed {
		closing = ""
	}
	return fmt.Sprintf("%s%s%s", comment.CommentStyle.Opening, comment.Content, closing)
}

// Continuation marker leads to the combination of pieces from both current and the next statement.
type StatementContinue struct {
	Style string
}

func (cont *StatementContinue) DebugString() string {
	return fmt.Sprintf("Continue[%s]", cont.Style)
}
func (cont *StatementContinue) TextString() string {
	return cont.Style
}

/*
Statement is made of leading indentation, suffix ending, pieces of texts, comments, and continuation markers.
Statement is a document node.
*/
type Statement struct {
	Indent string        // the leading spaces or tabs that indent the statement
	Pieces []interface{} // pieces must be pointers to Text, Comment, or StatementContinue
	Ending string        // the suffix (such as new-line character) that marks end of the statement
}

func (stmt *Statement) DebugString() string {
	var out bytes.Buffer
	for _, piece := range stmt.Pieces {
		out.WriteString("[" + piece.(EntityDebug).DebugString() + "]")
	}
	return fmt.Sprintf("Indent[%s] Pieces%s End[%v]", stmt.Indent, out.String(), []byte(stmt.Ending))
}
func (stmt *Statement) TextString() string {
	var out bytes.Buffer
	out.WriteString(stmt.Indent)
	for _, piece := range stmt.Pieces {
		out.WriteString(piece.(EntityToText).TextString())
	}
	out.WriteString(stmt.Ending)
	return out.String()
}

/*
Section's opening and closing are determined by markers. Optionally, the markers surround statements.
Section is a document node with leaves being its content; if a leaf is another Section, then it is a nested section.
Original text is recovered from DocumentNode rather than Section, thus Section does not support EntityToText.
*/
type Section struct {
	FirstStatement               *Statement
	OpeningPrefix, OpeningSuffix string
	ClosingPrefix, ClosingSuffix string
	FinalStatement               *Statement

	/*
	 The following two flags are used when section's opening is marked by both prefix and suffix,
	 or its closing is marked by both prefix and suffix.
	*/

	StatementCounterAtOpening int
	MissingOpeningStatement   bool
	StatementCounterAtClosing int
	MissingClosingStatement   bool
}

func (sect *Section) DebugString() string {
	beginStmtStr := ""
	if sect.FirstStatement != nil {
		beginStmtStr = sect.FirstStatement.DebugString()
	}
	endStmtStr := ""
	if sect.FinalStatement != nil {
		endStmtStr = sect.FinalStatement.DebugString()
	}
	return fmt.Sprintf("Section %s%s%s Closing with %s%s%s",
		sect.OpeningPrefix, beginStmtStr, sect.OpeningSuffix,
		sect.ClosingPrefix, endStmtStr, sect.ClosingSuffix)
}

/*
Document nodes together make a tree of entities that represent the original document text.
If a node comes with branches, then the node must be a Section.
Only Statement and Section can be nodes.
*/
type DocumentNode struct {
	Parent *DocumentNode
	Obj    interface{} // node content (object) must be a pointer
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

func (node *DocumentNode) TextString() string {
	var out bytes.Buffer
	section, isSection := node.Obj.(*Section)
	if isSection {
		// Write section opening prefix, first statement, and suffix.
		out.WriteString(section.OpeningPrefix)
		if section.FirstStatement != nil {
			out.WriteString(section.FirstStatement.TextString())
		}
		out.WriteString(section.OpeningSuffix)
	} else if node.Obj != nil {
		out.WriteString(node.Obj.(EntityToText).TextString())
	}
	for _, leaf := range node.Leaves {
		out.WriteString(leaf.TextString())
	}
	if isSection {
		// Write section closing prefix, final statement, and suffix.
		out.WriteString(section.ClosingPrefix)
		if section.FinalStatement != nil {
			out.WriteString(section.FinalStatement.TextString())
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
