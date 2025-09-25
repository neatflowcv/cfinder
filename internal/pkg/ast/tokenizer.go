package ast

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
)

type TokenKind int

const (
	TokenKindOther TokenKind = iota + 1
	TokenKindNewLine
	TokenKindLineContinuation
	TokenKindOpenBrace
	TokenKindCloseBrace
	TokenKindOpenParenthesis
	TokenKindCloseParenthesis
	TokenKindSemicolon
	TokenKindComma
	TokenKindPreprocessor
	TokenKindSpace
	TokenKindSingleLineComment
	TokenKindOpenMultiLineComment
	TokenKindCloseMultiLineComment
	TokenKindAsterisk
	TokenKindDoubleQuote
)

type Token struct {
	Kind  TokenKind
	Value string
}

func Tokenize(reader io.Reader) []*Token {
	bufReader := bufio.NewReader(reader)

	var ret []*Token

	var buf bytes.Buffer

	for {
		char, err := bufReader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatal("readChar", err)
		}

		kind := detectKind(char, bufReader)
		if kind == TokenKindOther {
			buf.WriteByte(char)

			continue
		}

		if buf.Len() > 0 {
			ret = append(ret, &Token{Kind: TokenKindOther, Value: buf.String()})
			buf.Reset()
		}

		if kind == TokenKindSpace {
			continue
		}

		ret = append(ret, &Token{Kind: kind, Value: ""})
	}

	if buf.Len() > 0 {
		ret = append(ret, &Token{Kind: TokenKindOther, Value: buf.String()})
	}

	return ret
}

func detectKind(char byte, bufReader *bufio.Reader) TokenKind { //nolint:cyclop,funlen
	switch char {
	case '\n':
		return TokenKindNewLine

	case '\t':
		return TokenKindSpace

	case ' ':
		return TokenKindSpace

	case '{':
		return TokenKindOpenBrace

	case '}':
		return TokenKindCloseBrace

	case '(':
		return TokenKindOpenParenthesis

	case ')':
		return TokenKindCloseParenthesis

	case ';':
		return TokenKindSemicolon

	case '#':
		return TokenKindPreprocessor

	case '\\':
		return TokenKindLineContinuation

	case '/':
		more, err := bufReader.Peek(1)
		if err != nil {
			return TokenKindOther
		}

		if more[0] == '/' {
			_, _ = bufReader.Discard(1)

			return TokenKindSingleLineComment
		}

		if more[0] == '*' {
			_, _ = bufReader.Discard(1)

			return TokenKindOpenMultiLineComment
		}

		return TokenKindOther

	case '*':
		more, err := bufReader.Peek(1)
		if err != nil {
			return TokenKindAsterisk
		}

		if more[0] == '/' {
			_, _ = bufReader.Discard(1)

			return TokenKindCloseMultiLineComment
		}

		return TokenKindAsterisk

	case ',':
		return TokenKindComma

	case '"':
		return TokenKindDoubleQuote

	default:
		return TokenKindOther
	}
}
