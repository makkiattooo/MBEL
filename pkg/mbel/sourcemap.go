package mbel

// SourceLocation represents a position in source code
type SourceLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

// SourceMap maps keys to their source locations
type SourceMap map[string]SourceLocation

// BuildSourceMap creates a source map from a parsed program
func BuildSourceMap(p *Program, filename string) SourceMap {
	sm := make(SourceMap)

	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *AssignStatement:
			sm[s.Name] = SourceLocation{
				File:   filename,
				Line:   s.Token.Line,
				Column: s.Token.Column,
			}
		case *MetadataStatement:
			sm["@"+s.Key] = SourceLocation{
				File:   filename,
				Line:   s.Token.Line,
				Column: s.Token.Column,
			}
		}
	}

	return sm
}
