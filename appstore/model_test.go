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
		err  error
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
			err:  errors.New("json: cannot unmarshal object into Go struct field foo.ID of type json.Number"),
			out:  foo{},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := foo{}
			err := json.Unmarshal(v.in, &out)

			if err != nil {
				if err.Error() != v.err.Error() {
					t.Errorf("input: %s, get: %s, want: %s\n", v.in, err, v.err)
				}
				return
			}

			if !reflect.DeepEqual(out, v.out) {
				t.Errorf("input: %s, get: %v, want: %v\n", v.in, out, v.out)
			}
		})
	}
}
