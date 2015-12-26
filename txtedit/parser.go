package txtedit

import "fmt"

func (an *Analyser) LookFor(match []string) (string, int) {
	for _, style := range match {
		if len(style) == 1 {
			if an.textInput[an.here] == style[0] {
				return style, 1
			} else {
				continue
			}
		} else {
			if an.here+len(style) > len(an.textInput) {
				continue
			} else if string(an.textInput[an.here:an.here+len(style)]) != style {
				continue
			} else {
				return style, len(style)
			}
		}
	}
	return "", 0
}

func (an *Analyser) LookForSpaces() (string, int) {
	pos := an.here
	for ; pos < len(an.textInput); pos++ {
		if an.textInput[pos] != ' ' && an.textInput[pos] != '\t' {
			break
		}
	}
	return an.textInput[an.here:pos], pos - an.here
}

func (an *Analyser) SetQuote(style string) {
	if an.commentContext != nil {
		// Add quote into comment content
		fmt.Println("adding quote into comment content")
		an.commentContext.Content += style
		return
	}
	an.createTextIfNil()
	if an.textContext.QuoteStyle == "" {
		// Begin to quote
		fmt.Println("quote begins here")
		an.textContext.QuoteStyle = style
	} else {
		if an.textContext.QuoteStyle == style {
			// Closing a quote
			fmt.Println("quote ends here")
			an.endText()
		} else {
			// Just content
			fmt.Println("quote is content")
			an.saveMissedCharacters()
			an.textContext.Text += style
		}
	}
}

func (an *Analyser) Analyse() {

	var adv int
	var spaces string
	for an.here = 0; an.here < len(an.textInput); an.here += adv {
		var style string
		if style, adv = an.LookFor(an.config.CommentBeginningMarkers); adv > 0 {
			fmt.Println("Comment: " + style)
			an.createCommentIfNil(style)
			an.lastBranchPosition = an.here + adv
		} else if style, adv = an.LookFor(an.config.TextQuoteStyle); adv > 0 {
			fmt.Println("Quote: " + style)
			an.SetQuote(style)
			an.lastBranchPosition = an.here + adv
		} else if spaces, adv = an.LookForSpaces(); adv > 0 {
			fmt.Println("Spaces: ", adv, spaces)
			an.saveSpaces(spaces)
			an.lastBranchPosition = an.here + adv
		} else if style, adv = an.LookFor(an.config.StatementContinuationMarkers); adv > 0 {
			fmt.Println("StmtContinue: " + style)
			an.ContinueStmt(style)
			an.lastBranchPosition = an.here + adv
		} else if style, adv = an.LookFor(an.config.StatementEndingMarkers); adv > 0 {
			fmt.Println("StmtEnd: " + style)
			an.endStatement(style)
			an.lastBranchPosition = an.here + adv
		} else if style, adv = an.LookFor(an.config.SectionEndingSuffixes); adv > 0 {
			fmt.Println("SectEndSuffix: " + style)
			an.EndSectionSetSuffix(style)
			an.lastBranchPosition = an.here + adv
		} else if style, adv = an.LookFor(an.config.SectionEndingPrefixes); adv > 0 {
			fmt.Println("SectEndPrefix: " + style)
			an.EndSectionSetPrefix(style)
			an.lastBranchPosition = an.here + adv
		} else if style, adv = an.LookFor(an.config.SectionBeginningSuffixes); adv > 0 {
			fmt.Println("SectBeginSuffix: " + style)
			an.BeginSectionSetSuffix(style)
			an.lastBranchPosition = an.here + adv
		} else if style, adv = an.LookFor(an.config.SectionBeginningPrefixes); adv > 0 {
			fmt.Println("SectBeginPrefix: " + style)
			an.BeginSectionSetPrefix(style)
			an.lastBranchPosition = an.here + adv
		} else {
			fmt.Println("text '"+string(an.textInput[an.here])+"' does not match any condition", an.lastBranchPosition, an.here, an.commentContext, an.textContext, an.statementContext)
			adv = 1
		}
	}
	fmt.Println("Analyse finished", an.lastBranchPosition, an.here)
	an.saveMissedCharacters()
	fmt.Println("Analyse will end stmt for one last time")
	an.endStatement("")
}
