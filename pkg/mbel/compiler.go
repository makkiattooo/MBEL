package mbel

import (
	"fmt"
)

// Compiler transforms AST into a runtime map
type Compiler struct {
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(node Node) (interface{}, error) {
	switch n := node.(type) {
	case *Program:
		return c.compileProgram(n)
	case *MetadataStatement:
		return n, nil // Metadata handled by Program
	case *SectionStatement:
		return n, nil // Sections handled by Program
	case *AssignStatement:
		return c.compileAssign(n)
	case *StringLiteral:
		return n.Value, nil
	case *BlockExpression:
		return c.compileBlock(n)
	case *ImportStatement:
		return n, nil // Imports handled by Program
	case *TermDefinition:
		return n, nil // Terms handled by Program
	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

func (c *Compiler) compileProgram(p *Program) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	metadata := make(map[string]string)

	// First pass to get metadata (especially namespace)
	for _, stmt := range p.Statements {
		if ms, ok := stmt.(*MetadataStatement); ok {
			metadata[ms.Key] = ms.Value
		}
	}

	currentSection := ""

	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *SectionStatement:
			currentSection = s.Name
		case *AssignStatement:
			val, err := c.Compile(s)
			if err != nil {
				return nil, err
			}

			key := s.Name
			if currentSection != "" {
				key = currentSection + "." + s.Name
			}
			result[key] = val
		}
	}

	if len(metadata) > 0 {
		result["__meta"] = metadata
	}

	// Export AI annotations
	if len(p.AIAnnotations) > 0 {
		aiMap := make(map[string][]map[string]string)
		for _, ann := range p.AIAnnotations {
			entry := map[string]string{
				"type":  ann.Type,
				"value": ann.Value,
			}
			if ann.ForKey != "" {
				aiMap[ann.ForKey] = append(aiMap[ann.ForKey], entry)
			} else {
				aiMap["__global"] = append(aiMap["__global"], entry)
			}
		}
		result["__ai"] = aiMap
	}

	// Export imports
	if len(p.Imports) > 0 {
		result["__imports"] = p.Imports
	}

	// Export terms
	if len(p.Terms) > 0 {
		terms := make(map[string]string)
		for name, def := range p.Terms {
			if sl, ok := def.Value.(*StringLiteral); ok {
				terms[name] = sl.Value
			}
		}
		result["__terms"] = terms
	}

	return result, nil
}

func (c *Compiler) compileAssign(node *AssignStatement) (interface{}, error) {
	return c.Compile(node.Value)
}

// RangeCase represents a compiled numeric range condition
type RangeCase struct {
	Start int
	End   int
	Value string
}

// RuntimeBlock represents a compiled logic block ready for execution
type RuntimeBlock struct {
	Argument   string
	Cases      map[string]string // keyword conditions: "one", "other", "0"
	RangeCases []RangeCase       // numeric range conditions: [2..4]
}

// Resolve finds the matching value for given argument
func (rb *RuntimeBlock) Resolve(arg interface{}) string {
	// Try string match first
	if strArg, ok := arg.(string); ok {
		if val, exists := rb.Cases[strArg]; exists {
			return val
		}
	}

	// Try numeric match
	var numArg int
	switch v := arg.(type) {
	case int:
		numArg = v
	case int64:
		numArg = int(v)
	case float64:
		numArg = int(v)
	default:
		// Fall through to other
	}

	// Check exact number match
	numStr := fmt.Sprintf("%d", numArg)
	if val, exists := rb.Cases[numStr]; exists {
		return val
	}

	// Check range matches
	for _, rc := range rb.RangeCases {
		if numArg >= rc.Start && numArg <= rc.End {
			return rc.Value
		}
	}

	// Check plural categories (hardcoded PL/EN)
	pluralCat := ResolvePluralCategory("pl", numArg)
	if val, exists := rb.Cases[pluralCat]; exists {
		return val
	}

	// Fallback to "other"
	if val, exists := rb.Cases["other"]; exists {
		return val
	}

	return ""
}

// ResolvePluralCategory returns CLDR plural category for a number
// Hardcoded for Polish and English
func ResolvePluralCategory(lang string, n int) string {
	switch lang {
	case "pl":
		// Polish rules:
		// one: n == 1
		// few: n % 10 in 2..4 AND n % 100 NOT in 12..14
		// many: n != 1 AND n % 10 in 0..1 OR n % 10 in 5..9 OR n % 100 in 12..14
		// other: fractions (not handled here)
		if n == 1 {
			return "one"
		}
		mod10 := n % 10
		mod100 := n % 100
		if mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14) {
			return "few"
		}
		return "many"

	case "en":
		// English rules:
		// one: n == 1
		// other: everything else
		if n == 1 {
			return "one"
		}
		return "other"

	default:
		// Default to simple one/other
		if n == 1 {
			return "one"
		}
		return "other"
	}
}

func (c *Compiler) compileBlock(node *BlockExpression) (*RuntimeBlock, error) {
	rb := &RuntimeBlock{
		Argument:   node.Argument,
		Cases:      make(map[string]string),
		RangeCases: []RangeCase{},
	}

	for _, bc := range node.Cases {
		if bc.IsRange {
			rb.RangeCases = append(rb.RangeCases, RangeCase{
				Start: bc.RangeStart,
				End:   bc.RangeEnd,
				Value: bc.Value,
			})
		} else {
			rb.Cases[bc.Condition] = bc.Value
		}
	}

	return rb, nil
}
