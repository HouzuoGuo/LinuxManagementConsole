package txtedit

import "fmt"

func (an *Analyser) LookFor(match []string) (string, int) {
	for _, style := range match {
		if len(style) == 1 {
			if an.textInput[an.positionHere] == style[0] {
				return style, 1
			} else {
				continue
			}
		} else {
			if an.positionHere+len(style) > len(an.textInput) {
				continue
			} else if string(an.textInput[an.positionHere:an.positionHere+len(style)]) != style {
				continue
			} else {
				return style, len(style)
			}
		}
	}
	return "", 0
}

func (an *Analyser) LookForSpaces() (string, int) {
	pos := an.positionHere
	for ; pos < len(an.textInput); pos++ {
		if an.textInput[pos] != ' ' && an.textInput[pos] != '\t' {
			break
		}
	}
	return an.textInput[an.positionHere:pos], pos - an.positionHere
}

func (an *Analyser) SetQuote(style string) {
	if an.contextComment != nil {
		// Add quote into comment content
		fmt.Println("adding quote into comment content")
		an.contextComment.Content += style
		return
	}
	an.createTextIfNil()
	if an.contextText.QuoteStyle == "" {
		// Begin to quote
		fmt.Println("quote begins here")
		an.contextText.QuoteStyle = style
	} else {
		if an.contextText.QuoteStyle == style {
			// Closing a quote
			fmt.Println("quote ends here")
			an.endText()
		} else {
			// Just content
			fmt.Println("quote is content")
			an.saveMissedCharacters()
			an.contextText.Text += style
		}
	}
}

func (an *Analyser) Analyse() {
	var advance int
	var spaces string
	for an.positionHere = 0; an.positionHere < len(an.textInput); an.positionHere += advance {
		var style string
		if style, advance = an.LookFor(an.config.CommentBeginningMarkers); advance > 0 {
			fmt.Println("Comment: " + style)
			an.createCommentIfNil(style)
			an.positionLastBranch = an.positionHere + advance
		} else if style, advance = an.LookFor(an.config.TextQuoteStyle); advance > 0 {
			fmt.Println("Quote: " + style)
			an.SetQuote(style)
			an.positionLastBranch = an.positionHere + advance
		} else if spaces, advance = an.LookForSpaces(); advance > 0 {
			fmt.Println("Spaces: ", advance, spaces)
			an.saveSpaces(spaces)
			an.positionLastBranch = an.positionHere + advance
		} else if style, advance = an.LookFor(an.config.StatementContinuationMarkers); advance > 0 {
			fmt.Println("StmtContinue: " + style)
			an.continueStatement(style)
			an.positionLastBranch = an.positionHere + advance
		} else if style, advance = an.LookFor(an.config.StatementEndingMarkers); advance > 0 {
			fmt.Println("StmtEnd: " + style)
			an.endStatement(style)
			an.positionLastBranch = an.positionHere + advance
		} else if style, advance = an.LookFor(an.config.SectionEndingSuffixes); advance > 0 {
			fmt.Println("SectEndSuffix: " + style)
			an.setSectionEndSuffix(style)
			an.positionLastBranch = an.positionHere + advance
		} else if style, advance = an.LookFor(an.config.SectionEndingPrefixes); advance > 0 {
			fmt.Println("SectEndPrefix: " + style)
			an.setSectionEndPrefix(style)
			an.positionLastBranch = an.positionHere + advance
		} else if style, advance = an.LookFor(an.config.SectionBeginningSuffixes); advance > 0 {
			fmt.Println("SectBeginSuffix: " + style)
			an.setSectionBeginSuffix(style)
			an.positionLastBranch = an.positionHere + advance
		} else if style, advance = an.LookFor(an.config.SectionBeginningPrefixes); advance > 0 {
			fmt.Println("SectBeginPrefix: " + style)
			an.setSectionBeginPrefix(style)
			an.positionLastBranch = an.positionHere + advance
		} else {
			advance = 1
		}
	}
	an.debug.Printf("Analyse: will end statement for one last time")
	an.endStatement("")
}
