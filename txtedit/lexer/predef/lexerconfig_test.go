package predef

import (
	"fmt"
	"github.com/HouzuoGuo/LinuxManagementConsole/txtedit/lexer"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var sampleTextLocation = os.Getenv("GOPATH") + "/src/github.com/HouzuoGuo/LinuxManagementConsole/txtedit/lexer/predef/samples/"

var samples = []struct {
	config   lexer.LexerConfig
	fileName string
}{
	{Sysconfig, "sysconfig"},
	{Systemd, "systemd"},
	{Sysctl, "sysctl"},
	{CronAllowDeny, "cron-allow-deny"},
	{Cron, "cron"},
	{Hosts, "hosts"},
	{Login, "login"},
	{Nsswitch, "nsswitch"},
	{Httpd, "httpd"},
	{NamedZone, "named-zone"},
	{Named, "named"},
	{Dhcpd, "dhcpd"},
	{Ntpd, "ntpd"},
	{Limits, "limits"},
}

func GetTextAround(str string, pos, length int) string {
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

		an := lexer.NewLexer(txtInputStr, &sample.config, &lexer.LexerDebugNoop{})
		fmt.Println("@@@@@@@@@@@@@@Going to analyse", sample.fileName)
		rootNode := an.Run()
		reproducedText := rootNode.VerbatimText()
		fmt.Println(lexer.DebugNode(rootNode, 0))
		lenOriginal := len(txtInputStr)
		lenReproduced := len(reproducedText)
		if lenReproduced >= lenOriginal {
			for i, ch := range txtInputStr {
				if ch != rune(reproducedText[i]) {
					t.Fatalf("Mismatch in file %s, at position %d\n====should read====\n%s\n====reproduced====\n%s\n",
						sample.fileName, i, GetTextAround(txtInputStr, i, 32), GetTextAround(reproducedText, i, 32))
				}
			}
		} else {
			for i, ch := range reproducedText {
				if ch != rune(txtInput[i]) {
					t.Fatalf("Mismatch in file %s, at position %d\n====should read====\n%s\n====reproduced====\n%s\n",
						sample.fileName, i, GetTextAround(txtInputStr, i, 32), GetTextAround(reproducedText, i, 32))
				}
			}
		}
		if lenReproduced > lenOriginal {
			t.Fatalf("Reproduced text is longer, extra is:\n%s", reproducedText[lenOriginal+1:])
		} else if lenReproduced < lenOriginal {
			t.Fatalf("Original text is longer, extra is:\n%s", txtInputStr[lenReproduced+1:])
		}
	}
}
