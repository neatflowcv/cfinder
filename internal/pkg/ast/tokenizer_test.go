package ast_test

import (
	"bytes"
	"testing"

	"github.com/neatflowcv/cfinder/internal/pkg/ast"
	"github.com/stretchr/testify/require"
)

func TestTokenizer_main(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte(`
		int main() {
			return 0;
		}
	`))

	tokens := ast.Tokenize(reader)

	require.Len(t, tokens, 13)
}

func TestTokenizer_preprocessor(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte(`
#define LOG(...)                                                               \
    {                                                                          \
        syslog(__VA_ARGS__);                                                   \
    }
	`))

	tokens := ast.Tokenize(reader)

	require.Len(t, tokens, 21)
}

func TestTokenizer_comment(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte(`//int main() {
	/*int main() {
		return 0;
	}*/ int main() {
		return 0;
	}
	`))

	tokens := ast.Tokenize(reader)

	require.Len(t, tokens, 32)
	require.Equal(t, ast.TokenKindSingleLineComment, tokens[0].Kind)
	require.Equal(t, ast.TokenKindOpenMultiLineComment, tokens[7].Kind)
	require.Equal(t, ast.TokenKindCloseMultiLineComment, tokens[19].Kind)
}
