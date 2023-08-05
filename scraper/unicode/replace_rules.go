package unicode

type replaceRule struct {
	originalChar rune
	replaceChar  rune
}

var rules = []replaceRule{
	{
		originalChar: '\u00A0',
		replaceChar:  ' ',
	},
	{
		originalChar: '\u2013',
		replaceChar:  '-',
	},
}
