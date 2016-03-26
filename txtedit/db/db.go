package db

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit/analyser"

type NodePtr struct {
	FileID int
	Node   *analyser.DocumentNode
}

type DB struct {
	GlobPatterns []string
	FileRoot     map[string]*analyser.DocumentNode
}
