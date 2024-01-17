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

func TestHashSet(t *testing.T) {
	tests := map[string]struct {
		input       []resp.Value
		expected    resp.Value
		fails       bool
		expectedErr resp.Value
	}{
		"HSET_SUCCESS": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "users",
				},
				{
					Typ: "bulk",
					Blk: "brahim",
				},
				{
					Typ: "bulk",
					Blk: "admin",
				},
			},
			expected: resp.Value{Typ: "string", Str: "OK"},
			fails:    false,
		},
		"HSET_INVALID_ARGS": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "users",
				},
				{
					Typ: "bulk",
					Blk: "brahim",
				},
			},
			expectedErr: resp.Value{Typ: "error", Str: "ERR invalid arguments for 'hset' command"},
			fails:       true,
		},
	}

	for _, test := range tests {
		results := hset(test.input)
		if test.fails == true {
			assert.Equal(t, test.expectedErr, results)
		} else {
			assert.Equal(t, test.expected, results)
		}
	}
}

func TestHashGet(t *testing.T) {

	tests := map[string]struct {
		input       []resp.Value
		expected    resp.Value
		fails       bool
		expectedErr resp.Value
	}{

		"HGET_SUCCESS": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "users",
				},
				{
					Typ: "bulk",
					Blk: "brahim",
				},
			},
			expected: resp.Value{Typ: "bulk", Blk: "admin"},
			fails:    false,
		},
		"HGET_INVALID_ARGUMENT": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "users",
				},
			},
			expectedErr: resp.Value{Typ: "error", Str: "ERR invalid arguments for 'hget' command"},
			fails:       true,
		},
		"HGET_NOT_FOUND": {
			input: []resp.Value{
				{
					Typ: "bulk",
					Blk: "usersx",
				},
				{
					Typ: "bulk",
					Blk: "do_not_exist",
				},
			},
			expected: resp.Value{Typ: "null"},
			fails:    false,
		},
	}

	args := []resp.Value{
		{
			Typ: "bulk",
			Blk: "users",
		},
		{
			Typ: "bulk",
			Blk: "brahim",
		},
		{
			Typ: "bulk",
			Blk: "admin",
		},
	}

	_ = hset(args)

	for title, test := range tests {
		fmt.Println(">>>", title)
		results := hget(test.input)
		if test.fails == true {
			assert.Equal(t, test.expectedErr, results)
		} else {
			assert.Equal(t, test.expected, results)
		}
	}

}
