package domain

type Symbol struct {
	Kind SymbolKind
	Name string
	Path string
	Line int
}

type SymbolKind int

const (
	FunctionDefinition SymbolKind = iota + 1
	FunctionDeclaration
	FunctionCall
)

func NewSymbol(kind SymbolKind, name string, path string, line int) *Symbol {
	return &Symbol{Path: path, Name: name, Kind: kind, Line: line}
}
