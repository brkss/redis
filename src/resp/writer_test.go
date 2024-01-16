package resp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	writeFunc func(p []byte) (n int, err error)
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return m.writeFunc(p)
}

func TestNewWriter(t *testing.T) {
	mock := &mockWriter{}
	writer := NewWriter(mock)

	assert.Equal(t, writer.writer, mock)
}

func TestWrite(t *testing.T) {
	tests := map[string]struct {
		writeFunc   func(p []byte) (n int, err error)
		expectedErr error
	}{
		"success": {
			writeFunc: func(p []byte) (n int, err error) {
				return len(p), nil
			},
			expectedErr: nil,
		},
		"error during write": {
			writeFunc: func(p []byte) (n int, err error) {
				return 0, errors.New("mocked write error")
			},
			expectedErr: errors.New("mocked write error"),
		},
	}

	for _, test := range tests {
		mock := &mockWriter{test.writeFunc}
		writer := &Writer{writer: mock}
		val := Value{}

		err := writer.Write(val)

		assert.Equal(t, err, test.expectedErr)
	}

}

func TestMarshal(t *testing.T) {

	tests := map[string]struct {
		expected string
		input    Value
		fails    bool
	}{
		"string": {
			expected: "+justasimplestring\r\n",
			input: Value{
				Typ: "string",
				Str: "justasimplestring",
			},
			fails: false,
		},
		"bulk": {
			expected: "*3\r\n$3\r\nSET\r\n$4\r\njack\r\n$5\r\nadmin\r\n",
			input: Value{
				Typ: "array",
				Arr: []Value{
					{
						Typ: "bulk",
						Blk: "SET",
					},
					{
						Typ: "bulk",
						Blk: "jack",
					},
					{
						Typ: "bulk",
						Blk: "admin",
					},
				},
			},
			fails: false,
		},
		"null": {
			expected: "$-1\r\n",
			input: Value{
				Typ: "null",
			},
			fails: false,
		},
		"error": {
			expected: "-Invalid Command\r\n",
			input: Value{
				Typ: "error",
				Str: "Invalid Command",
			},
			fails: false,
		},
		"invalid type": {
			expected: "",
			input: Value{
				Typ: "wewe",
			},
			fails: true,
		},
	}

	for _, test := range tests {
		result := test.input.Marshal()
		if test.fails {
			assert.Equal(t, []byte{}, result)
		} else {
			assert.Equal(t, test.expected, string(result))
		}
	}

}
