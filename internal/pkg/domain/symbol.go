package domain

type Symbol struct {
	Name string
	Kind SymbolKind
	Line int
}

type SymbolKind int

const (
	FunctionDefinition SymbolKind = iota + 1
	FunctionDeclaration
	FunctionCall
)

func NewSymbol(name string, kind SymbolKind, line int) *Symbol {
	return &Symbol{Name: name, Kind: kind, Line: line}
}
