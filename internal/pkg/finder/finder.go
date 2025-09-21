package finder

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	"github.com/neatflowcv/cfinder/internal/pkg/domain"
)

type Finder struct {
	reader *bufio.Reader
}

func NewFinder(reader io.Reader) *Finder {
	return &Finder{reader: bufio.NewReader(reader)}
}

func (f *Finder) FindSymbol(path string, symbol string) []*domain.Symbol { //nolint:cyclop
	iter := 0

	var locations []*domain.Symbol

	inComment := false

	for {
		part, err := f.reader.ReadBytes('\n')
		if err != nil {
			break
		}

		iter++

		pos := bytes.Index(part, []byte("//"))
		if pos != -1 {
			part = part[:pos]
		}

		pos = bytes.Index(part, []byte("/*"))
		if !inComment && pos != -1 {
			inComment = true
			part = part[:pos]
		}

		pos = bytes.Index(part, []byte("*/"))
		if inComment && pos != -1 {
			inComment = false
			part = part[pos+2:]
		}

		part = bytes.TrimSpace(part)
		if len(part) == 0 {
			continue
		}

		if inComment {
			continue
		}

		// 심볼 앞은 공백이거나 없어야 하고, 심볼 뒤는 공백이거나 "("이어야 함
		pattern := `(^|\s)` + regexp.QuoteMeta(symbol) + `(\s|\(|$)`

		matched, err := regexp.Match(pattern, part)
		if err == nil && matched {
			locations = append(locations, domain.NewSymbol(domain.FunctionCall, symbol, path, iter))
		}
	}

	return locations
}
