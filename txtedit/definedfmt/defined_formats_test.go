package definedfmt

import (
	"fmt"
	"github.com/HouzuoGuo/LinuxManagementConsole/txtedit"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var sampleTextLocation = os.Getenv("GOPATH") + "/src/github.com/HouzuoGuo/LinuxManagementConsole/txtedit/samples/"

var samples = []struct {
	config   txtedit.AnalyserConfig
	fileName string
}{
	{Sysconfig, "sysconfig"},
	{Systemd, "systemd"},
	{Sysctl, "sysctl"}}

func GetTextAround(str string, pos, length int) (ret string) {
	startPos := pos - length
	if startPos < 0 {
		startPos = 0
	}
	endPos := pos + length
	if endPos >= len(str) {
		endPos = len(str)
	}
	return str[startPos:endPos]
}

func TestTextBreakdown(t *testing.T) {
	for _, sample := range samples {
		txtInput, err := ioutil.ReadFile(path.Join(sampleTextLocation + sample.fileName))
		if err != nil {
			t.Fatal(err)
		}
		txtInputStr := string(txtInput)

		analyser := txtedit.NewAnalyser(txtInputStr, &sample.config, &txtedit.PrintDebugger{})
		rootNode := analyser.Run()
		reproducedText := rootNode.TextString()
		fmt.Println(txtedit.DebugNode(rootNode, 0))
		for i, ch := range txtInputStr {
			if ch != rune(reproducedText[i]) {
				t.Fatalf("Mismatch in file %s, at position %d\n====should read====\n%s\n====reproduced====\n%s\n",
					sample.fileName, i, GetTextAround(txtInputStr, i, 32), GetTextAround(reproducedText, i, 32))
			}
		}
	}
}
