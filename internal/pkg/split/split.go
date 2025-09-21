package split

import (
	"bufio"
	"bytes"
	"io"
)

func Split(content []byte, delimiters [][]byte) [][]byte {
	reader := bufio.NewReader(bytes.NewReader(content))

	return SplitWithReader(reader, delimiters)
}

func SplitWithReader(reader *bufio.Reader, delimiters [][]byte) [][]byte { //nolint:cyclop
	var (
		ret [][]byte
		buf bytes.Buffer
	)

	for {
		_, err := reader.Peek(1)
		if err == io.EOF {
			break
		}

		isFound := false

		for _, delimiter := range delimiters {
			cont, err := reader.Peek(len(delimiter))
			if err != nil {
				continue
			}

			if bytes.Equal(delimiter, cont) {
				copied := append([]byte{}, buf.Bytes()...)
				if len(copied) > 0 {
					ret = append(ret, copied)
				}

				ret = append(ret, delimiter)

				buf.Reset()

				_, _ = reader.Discard(len(delimiter))

				isFound = true

				break
			}
		}

		if !isFound {
			ch, err := reader.ReadByte()
			if err != nil {
				break
			}

			buf.WriteByte(ch)
		}
	}

	copied := append([]byte{}, buf.Bytes()...)
	if len(copied) > 0 {
		ret = append(ret, copied)
	}

	if bytes.Equal(ret[0], []byte("")) {
		ret = ret[1:]
	}

	return ret
}
