package txtedit

func (an *Analyser) LookFor(match []string) (string, int) {
	for _, style := range match {
		if an.here + len(match) >= len(an.text) {
			continue
		} else if string(an.text[an.here:an.here + len(match)]) != match {
			continue
		} else {
			return style, len(style)
		}
	}
	return "", 0
}

func (an *Analyser) LookForSpaces() int {
	pos := an.here
	for ; pos < len(an.text); pos++ {
		switch an.text[pos] {
		case ' ':
			continue
		case '\t':
			continue
		default:
			break
		}
	}
	return pos - an.here
}

func (an *Analyser) Analyse() {
	var adv int
	for an.here = 0; an.here < len(an.text); an.here += adv {
		var style string
		if style, adv = an.LookFor(an.Style.CommentBegin); adv > 0 {
			an.NewComment(style)
		} else if style, adv = an.LookFor(an.Style.Quote); adv > 0 {
			an.SetQuote(style)
		} else if adv = an.LookForSpaces(); adv > 0 {
			// Either
			an.SetTrailingSpaces(adv)
			an.SetIndent(adv)
		} else if style, adv = an.LookFor(an.Style.StmtContinue); adv > 0 {
			an.ContinueStmt(style, true)
		} else if style, adv = an.LookFor(an.Style.StmtEnd); adv > 0 {
			if !an.ContinueStmt {
				an.EndStmt(true)
			}
		} else if style, adv = an.LookFor(an.Style.SectBeginPrefix); adv > 0 {
			an.NewSection(style)
		} else if style, adv = an.LookFor(an.Style.SectBeginSuffix); adv > 0 {
			an.SetSectBeginSuffix(style)
		} else if style, adv = an.LookFor(an.Style.SectEndPrefix); adv > 0 {
			an.SetSectEndPrefix(style)
		} else if style, adv = an.LookFor(an.Style.SectEndSuffix); adv > 0 {
			an.SetSectEndSuffix(style)
		} else {
			adv = 1
		}
	}
}