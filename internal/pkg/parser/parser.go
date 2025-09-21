package parser

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"sort"

	"github.com/neatflowcv/cfinder/internal/pkg/domain"
	"github.com/neatflowcv/cfinder/internal/pkg/pair"
	"github.com/neatflowcv/cfinder/internal/pkg/split"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(content []byte) []*domain.Symbol { //nolint:funlen
	delimiters := [][]byte{
		[]byte("//"),
		[]byte("/*"),
		[]byte("*/"),
		[]byte("{"),
		[]byte("}"),
		[]byte("("),
		[]byte(")"),
		[]byte(";"),
		[]byte(","),
		[]byte("#"),
		[]byte(" "),
		[]byte("\\\n"),
		[]byte("\t"),
		[]byte("\n"),
	}
	sps := split.Split(content, delimiters)

	for idx, part := range sps {
		if bytes.Equal(part, []byte("")) {
			panic(fmt.Sprintf("empty part at %d", idx))
		}
	}

	var deleted []int

	for idx, part := range sps {
		if slices.ContainsFunc([][]byte{[]byte(" "), []byte("\t")}, func(d []byte) bool {
			return slices.Equal(d, part)
		}) {
			deleted = append(deleted, idx)

			continue
		}
	}

	filtered := p.dropIndexes(deleted, sps)

	// 주석 지우기
	matches := pair.FindPair(filtered, []byte("/*"), []byte("*/"))
	filtered = p.dropMatchesWithoutEnter(filtered, matches)

	// 주석 지우기
	matches = pair.FindPair(filtered, []byte("//"), []byte("\n"))
	filtered = p.dropMatchesWithoutEnter(filtered, matches)

	// # 지우기
	matches = pair.FindPair(filtered, []byte("#"), []byte("\n"))
	filtered = p.dropMatchesWithoutEnter(filtered, matches)

	matches = pair.FindPair(filtered, []byte("("), []byte(")"))

	prev := 0
	line := 1

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Start < matches[j].Start
	})

	var ret []*domain.Symbol

	for _, match := range matches {
		for prev < match.Start {
			if bytes.Equal(filtered[prev], []byte("\n")) || bytes.Equal(filtered[prev], []byte("\\\n")) {
				line++
			}

			prev++
		}

		if !isFunctionName(filtered[match.Start-1]) {
			continue
		}

		ret = append(ret, domain.NewSymbol(string(filtered[match.Start-1]), domain.FunctionCall, line))
	}

	return ret

	// matchIdx := 0
	// line := 1

	// for idx, part := range filtered {
	// 	if bytes.Equal(part, []byte("\n")) {
	// 		line++

	// 		continue
	// 	}

	// 	if bytes.Equal(part, []byte("\\\n")) {
	// 		line++

	// 		continue
	// 	}

	// 	if matchIdx < len(matches) && idx == matches[matchIdx].Start {
	// 		matchIdx++
	// 	}

	// 	if !isFunctionName(filtered[idx-1]) {
	// 		continue
	// 	}

	// 	ret = append(ret, domain.NewSymbol(string(filtered[idx-1]), domain.FunctionCall, line))
	// }

	// return ret
}

func (p *Parser) ParseFile(reader io.Reader) ([]*domain.Symbol, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("ReadAll: %w", err)
	}

	return p.Parse(content), nil
}

func isFunctionName(part []byte) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_-]+$")

	ret := re.Match(part)
	if !ret {
		return false
	}

	if slices.ContainsFunc([][]byte{[]byte("if"), []byte("for"), []byte("while")}, func(d []byte) bool {
		return slices.Equal(d, part)
	}) {
		return false
	}

	return true
}

func (p *Parser) dropMatchesWithoutEnter(filtered [][]byte, matches []*pair.Pair) [][]byte {
	enterDelimiters := [][]byte{[]byte("\n"), []byte("//\n")}

	var deleted []int

	for _, match := range matches {
		for idx := match.Start; idx <= match.End; idx++ {
			if slices.ContainsFunc(enterDelimiters, func(e []byte) bool {
				return slices.Equal(e, filtered[idx])
			}) {
				continue
			}

			deleted = append(deleted, idx)
		}
	}

	ret := p.dropIndexes(deleted, filtered)

	return ret
}

func (*Parser) dropIndexes(deleted []int, data [][]byte) [][]byte {
	slices.Sort(deleted)
	deleted = slices.Compact(deleted)
	deletedIndex := 0

	var ret [][]byte

	for i, item := range data {
		if deletedIndex < len(deleted) && deleted[deletedIndex] == i {
			deletedIndex++

			continue
		}

		ret = append(ret, item)
	}

	return ret
}
