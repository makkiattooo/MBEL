package mbel

import (
	"bytes"
	"fmt"
)

// Node represents any node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node (e.g. key = value, @meta)
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node (e.g. "value", { block })
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements    []Statement
	AIAnnotations []*AIAnnotation            // Extracted AI_Context, AI_Tone, etc.
	Terms         map[string]*TermDefinition // -term-name definitions
	Imports       []string                   // @import namespaces
}

// AIAnnotation represents structured AI metadata from comments
// # AI_Context: Button on login screen
// # AI_Tone: Motivating, short
type AIAnnotation struct {
	Type  string // "Context", "Tone", "Constraints", "Examples"
	Value string
	Line  int
	// ForKey is set when annotation appears directly before an assignment
	ForKey string
}

func (a *AIAnnotation) String() string {
	return fmt.Sprintf("# AI_%s: %s", a.Type, a.Value)
}

// TermDefinition represents -term-name = "value"
type TermDefinition struct {
	Token Token
	Name  string // without the leading "-"
	Value Expression
}

func (td *TermDefinition) statementNode()       {}
func (td *TermDefinition) TokenLiteral() string { return td.Token.Literal }
func (td *TermDefinition) String() string {
	return fmt.Sprintf("-%s = %s\n", td.Name, td.Value.String())
}

// TermReference represents {-term-name} usage in strings
type TermReference struct {
	Token Token
	Name  string
}

func (tr *TermReference) expressionNode()      {}
func (tr *TermReference) TokenLiteral() string { return tr.Token.Literal }
func (tr *TermReference) String() string       { return "{-" + tr.Name + "}" }

// ImportStatement represents @import namespace
type ImportStatement struct {
	Token     Token
	Namespace string
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	return fmt.Sprintf("@import %s\n", is.Namespace)
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// MetadataStatement represents @key: value
type MetadataStatement struct {
	Token Token // The '@' token
	Key   string
	Value string // e.g., "pl", "1.0"
}

func (ms *MetadataStatement) statementNode()       {}
func (ms *MetadataStatement) TokenLiteral() string { return ms.Token.Literal }
func (ms *MetadataStatement) String() string {
	return fmt.Sprintf("@%s: %s\n", ms.Key, ms.Value)
}

// SectionStatement represents [section_name]
type SectionStatement struct {
	Token Token // The '[' token
	Name  string
}

func (ss *SectionStatement) statementNode()       {}
func (ss *SectionStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SectionStatement) String() string {
	return fmt.Sprintf("[%s]\n", ss.Name)
}

// AssignStatement represents key = value or key(arg) { ... }
type AssignStatement struct {
	Token Token // The IDENT token
	Name  string
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	return fmt.Sprintf("%s = %s\n", as.Name, as.Value.String())
}

// StringLiteral represents a string value "..."
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BlockExpression represents a logic block { [0] => "...", [other] => "..." }
type BlockExpression struct {
	Token    Token  // The '{' token
	Argument string // The variable name, e.g. "n" in count(n)
	Cases    []*BlockCase
}

func (be *BlockExpression) expressionNode()      {}
func (be *BlockExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BlockExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + be.Argument + ") {\n")
	for _, c := range be.Cases {
		out.WriteString(c.String())
	}
	out.WriteString("}")
	return out.String()
}

type BlockCase struct {
	Condition  string // "0", "other", "male", "one", "few", "many"
	Value      string // The resulting string
	IsRange    bool   // true if this is a numeric range [2..4]
	RangeStart int    // Start of range (inclusive)
	RangeEnd   int    // End of range (inclusive)
}

func (bc *BlockCase) String() string {
	return fmt.Sprintf("\t[%s] => \"%s\"\n", bc.Condition, bc.Value)
}
