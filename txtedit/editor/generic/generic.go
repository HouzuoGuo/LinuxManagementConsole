package generic

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit/lexer"

type GenericNode struct {
	Thing interface{}
	Node  *lexer.DocumentNode
}

type StrNode struct {
	Str  string
	Node *lexer.DocumentNode
}

type IntNode struct {
	Int  int
	Node *lexer.DocumentNode
}
