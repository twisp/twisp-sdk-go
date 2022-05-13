package client

import (
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

type mergeTestCase struct {
	vars  map[string]any
	strct *TestStruct

	expected map[string]any
	message  string
}

type TestStruct struct {
	Field   string           `json:"field"`
	Nested  TestStructNested `json:"nested"`
	Collide string           `json:"collide"`
}

type TestStructNested struct {
	FieldNested   string `json:"fieldNested"`
	NestedCollide string `json:"nestedCollide"`
}

func TestMerge(t *testing.T) {
	testCases := []mergeTestCase{
		{
			vars:  map[string]any{"field": "field"},
			strct: nil,

			expected: map[string]any{"field": "field"},
			message:  "doesnt merge when strct is nil",
		},
		{
			vars: map[string]any{},
			strct: &TestStruct{
				Field: "field",
				Nested: TestStructNested{
					FieldNested:   "fieldNested",
					NestedCollide: "nestedCollide",
				},
				Collide: "collide",
			},

			expected: map[string]any{"field": "field", "collide": "collide", "nested": map[string]any{"nestedCollide": "nestedCollide", "fieldNested": "fieldNested"}},
			message:  "doesnt merge when map is empty",
		},
		{
			vars: map[string]any{"field": "this", "collide": "should", "nested": map[string]any{"nestedCollide": "not", "fieldNested": "show"}},
			strct: &TestStruct{
				Field: "field",
				Nested: TestStructNested{
					FieldNested:   "fieldNested",
					NestedCollide: "nestedCollide",
				},
				Collide: "collide",
			},

			expected: map[string]any{"field": "field", "collide": "collide", "nested": map[string]any{"nestedCollide": "nestedCollide", "fieldNested": "fieldNested"}},
			message:  "favors struct",
		},
		{
			vars: map[string]any{"field1": "field1", "nested": map[string]any{"fieldNested1": "fieldNested1"}, "anotherNested": map[string]any{"field": "field"}},
			strct: &TestStruct{
				Field: "field",
				Nested: TestStructNested{
					FieldNested:   "fieldNested",
					NestedCollide: "nestedCollide",
				},
				Collide: "collide",
			},

			expected: map[string]any{"field": "field", "field1": "field1", "collide": "collide", "nested": map[string]any{"nestedCollide": "nestedCollide", "fieldNested": "fieldNested", "fieldNested1": "fieldNested1"}, "anotherNested": map[string]any{"field": "field"}},
			message:  "when no collisions produces superset",
		},
	}

	for i, testCase := range testCases {
		req := &graphql.Request{Variables: testCase.strct}
		err := merge(testCase.vars, req)
		require.NoError(t, err, "test case %d %s", i, testCase.message)
		require.EqualValues(t, testCase.expected, req.Variables, "test case %d %s", i, testCase.message)
	}
}
