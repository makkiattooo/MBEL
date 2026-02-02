package mbel

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `@lang: pl
@version: 1.0

# Context: Main Title
title = "Hello World"
count(n) {
	[other] => "{n} items"
}

dotted.key.test = "Value"

description = """
Line 1
Line 2
"""
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TOKEN_AT, "@"},
		{TOKEN_IDENT, "lang"},
		{TOKEN_COLON, ":"},
		{TOKEN_IDENT, "pl"},
		{TOKEN_NEWLINE, ""},

		{TOKEN_AT, "@"},
		{TOKEN_IDENT, "version"},
		{TOKEN_COLON, ":"},
		{TOKEN_NUMBER, "1.0"},
		{TOKEN_NEWLINE, ""},

		{TOKEN_NEWLINE, ""}, // Empty line
		{TOKEN_COMMENT, " Context: Main Title"},

		{TOKEN_IDENT, "title"},
		{TOKEN_ASSIGN, "="},
		{TOKEN_STRING, "Hello World"},
		{TOKEN_NEWLINE, ""},

		{TOKEN_IDENT, "count"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "n"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_LBRACE, "{"},
		{TOKEN_NEWLINE, ""},

		{TOKEN_LBRACKET, "["},
		{TOKEN_IDENT, "other"},
		{TOKEN_RBRACKET, "]"},
		{TOKEN_ARROW, "=>"},
		{TOKEN_STRING, "{n} items"},
		{TOKEN_NEWLINE, ""},

		{TOKEN_RBRACE, "}"},
		{TOKEN_NEWLINE, ""},
		{TOKEN_NEWLINE, ""},

		{TOKEN_IDENT, "dotted.key.test"},
		{TOKEN_ASSIGN, "="},
		{TOKEN_STRING, "Value"},
		{TOKEN_NEWLINE, ""},
		{TOKEN_NEWLINE, ""},

		{TOKEN_IDENT, "description"},
		{TOKEN_ASSIGN, "="},
		{TOKEN_STRING, "\nLine 1\nLine 2\n"},
		{TOKEN_NEWLINE, ""},

		{TOKEN_EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		// Skip NEWLINE tokens if my logic implies they are skipped whitepace
		// Current logic: skipWhitespace consumes \n. So NextToken never returns TOKEN_NEWLINE.
		// So I should remove TOKEN_NEWLINE from expected tests.

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
