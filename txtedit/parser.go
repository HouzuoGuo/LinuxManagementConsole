package txtedit
import "fmt"

func (an *Analyser) LookFor(match []string) (string, int) {
	for _, style := range match {
		if an.here + len(match) > len(an.text) {
			continue
		} else if string(an.text[an.here:an.here + len(match)]) != style {
			continue
		} else {
			return style, len(style)
		}
	}
	return "", 0
}

func (an *Analyser) LookForSpaces() (string, int) {
	pos := an.here
	for ; pos < len(an.text); pos++ {
		if an.text[pos] != ' ' && an.text[pos] != '\t' {
			break
		}
	}
	return an.text[an.here: pos], pos - an.here
}

func (an *Analyser) Analyse() {
	var adv int
	var spaces string
	for an.here = 0; an.here < len(an.text); an.here += adv {
		var style string
		if style, adv = an.LookFor(an.Style.CommentBegin); adv > 0 {
			fmt.Println("Comment: " + style)
			an.NewComment(style)
			an.lastBranch = an.here + adv
		} else if style, adv = an.LookFor(an.Style.Quote); adv > 0 {
			fmt.Println("Quote: " + style)
			an.lastBranch = an.here+ adv
		} else if spaces, adv = an.LookForSpaces(); adv > 0 {
			fmt.Println("Spaces: ", adv, spaces)
			an.storeSpaces(spaces)
			an.lastBranch = an.here + adv
		} else if style, adv = an.LookFor(an.Style.StmtContinue); adv > 0 {
			fmt.Println("StmtContinue: " + style)
			an.lastBranch = an.here+ adv
		} else if style, adv = an.LookFor(an.Style.StmtEnd); adv > 0 {
			fmt.Println("StmtEnd: " + style)
			an.EndStmt()
			an.lastBranch = an.here+ adv
		} else if style, adv = an.LookFor(an.Style.SectBeginPrefix); adv > 0 {
			fmt.Println("SectBeginPrefix: " + style)
			an.lastBranch = an.here+ adv
		} else if style, adv = an.LookFor(an.Style.SectBeginSuffix); adv > 0 {
			fmt.Println("SectBeginSuffix: " + style)
			an.lastBranch = an.here+ adv
		} else if style, adv = an.LookFor(an.Style.SectEndPrefix); adv > 0 {
			fmt.Println("SectEndPrefix: " + style)
			an.lastBranch = an.here+ adv
		} else if style, adv = an.LookFor(an.Style.SectEndSuffix); adv > 0 {
			fmt.Println("SectEndSuffix: " + style)
			an.lastBranch = an.here+ adv
		} else {
			fmt.Println("text '" + string(an.text[an.here]) + "' does not match any condition")
			adv = 1
		}
	}
	an.EndStmt()
}