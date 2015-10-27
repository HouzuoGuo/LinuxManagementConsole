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

func (an *Analyser) LookForSpaces() (string, int) {
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
	return an.text[an.here: pos], pos - an.here
}

func (an *Analyser) Analyse() {
	var adv int
	var spaces string
	for an.here = 0; an.here < len(an.text); an.here += adv {
		var style string
		if style, adv = an.LookFor(an.Style.CommentBegin); adv > 0 {
			an.NewComment(style)
		} else if style, adv = an.LookFor(an.Style.Quote); adv > 0 {
			an.SetQuote(style)
		} else if spaces, adv = an.LookForSpaces(); adv > 0 {
			an.SetTrailingSpacesOrIndent(spaces)
		} else if style, adv = an.LookFor(an.Style.StmtContinue); adv > 0 {
			an.ContinueStmt(style)
		} else if style, adv = an.LookFor(an.Style.StmtEnd); adv > 0 {
			an.EndStmt(style)
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