package handler

import (
	"fmt"
	"testing"

	"github.com/brkss/redis/src/resp"
	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {

	tests := map[string]struct {
		input       []resp.Value
		expected    resp.Value
		expectedErr resp.Value
		fails       bool
		f           string
	}{
		"SET_SUCCESS": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "jack",
				},
				{
					Typ: "bulk",
					Blk: "admin",
				},
			},
			expected: resp.Value{Typ: "string", Str: "OK"},
			fails:    false,
			f:        "SET",
		},
		"SET_INVALID_ARGUMENT": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "jack",
				},
			},
			expectedErr: resp.Value{Typ: "error", Str: "ERR wrong number of argements for 'set' command"},
			fails:       true,
			f:           "SET",
		},
		"GET_SUCCESS": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "jack",
				},
			},
			expected: resp.Value{Typ: "bulk", Blk: "admin"},
			fails:    false,
			f:        "GET",
		},
		"GET_INVALID_ARGEMENTS": {
			input:       []resp.Value{},
			expectedErr: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"},

			fails: true,
			f:     "GET",
		},
		"GET_NOT_FOUND": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "hack",
				},
			},
			expectedErr: resp.Value{Typ: "null"},

			fails: true,
			f:     "GET",
		},
	}

	for title, test := range tests {
		fmt.Println("test >> ", title)
		if test.f == "SET" { // SET
			res := set(test.input)
			if test.fails == true {
				assert.Equal(t, test.expectedErr, res)
			} else {
				// check result
				assert.Equal(t, test.expected, res)
				// // test get
				// val := get([]resp.Value{test.input[0]})
				// assert.Equal(t, val, test.input[1])
			}
		} else { // GET
			val := get(test.input)
			if test.fails == true {
				assert.Equal(t, val, test.expectedErr)
			} else {
				assert.Equal(t, val, test.expected)
			}
		}

	}

}
