package resp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected Value
		fails    bool
	}{
		"correct parsing": {
			input: "*3\r\n$3\r\nSET\r\n$4\r\njack\r\n$5\r\nadmin\r\n",
			expected: Value{
				typ: "array",
				arr: []Value{
					{
						typ: "bulk",
						blk: "SET",
					},
					{
						typ: "bulk",
						blk: "jack",
					},
					{
						typ: "bulk",
						blk: "admin",
					},
				},
			},
			fails: false,
		},
		"incorrect type": {
			input:    "/3\r\n$3\r\nSET\r\n$4\r\njack\r\n$5\r\nadmin\r\n",
			expected: Value{},
			fails:    true,
		},
	}
	for _, test := range tests {
		resp := NewReader(strings.NewReader(test.input))
		val, err := resp.Read()
		if test.fails {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.expected, val)
		}

	}

}
