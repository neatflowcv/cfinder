package split_test

import (
	"testing"

	"github.com/neatflowcv/cfinder/internal/pkg/split"
	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      []byte
		delimiters [][]byte
		expected   [][]byte
	}{
		{
			name:       "space",
			input:      []byte("abc def ghi"),
			delimiters: [][]byte{[]byte(" ")},
			expected:   [][]byte{[]byte("abc"), []byte(" "), []byte("def"), []byte(" "), []byte("ghi")},
		},
		{
			name:       "hash",
			input:      []byte("#"),
			delimiters: [][]byte{[]byte("#")},
			expected:   [][]byte{[]byte("#")},
		},
		{
			name:       "comment",
			input:      []byte("/***"),
			delimiters: [][]byte{[]byte("/*")},
			expected:   [][]byte{[]byte("/*"), []byte("**")},
		},
		{
			name:       "enter",
			input:      []byte("\n\n\n"),
			delimiters: [][]byte{[]byte("\n")},
			expected:   [][]byte{[]byte("\n"), []byte("\n"), []byte("\n")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, split.Split(test.input, test.delimiters))
		})
	}
}

func TestSplit_example(t *testing.T) {
	t.Parallel()

	content := []byte(`#include <stdio.h>

int add(int a, int b) {
    return a + b;
}

int main(void) {
    int x = 5, y = 7;
    int result = add(x, y);

    printf("%d + %d = %d\n", x, y, result);
    return 0;
}`)

	ret := split.Split(content, [][]byte{[]byte("{}();#, \t\n")})

	require.NotEmpty(t, ret)
}
