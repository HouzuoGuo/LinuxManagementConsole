package db

import "github.com/HouzuoGuo/LinuxManagementConsole/txtedit/analyser"

type NodePtr struct {
	FileName string
	Node     *analyser.DocumentNode
}

type DB struct {
	GlobPatterns []string
	FileRoot     map[string]*analyser.DocumentNode
}
