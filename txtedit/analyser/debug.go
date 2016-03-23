package analyser

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
	if node.Obj == nil {
		out.WriteString(prefixIndent + "Node - (empty)")
	} else {
		out.WriteString(prefixIndent + "Node - " + node.Obj.(EntityDebug).DebugString())
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

// Offer debugging capabilities to an Analyser.
type AnalyserDebugger interface {
	Printf(format string, msg ...interface{})
}

// An AnalyzerDebugger implementation that is silent and does nothing.
type NoopDebugger struct {
}

func (debug *NoopDebugger) Printf(format string, msg ...interface{}) {
	// Intentionally left blank
}

// An AnalyzerDebugger implementation that prints messages to standard output.
type PrintDebugger struct {
}

func (debug *PrintDebugger) Printf(format string, msg ...interface{}) {
	fmt.Printf(format, msg...)
	fmt.Println()
}
