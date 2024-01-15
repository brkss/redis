package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
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

// Read: read parses Value from bufio reader
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, nil
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		err := fmt.Errorf("Unkown type %v\n", string(_type))
		return Value{}, err
	}
}

// readArray: reads array from bufio reader
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of the array
	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	v.arr = make([]Value, 0)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		// append parsed value to array
		v.arr = append(v.arr, val)
	}

	return v, nil
}

// readBulk: read and parse a bulk string from the bufio reader
func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	str := make([]byte, length)
	v.blk = string(str)

	// read the trailing CRLF
	r.readLine()

	return v, nil
}
