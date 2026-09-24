package appstore

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestNumericString_UnmarshalJSON(t *testing.T) {
	type foo struct {
		ID NumericString
	}

	tests := []struct {
		name string
		in   []byte
		err  *json.UnmarshalTypeError
		out  foo
	}{
		{
			name: "string case",
			in:   []byte("{\"ID\":\"8080\"}"),
			err:  nil,
			out:  foo{ID: "8080"},
		},
		{
			name: "number case",
			in:   []byte("{\"ID\":8080}"),
			err:  nil,
			out:  foo{ID: "8080"},
		},
		{
			name: "object case",
			in:   []byte("{\"ID\":{\"Num\": 8080}}"),
			err:  &json.UnmarshalTypeError{Value: "object", Type: reflect.TypeFor[json.Number]()},
			out:  foo{},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := foo{}
			err := json.Unmarshal(v.in, &out)

			if v.err != nil {
				// Compare only Value and Type: Go 1.27 no longer adds the enclosing
				// struct field to errors returned from UnmarshalJSON.
				var typeErr *json.UnmarshalTypeError
				if !errors.As(err, &typeErr) || typeErr.Value != v.err.Value || typeErr.Type != v.err.Type {
					t.Errorf("input: %s, get: %v, want: %v\n", v.in, err, v.err)
				}
				return
			}
			if err != nil {
				t.Errorf("input: %s, get: %s, want: no error\n", v.in, err)
				return
			}

			if !reflect.DeepEqual(out, v.out) {
				t.Errorf("input: %s, get: %v, want: %v\n", v.in, out, v.out)
			}
		})
	}
}
