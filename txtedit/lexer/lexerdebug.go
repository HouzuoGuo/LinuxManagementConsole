package lexer

import (
	"bytes"
	"fmt"
	"strings"
)

// Describe a document node and all of its leaves recursively in a human-readable presentation.
func DebugNode(node *DocumentNode, indent int) string {
	var out bytes.Buffer
	prefixIndent := strings.Repeat(" ", indent)
	if node == nil {
		out.WriteString(prefixIndent + "(nil node)\n")
		return out.String()
	}
	if node.Entity == nil {
		out.WriteString(prefixIndent + "Node - (empty)")
	} else if sect, ok := node.Entity.(*Section); ok {
		// Section does not implement ContainVerbatimText
		out.WriteString(prefixIndent + "Node - " + sect.DebugInfo())
	} else {
		out.WriteString(prefixIndent + "Node - " + node.Entity.(ContainVerbatimText).DebugInfo())
	}
	// Recursively descent into leaves
	if len(node.Leaves) > 0 {
		out.WriteString(" -->\n")
		for _, leaf := range node.Leaves {
			out.WriteString(DebugNode(leaf, indent+2))
		}
	} else {
		out.WriteRune('\n')
	}
	return out.String()
}

// Offer debugging capabilities to a lexer.
type LexerDebug interface {
	Printfln(format string, msg ...interface{})
}

// Not to debug a lexer by discarding debug messages.
type LexerDebugNoop struct {
}

// Do nothing.
func (debug *LexerDebugNoop) Printfln(format string, msg ...interface{}) {
	// Intentionally left blank
}

// Debug a lexer by printing debug messages to standard output.
type LexerDebugStdout struct {
}

// Print debug messages to standard output.
func (debug *LexerDebugStdout) Printfln(format string, msg ...interface{}) {
	fmt.Printf(format, msg...)
	fmt.Println()
}
