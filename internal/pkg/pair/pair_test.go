package pair_test

import (
	"testing"

	"github.com/neatflowcv/cfinder/internal/pkg/pair"
	"github.com/stretchr/testify/require"
)

func TestFindPair(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		content  [][]byte
		start    []byte
		end      []byte
		expected []*pair.Pair
	}{
		{
			name:     "simple",
			content:  [][]byte{[]byte("("), []byte(")"), []byte("("), []byte(")")},
			start:    []byte("("),
			end:      []byte(")"),
			expected: []*pair.Pair{pair.NewPair(0, 1), pair.NewPair(2, 3)},
		},
		{
			name:     "simple2",
			content:  [][]byte{[]byte("("), []byte("asdf"), []byte(")")},
			start:    []byte("("),
			end:      []byte(")"),
			expected: []*pair.Pair{pair.NewPair(0, 2)},
		},
		{
			name:     "comment",
			content:  [][]byte{[]byte("/*"), []byte("*/"), []byte("/*"), []byte("*/")},
			start:    []byte("/*"),
			end:      []byte("*/"),
			expected: []*pair.Pair{pair.NewPair(0, 1), pair.NewPair(2, 3)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ret := pair.FindPair(test.content, test.start, test.end)

			require.Equal(t, test.expected, ret)
		})
	}
}
