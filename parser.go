package mbel

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	l                    *Lexer
	curToken             Token
	peekToken            Token
	errors               []string
	pendingAIAnnotations []*AIAnnotation // AI annotations waiting to be attached to next key
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// Extract AI annotations from comments, skip other comments
	for p.curToken.Type == TOKEN_COMMENT {
		if ann := p.parseAIAnnotation(p.curToken); ann != nil {
			p.pendingAIAnnotations = append(p.pendingAIAnnotations, ann)
		}
		p.curToken = p.peekToken
		p.peekToken = p.l.NextToken()
	}
}

// parseAIAnnotation extracts structured data from AI_* comments
// Supports formats like:
// @AI_Context: "some text"
// @AI_MaxLength: 50
// Returns nil if comment is not an AI annotation
func (p *Parser) parseAIAnnotation(tok Token) *AIAnnotation {
	text := strings.TrimSpace(tok.Literal)

	// Check for AI_ prefix
	if !strings.HasPrefix(text, "AI_") {
		return nil
	}

	// Find the colon separator
	colonIdx := strings.Index(text, ":")
	if colonIdx == -1 {
		return nil
	}

	// Extract type (Context, Tone, Constraints, Examples)
	aiType := strings.TrimPrefix(text[:colonIdx], "AI_")
	value := strings.TrimSpace(text[colonIdx+1:])

	// Handle multi-line values in curly braces
	if strings.HasPrefix(value, "{") {
		// Value spans multiple lines until closing }
		// Collect lines until we find the closing }
		lines := []string{value}
		bracketCount := 1

		for bracketCount > 0 && p.peekToken.Type != TOKEN_EOF {
			p.nextToken() // Move to next token
			line := p.curToken.Literal

			// Count brackets
			for _, ch := range line {
				if ch == '{' {
					bracketCount++
				} else if ch == '}' {
					bracketCount--
				}
			}

			lines = append(lines, line)
			if bracketCount <= 0 {
				break
			}
		}

		// Join all lines and parse as JSON/YAML-like
		fullValue := strings.Join(lines, "\n")
		value = parseMultiLineAnnotation(fullValue)
	}

	return &AIAnnotation{
		Type:  aiType,
		Value: value,
		Line:  tok.Line,
	}
}

// parseMultiLineAnnotation handles multi-line annotation values
// Supports both JSON and simple YAML-like formats
func parseMultiLineAnnotation(value string) string {
	// Remove outer braces
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		value = strings.TrimSpace(value[1 : len(value)-1])
	}

	// For now, return as-is. Could be extended to parse YAML/JSON
	// This allows storing structured data
	return value
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{
		Terms: make(map[string]*TermDefinition),
	}
	program.Statements = []Statement{}

	for p.curToken.Type != TOKEN_EOF {
		stmt := p.parseStatement(program)
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else if p.curToken.Type != TOKEN_NEWLINE && p.curToken.Type != TOKEN_EOF {
			// If statement parsing failed and it wasn't just an empty line,
			// we need to skip to the next safe point
			p.synchronize()
		}
		p.nextToken()
	}

	// Add any remaining pending annotations
	program.AIAnnotations = append(program.AIAnnotations, p.pendingAIAnnotations...)
	p.pendingAIAnnotations = nil

	return program
}

// synchronize skips tokens until a safe state (statement boundary) is found
// Used for error recovery
func (p *Parser) synchronize() {
	for p.curToken.Type != TOKEN_EOF {
		if p.curToken.Type == TOKEN_NEWLINE {
			return
		}

		// Check if next token starts a statement
		switch p.peekToken.Type {
		case TOKEN_IDENT, TOKEN_AT, TOKEN_LBRACKET:
			return
		}

		p.nextToken()
	}
}

func (p *Parser) parseStatement(program *Program) Statement {
	switch p.curToken.Type {
	case TOKEN_AT:
		return p.parseMetadataOrImport(program)
	case TOKEN_IDENT:
		stmt := p.parseAssignStatement(program)
		if stmt == nil {
			return nil
		}
		return stmt
	case TOKEN_LBRACKET:
		stmt := p.parseSectionStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	case TOKEN_NEWLINE:
		return nil // Skip empty lines / separators
	default:
		return nil
	}
}

func (p *Parser) parseSectionStatement() *SectionStatement {
	stmt := &SectionStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_IDENT) {
		return nil
	}
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(TOKEN_RBRACKET) {
		return nil
	}

	return stmt
}

// parseMetadataOrImport handles both @key: value and @import namespace
func (p *Parser) parseMetadataOrImport(program *Program) Statement {
	startToken := p.curToken

	if !p.expectPeek(TOKEN_IDENT) {
		return nil
	}

	key := p.curToken.Literal

	// Check for @import directive
	if key == "import" {
		p.nextToken()
		if p.curToken.Type == TOKEN_IDENT {
			program.Imports = append(program.Imports, p.curToken.Literal)
			return &ImportStatement{Token: startToken, Namespace: p.curToken.Literal}
		}
		p.peekError(TOKEN_IDENT)
		return nil
	}

	// Regular metadata @key: value
	stmt := &MetadataStatement{Token: startToken, Key: key}

	if !p.expectPeek(TOKEN_COLON) {
		return nil
	}

	p.nextToken()
	// Metadata value can be IDENT (e.g. pl) or NUMBER (e.g. 1.0) or STRING
	if p.curToken.Type == TOKEN_IDENT || p.curToken.Type == TOKEN_NUMBER || p.curToken.Type == TOKEN_STRING {
		stmt.Value = p.curToken.Literal
	} else {
		p.peekError(TOKEN_STRING)
		return nil
	}

	return stmt
}

func (p *Parser) parseAssignStatement(program *Program) *AssignStatement {
	stmt := &AssignStatement{Token: p.curToken}
	stmt.Name = p.curToken.Literal

	// Attach any pending AI annotations to this key
	if len(p.pendingAIAnnotations) > 0 {
		for _, ann := range p.pendingAIAnnotations {
			ann.ForKey = stmt.Name
		}
		program.AIAnnotations = append(program.AIAnnotations, p.pendingAIAnnotations...)
		p.pendingAIAnnotations = nil
	}

	if p.peekTokenIs(TOKEN_ASSIGN) {
		p.nextToken() // move to =
		p.nextToken() // move to value
		stmt.Value = p.parseExpression()
		if stmt.Value == nil {
			p.errors = append(p.errors, fmt.Sprintf("Expected expression after = at line %d", p.curToken.Line))
			return nil
		}
		return stmt
	} else if p.peekTokenIs(TOKEN_LPAREN) {
		// handle block: key(arg) { ... }
		return p.parseBlockAssignStatement(stmt)
	} else {
		p.peekError(TOKEN_ASSIGN) // or LPAREN
		return nil
	}
}

func (p *Parser) parseBlockAssignStatement(stmt *AssignStatement) *AssignStatement {
	// Current is IDENT (name)
	// Peek is LPAREN
	p.nextToken() // move to LPAREN

	if !p.expectPeek(TOKEN_IDENT) {
		return nil
	}
	argName := p.curToken.Literal

	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}

	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}

	block := &BlockExpression{Token: p.curToken, Argument: argName}
	block.Cases = p.parseBlockCases()

	stmt.Value = block
	return stmt
}

func (p *Parser) parseBlockCases() []*BlockCase {
	cases := []*BlockCase{}

	for !p.peekTokenIs(TOKEN_RBRACE) && !p.peekTokenIs(TOKEN_EOF) {
		p.nextToken()

		if p.curToken.Type == TOKEN_NEWLINE {
			continue
		}

		if p.curToken.Type == TOKEN_LBRACKET {
			// [condition] => "value" or [2..4] => "value"
			bc := &BlockCase{}

			p.nextToken() // move to condition start

			if p.curToken.Type == TOKEN_NUMBER {
				startNum := p.curToken.Literal

				// Check for range [2..4]
				if p.peekTokenIs(TOKEN_DOT_RANGE) {
					p.nextToken() // consume ..
					if !p.expectPeek(TOKEN_NUMBER) {
						return nil
					}
					endNum := p.curToken.Literal

					// Parse as range
					start, err1 := strconv.Atoi(startNum)
					end, err2 := strconv.Atoi(endNum)
					if err1 != nil || err2 != nil {
						p.errors = append(p.errors, fmt.Sprintf("Invalid range numbers at line %d", p.curToken.Line))
						return nil
					}

					bc.IsRange = true
					bc.RangeStart = start
					bc.RangeEnd = end
					bc.Condition = fmt.Sprintf("%d..%d", start, end)
				} else {
					// Simple number condition
					bc.Condition = startNum
				}
			} else if p.curToken.Type == TOKEN_IDENT {
				// Keyword conditions: one, few, many, other, male, female, etc.
				bc.Condition = p.curToken.Literal
			} else {
				p.errors = append(p.errors, fmt.Sprintf("Expected condition at line %d, got %s", p.curToken.Line, p.curToken.Type))
				return nil
			}

			if !p.expectPeek(TOKEN_RBRACKET) {
				return nil
			}

			if !p.expectPeek(TOKEN_ARROW) {
				return nil
			}

			if !p.expectPeek(TOKEN_STRING) {
				return nil
			}
			bc.Value = p.curToken.Literal
			cases = append(cases, bc)
		}
	}

	if !p.expectPeek(TOKEN_RBRACE) {
		return nil
	}

	return cases
}

func (p *Parser) parseExpression() Expression {
	// Simple string literal expression
	if p.curToken.Type == TOKEN_STRING {
		return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	}
	// TODO: Support Number literal as value?
	return nil
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d", t, p.peekToken.Type, p.peekToken.Line)
	p.errors = append(p.errors, msg)
}
