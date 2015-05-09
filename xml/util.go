package xml // import "github.com/tdewolff/parse/xml"

import "github.com/tdewolff/parse"

var (
	ltEntityBytes          = []byte("&lt;")
	ampEntityBytes         = []byte("&amp;")
	singleQuoteEntityBytes = []byte("&#39;")
	doubleQuoteEntityBytes = []byte("&#34;")
)

// EscapeAttrVal returns the escape attribute value bytes without quotes.
func EscapeAttrVal(buf *[]byte, b []byte) []byte {
	singles := 0
	doubles := 0
	for i, c := range b {
		if c == '&' {
			if quote, _, ok := parse.QuoteEntity(b[i:]); ok {
				if quote == '"' {
					doubles++
				} else {
					singles++
				}
			}
		} else if c == '"' {
			doubles++
		} else if c == '\'' {
			singles++
		}
	}

	var quote byte
	var escapedQuote []byte
	if doubles > singles {
		quote = '\''
		escapedQuote = singleQuoteEntityBytes
	} else {
		quote = '"'
		escapedQuote = doubleQuoteEntityBytes
	}
	if len(b)+2 > cap(*buf) {
		*buf = make([]byte, 0, len(b)+2) // maximum size, not actual size
	}
	t := (*buf)[:len(b)+2] // maximum size, not actual size
	t[0] = quote
	j := 1
	start := 0
	for i, c := range b {
		if c == '&' {
			if entityQuote, n, ok := parse.QuoteEntity(b[i:]); ok {
				j += copy(t[j:], b[start:i])
				if entityQuote != quote {
					j += copy(t[j:], []byte{entityQuote})
				} else {
					j += copy(t[j:], escapedQuote)
				}
				start = i + n
			}
		} else if c == quote {
			j += copy(t[j:], b[start:i])
			j += copy(t[j:], escapedQuote)
			start = i + 1
		}
	}
	j += copy(t[j:], b[start:])
	t[j] = quote
	return t[:j+1]
}

// EscapeCDATAVal returns the escaped text bytes.
func EscapeCDATAVal(buf *[]byte, b []byte) ([]byte, bool) {
	n := 0
	for _, c := range b {
		if c == '<' || c == '&' {
			if c == '<' {
				n += 3 // &lt;
			} else {
				n += 4 // &amp;
			}
			if n > len("<![CDATA[]]>") {
				return b, true
			}
		}
	}
	if len(b)+n > cap(*buf) {
		*buf = make([]byte, 0, len(b)+n)
	}
	t := (*buf)[:len(b)+n]
	j := 0
	start := 0
	for i, c := range b {
		if c == '<' {
			j += copy(t[j:], b[start:i])
			j += copy(t[j:], ltEntityBytes)
			start = i + 1
		} else if c == '&' {
			j += copy(t[j:], b[start:i])
			j += copy(t[j:], ampEntityBytes)
			start = i + 1
		}
	}
	j += copy(t[j:], b[start:])
	return t[:j], false
}
