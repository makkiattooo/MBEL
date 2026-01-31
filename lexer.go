package mbel

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '\n':
		tok = newToken(TOKEN_NEWLINE, "", l.line, 0)
		// We leave line increment to next skipWhitespace/readChar or do it here?
		// If we do it here, the token line is the OLD line.
		tok.Line = l.line
		l.line++
		l.column = 0
	case '=':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = newToken(TOKEN_ARROW, string(ch)+string(l.ch), l.line, l.column)
		} else {
			tok = newToken(TOKEN_ASSIGN, string(l.ch), l.line, l.column)
		}
	case '@':
		tok = newToken(TOKEN_AT, string(l.ch), l.line, l.column)
	case '{':
		tok = newToken(TOKEN_LBRACE, string(l.ch), l.line, l.column)
	case '}':
		tok = newToken(TOKEN_RBRACE, string(l.ch), l.line, l.column)
	case '[':
		tok = newToken(TOKEN_LBRACKET, string(l.ch), l.line, l.column)
	case ']':
		tok = newToken(TOKEN_RBRACKET, string(l.ch), l.line, l.column)
	case '(':
		tok = newToken(TOKEN_LPAREN, string(l.ch), l.line, l.column)
	case ')':
		tok = newToken(TOKEN_RPAREN, string(l.ch), l.line, l.column)
	case ':':
		tok = newToken(TOKEN_COLON, string(l.ch), l.line, l.column)
	case ',':
		tok = newToken(TOKEN_COMMA, string(l.ch), l.line, l.column)
	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			tok = newToken(TOKEN_DOT_RANGE, string(ch)+string(l.ch), l.line, l.column)
		} else {
			tok = newToken(TOKEN_ILLEGAL, string(l.ch), l.line, l.column)
		}
	case '"':
		if l.isTripleQuote() {
			tok.Type = TOKEN_STRING
			tok.Literal = l.readTripleQuotedString()
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok.Type = TOKEN_STRING
			tok.Literal = l.readString()
			tok.Line = l.line
			tok.Column = l.column
		}
	case '#':
		tok.Type = TOKEN_COMMENT
		tok.Literal = l.readComment()
		tok.Line = l.line
		tok.Column = l.column
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = TOKEN_IDENT
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = TOKEN_NUMBER
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = newToken(TOKEN_ILLEGAL, string(l.ch), l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	// Skip ' ' \t \r but NOT \n
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for {
		if isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
			l.readChar()
		} else if l.ch == '.' {
			// Check if it's a range ".."
			if l.peekChar() == '.' {
				// It's a range, stop reading identifier
				break
			}
			// It's a single dot inside identifier
			l.readChar()
		} else {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	// Handle decimal point, but NOT range operator (..)
	if l.ch == '.' && l.peekChar() != '.' {
		l.readChar() // consume the dot
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) isTripleQuote() bool {
	if l.ch == '"' && l.readPosition < len(l.input) && l.input[l.readPosition] == '"' && l.readPosition+1 < len(l.input) && l.input[l.readPosition+1] == '"' {
		return true
	}
	return false
}

func (l *Lexer) readTripleQuotedString() string {
	l.readChar()
	l.readChar()
	l.readChar()

	position := l.position
	for {
		if l.ch == '"' && l.readPosition < len(l.input) && l.input[l.readPosition] == '"' && l.readPosition+1 < len(l.input) && l.input[l.readPosition+1] == '"' {
			break
		}
		if l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
	str := l.input[position:l.position]

	l.readChar()
	l.readChar()
	l.readChar()

	return str
}

func (l *Lexer) readComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType TokenType, literal string, line, col int) Token {
	return Token{Type: tokenType, Literal: literal, Line: line, Column: col}
}
