package parser_test

import (
	"testing"

	"github.com/neatflowcv/cfinder/internal/pkg/domain"
	"github.com/neatflowcv/cfinder/internal/pkg/parser"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()

	parser := parser.NewParser()

	symbols := parser.Parse("test.c", []byte(`    if (fp) {
        pclose(fp);
        fp = NULL;
    }`))

	require.Len(t, symbols, 1)
	require.Equal(t, "pclose", symbols[0].Name)
	require.Equal(t, "test.c", symbols[0].Path)
	require.Equal(t, 2, symbols[0].Line)
	require.Equal(t, domain.FunctionCall, symbols[0].Kind)
}
