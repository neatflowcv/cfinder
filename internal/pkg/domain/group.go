package domain

type Group struct {
	Definition   *Symbol
	Declarations []*Symbol
	Calls        []*Symbol
}

func NewGroups(symbols []*Symbol) []*Group {
	symbolGroup := groupBy(symbols, func(s *Symbol) string {
		return s.Name
	})

	var ret []*Group

	for _, group := range symbolGroup {
		var (
			definitions  []*Symbol
			declarations []*Symbol
			calls        []*Symbol
		)

		for _, symbol := range group {
			switch symbol.Kind {
			case FunctionDefinition:
				definitions = append(definitions, symbol)
			case FunctionDeclaration:
				declarations = append(declarations, symbol)
			case FunctionCall:
				calls = append(calls, symbol)
			}
		}

		if len(definitions) > 0 && len(declarations) == 0 && len(calls) == 0 {
			// define 함수일 수 있다.
			continue
		}

		if len(definitions) > 2 { //nolint:mnd
			panic("definitions > 2")
		}

		ret = append(ret, &Group{
			Definition:   getFirst(definitions),
			Declarations: declarations,
			Calls:        calls,
		})
	}

	return ret
}

func getFirst(s []*Symbol) *Symbol {
	if len(s) == 0 {
		return nil
	}

	return s[0]
}

func groupBy[T any, K comparable](s []T, keyFunc func(T) K) map[K][]T {
	groups := make(map[K][]T)

	for _, item := range s {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)
	}

	return groups
}
