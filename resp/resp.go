package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	STRING  = "+"
	ERROR   = "-"
	INTEGER = ":"
	BULK    = "$"
	ARRAY   = "*"
)

type Value struct {
	typ string
	str string
	num int
	blk string
	arr []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Resp {
	return &Resp{bufio.NewReader(rd)}
}

// readLine: reads the line from the buffer
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readInteger: read integer from reader !
func (r *Resp) readInteger() (x, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(i64), n, nil
}

func ParseInput(input string) (*Value, error) {

	reader := bufio.NewReader(strings.NewReader(input))
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if b != '$' {
		err := fmt.Errorf("Invalid type, expecting bulk strings ($) got : %c", b)
		return nil, err
	}

	size, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		return nil, err
	}

	// consume /r/n
	reader.ReadByte()
	reader.ReadByte()

}
