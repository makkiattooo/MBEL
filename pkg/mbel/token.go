package mbel

import "fmt"

type TokenType string

const (
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
	TOKEN_EOF     TokenType = "EOF"

	// Identifiers & Literals
	TOKEN_IDENT  TokenType = "IDENT"  // key_name
	TOKEN_STRING TokenType = "STRING" // "value", """multiline"""
	TOKEN_NUMBER TokenType = "NUMBER" // 1, 2, 0.5

	// Operators & Delimiters
	TOKEN_ASSIGN    TokenType = "="
	TOKEN_AT        TokenType = "@"
	TOKEN_LBRACE    TokenType = "{"
	TOKEN_RBRACE    TokenType = "}"
	TOKEN_LBRACKET  TokenType = "["
	TOKEN_RBRACKET  TokenType = "]"
	TOKEN_LPAREN    TokenType = "("
	TOKEN_RPAREN    TokenType = ")"
	TOKEN_ARROW     TokenType = "=>"
	TOKEN_COLON     TokenType = ":"
	TOKEN_COMMA     TokenType = ","
	TOKEN_DOT_RANGE TokenType = ".."

	// Special
	TOKEN_COMMENT TokenType = "COMMENT" // # Comment
	TOKEN_NEWLINE TokenType = "NEWLINE"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, '%s') at %d:%d", t.Type, t.Literal, t.Line, t.Column)
}
