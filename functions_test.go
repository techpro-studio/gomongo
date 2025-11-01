package gomongo

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"reflect"
	"strings"
	"testing"
)

func TestBuildUpdateDoc(t *testing.T) {
	type sample struct {
		Name   string `bson:"name"`
		Email  string `bson:"email,omitempty"`
		Age    int    `bson:"age"`
		Active bool   `bson:"active"`
		Note   string // no bson tag
	}

	tests := []struct {
		name     string
		input    sample
		expected bson.M
	}{
		{
			name: "all fields non-zero",
			input: sample{
				Name:   "Alice",
				Email:  "a@example.com",
				Age:    30,
				Active: true,
				Note:   "ignored",
			},
			expected: bson.M{
				"name":   "Alice",
				"email":  "a@example.com",
				"age":    30,
				"active": true,
			},
		},
		{
			name: "zero string and int omitted",
			input: sample{
				Name:   "",
				Email:  "",
				Age:    0,
				Active: false,
			},
			expected: bson.M{},
		},
		{
			name: "some fields zero, others set",
			input: sample{
				Name:   "Bob",
				Email:  "",
				Age:    25,
				Active: false,
			},
			expected: bson.M{
				"name": "Bob",
				"age":  25,
			},
		},
		{
			name: "bson tag with omitempty trimmed correctly",
			input: sample{
				Email: "x@example.com",
			},
			expected: bson.M{
				"email": "x@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildUpdateDoc(tt.input)

			// Check same number of keys
			if len(got) != len(tt.expected) {
				t.Errorf("[%s] expected %d keys, got %d (%v)", tt.name, len(tt.expected), len(got), got)
			}

			// Check all expected keys and values
			for k, v := range tt.expected {
				gotVal, ok := got[k]
				if !ok {
					t.Errorf("[%s] missing key %s", tt.name, k)
					continue
				}
				if !reflect.DeepEqual(gotVal, v) {
					t.Errorf("[%s] for key %s expected %v, got %v", tt.name, k, v, gotVal)
				}
			}

			// Ensure no field without bson tag
			for k := range got {
				if strings.Contains(k, "Note") {
					t.Errorf("[%s] found unexpected Note field in update: %v", tt.name, got)
				}
			}
		})
	}
}
