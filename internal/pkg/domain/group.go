package domain

type Group struct {
	Definition  *Symbol
	Declaration *Symbol
	Calls       []*Symbol
}

func NewGroups(symbols []*Symbol) []*Group {
	symbolGroup := groupBy(symbols, func(s *Symbol) string {
		return s.Name
	})

	var ret []*Group

	for _, group := range symbolGroup {
		var (
			definition  *Symbol
			declaration *Symbol
			calls       []*Symbol
		)

		for _, symbol := range group {
			switch symbol.Kind {
			case FunctionDefinition:
				definition = symbol
			case FunctionDeclaration:
				declaration = symbol
			case FunctionCall:
				calls = append(calls, symbol)
			}
		}

		ret = append(ret, &Group{
			Definition:  definition,
			Declaration: declaration,
			Calls:       calls,
		})
	}

	return ret
}

func groupBy[T any, K comparable](s []T, keyFunc func(T) K) map[K][]T {
	groups := make(map[K][]T)

	for _, item := range s {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)
	}

	return groups
}
