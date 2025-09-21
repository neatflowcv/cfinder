package pair

import (
	"bytes"
)

type Pair struct {
	Start int
	End   int
}

func NewPair(start int, end int) *Pair {
	return &Pair{Start: start, End: end}
}

func FindPair(content [][]byte, start, end []byte) []*Pair {
	var (
		ret   []*Pair
		stack []int
	)

	for idx, part := range content {
		switch {
		case bytes.Equal(part, start):
			stack = append(stack, idx)

		case bytes.Equal(part, end):
			if len(stack) == 0 {
				continue
			}

			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			ret = append(ret, NewPair(last, idx))
		}
	}

	return ret
}
