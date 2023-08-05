package unicode

import "strings"

func Normalize(text string) string {
	fixedText := text
	for _, rule := range rules {
		fixedText = replaceUnicodeChar(fixedText, rule.originalChar, rule.replaceChar)
	}
	return fixedText
}

func replaceUnicodeChar(str string, oldChar, newChar rune) string {
	return strings.Map(func(r rune) rune {
		if r == oldChar {
			return newChar
		}
		return r
	}, str)
}
