package txtedit

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
	CommentStyle string
	Content      string
}

func (comment *Comment) DebugString() string {
	return fmt.Sprintf("Comment[%s] Content[%s]", comment.CommentStyle, comment.Content)
}
func (comment *Comment) TextString() string {
	return fmt.Sprintf("%s%s", comment.CommentStyle, comment.Content)
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
	Pieces []interface{} // pieces can be Text, Comment, or StatementContinue
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
Section's beginning and ending are determined by markers. Optionally, the markers surround statements.
Section is a document node with leaves being its content; if a leaf is another Section, then it is a nested section.
Original text is recovered from DocumentNode rather than Section, thus Section does not support EntityToText.
*/
type Section struct {
	FirstStatement           *Statement
	BeginPrefix, BeginSuffix string
	EndPrefix, EndSuffix     string
	FinalStatement           *Statement
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
	return fmt.Sprintf("Section %s%s%s End %s%s%s",
		sect.BeginPrefix, beginStmtStr, sect.BeginSuffix,
		sect.EndPrefix, endStmtStr, sect.EndSuffix)
}

/*
Document nodes together make a tree of entities that represent the original document text.
If a node comes with branches, then the node must be a Section.
Only Statement and Section can be nodes.
*/
type DocumentNode struct {
	Parent *DocumentNode
	Obj    interface{}
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
		// Write section beginning prefix, first statement, and suffix.
		out.WriteString(section.BeginPrefix)
		if section.FirstStatement != nil {
			out.WriteString(section.FirstStatement.TextString())
		}
		out.WriteString(section.BeginSuffix)
	} else if node.Obj != nil {
		out.WriteString(node.Obj.(EntityToText).TextString())
	}
	for _, leaf := range node.Leaves {
		out.WriteString(leaf.TextString())
	}
	if isSection {
		// Write section ending prefix, final statement, and suffix.
		out.WriteString(section.EndPrefix)
		if section.FinalStatement != nil {
			out.WriteString(section.FinalStatement.TextString())
		}
		out.WriteString(section.EndSuffix)
	}
	return out.String()
}
