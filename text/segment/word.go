package segment

//go:generate go run gen_props.go --prefix wb --dataPath data/WordBreakProperty.txt --propertyName CR --propertyName LF --propertyName Newline --propertyName ZWJ --propertyName WSegSpace --propertyName Extend --propertyName Format --propertyName ALetter --propertyName Hebrew_Letter --propertyName MidLetter --propertyName MidNum --propertyName Single_Quote --propertyName Double_Quote --propertyName Numeric --propertyName Katakana --propertyName ExtendNumLet --propertyName Regional_Indicator --outputPath word_props.go

// Information symbol (U+2139) is both ALetter and Extended_Pictographic, so we need a separate table for Extended_Pictographic.
//go:generate go run gen_props.go --prefix wbe --dataPath data/emoji-data.txt --propertyName Extended_Pictographic --outputPath word_emoji_props.go

// WordBreaker finds breakpoints between words.
// This complies with https://www.unicode.org/reports/tr29/#Word_Boundaries
type WordBreaker struct {
	lastProp wbProp
}

func (wb *WordBreaker) ProcessRune(r rune) (canBreakBefore bool) {
	prop := wbPropForRune(r)

	// WB1 sot ÷ Any
	// WB2 Any ÷ eot
	// We don't need to implement these because we're only interested in non-empty segments.

	// WB3 CR × LF
	if prop == wbPropLF && wb.lastProp == wbPropCR {
		canBreakBefore = false
		goto done
	}

	// WB3a (Newline | CR | LF) ÷
	if wb.lastProp == wbPropNewline || wb.lastProp == wbPropCR || wb.lastProp == wbPropLF {
		canBreakBefore = true
		goto done
	}

	// WB3b   ÷ (Newline | CR | LF)

	/*
		Otherwise break before and after Newlines (including CR and LF)
		WB3a (Newline | CR | LF) ÷
		WB3b   ÷ (Newline | CR | LF)

		Do not break within emoji zwj sequences.
		WB3c ZWJ × \p{Extended_Pictographic}

		Keep horizontal whitespace together.
		WB3d WSegSpace × WSegSpace

		Ignore Format and Extend characters, except after sot, CR, LF, and Newline. (See Section 6.2, Replacing Ignore Rules.) This also has the effect of: Any × (Format | Extend | ZWJ)
		WB4 X (Extend | Format | ZWJ)* → X

		Do not break between most letters.
		WB5 AHLetter × AHLetter

		Do not break letters across certain punctuation.
		WB6 AHLetter × (MidLetter | MidNumLetQ) AHLetter
		WB7 AHLetter (MidLetter | MidNumLetQ) × AHLetter
		WB7a Hebrew_Letter × Single_Quote
		WB7b Hebrew_Letter × Double_Quote Hebrew_Letter
		WB7c Hebrew_Letter Double_Quote × Hebrew_Letter

		Do not break within sequences of digits, or digits adjacent to letters (“3a”, or “A3”).
		WB8 Numeric × Numeric
		WB9 AHLetter × Numeric
		WB10 Numeric × AHLetter

		Do not break within sequences, such as “3.2” or “3,456.789”.
		WB11 Numeric (MidNum | MidNumLetQ) × Numeric
		WB12 Numeric × (MidNum | MidNumLetQ) Numeric

		Do not break between Katakana.
		WB13 Katakana × Katakana

		Do not break from extenders.
		WB13a (AHLetter | Numeric | Katakana | ExtendNumLet) × ExtendNumLet
		WB13b ExtendNumLet × (AHLetter | Numeric | Katakana)

		Do not break within emoji flag sequences. That is, do not break between regional indicator (RI) symbols if there is an odd number of RI characters before the break point.
		WB15 sot (RI RI)* RI × RI
		WB16 [^RI] (RI RI)* RI × RI

		Otherwise, break everywhere (including around ideographs).
		WB999 Any ÷ Any
	*/

done:
	wb.lastProp = prop

	return canBreakBefore
}
