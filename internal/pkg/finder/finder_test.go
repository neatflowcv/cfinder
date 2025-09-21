package finder_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/neatflowcv/cfinder/internal/pkg/finder"
	"github.com/stretchr/testify/require"
)

//go:embed test_file.c.test
var testFileContent []byte

func TestFinder_FindSymbol(t *testing.T) {
	t.Parallel()

	finder := finder.NewFinder(bytes.NewReader(testFileContent))

	symbols := finder.FindSymbol("abc_def_ghijkl.c", "abc_def_ghijkl")

	require.Len(t, symbols, 1)
	require.Equal(t, "abc_def_ghijkl", symbols[0].Name)
	require.Equal(t, "abc_def_ghijkl.c", symbols[0].Path)
	require.Equal(t, 4, symbols[0].Line)
}
