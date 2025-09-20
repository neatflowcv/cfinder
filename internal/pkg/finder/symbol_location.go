package finder

type SymbolKind int

const (
	FunctionDefinition SymbolKind = iota + 1
	FunctionDeclaration
	FunctionCall
)

type SymbolLocation struct {
	Name string
	Kind SymbolKind
	Line int
}

func NewSymbolLocation(name string, kind SymbolKind, line int) *SymbolLocation {
	return &SymbolLocation{Name: name, Kind: kind, Line: line}
}
