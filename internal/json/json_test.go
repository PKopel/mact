package internal

import (
	"reflect"
	"testing"

	conf "github.com/PKopel/mact/internal/config"
)

func TestFindField(t *testing.T) {
	cases := map[string]struct {
		object   JSON
		path     []string
		result   JSON
		leftover []string
	}{
		"shallow, happy path": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			path: []string{"field1"},
			result: JSON{
				"inner": "value",
			},
			leftover: nil,
		},
		"deep, happy path": {
			object: JSON{
				"field1": &JSON{
					"inner1": &JSON{
						"inner2": "value",
					},
				},
				"field2": "value2",
			},
			path: []string{"field1", "inner1"},
			result: JSON{
				"inner2": "value",
			},
			leftover: nil,
		},
		"shallow, too long path": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			path: []string{"field1", "inner"},
			result: JSON{
				"inner": "value",
			},
			leftover: []string{"inner"},
		},
		"deep, too long path": {
			object: JSON{
				"field1": &JSON{
					"inner1": &JSON{
						"inner2": "value",
					},
				},
				"field2": "value2",
			},
			path: []string{"field1", "inner1", "inner2"},
			result: JSON{
				"inner2": "value",
			},
			leftover: []string{"inner2"},
		},
		"shallow, no match": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			path: []string{"field3"},
			result: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			leftover: []string{"field3"},
		},
		"deep, no match": {
			object: JSON{
				"field1": &JSON{
					"inner1": &JSON{
						"inner2": "value",
					},
				},
				"field2": "value2",
			},
			path: []string{"field1", "inner3"},
			result: JSON{
				"inner1": &JSON{
					"inner2": "value",
				},
			},
			leftover: []string{"inner3"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result, leftover := FindField(tc.object, tc.path)
			if !reflect.DeepEqual(result, tc.result) || !reflect.DeepEqual(leftover, tc.leftover) {
				t.Errorf("Expected %v %v, got %v %v", tc.result, tc.leftover, result, leftover)
			}
		})
	}
}

func TestApplyChanges(t *testing.T) {
	cases := map[string]struct {
		object  JSON
		changes []conf.Change
		result  JSON
	}{
		"shallow, one add": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Add,
					Field: "field3",
					Value: "new value",
				},
			},
			result: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
				"field3": "new value",
			},
		},
		"shallow, one modify": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Modify,
					Field: "field1",
					Value: "changed value",
				},
			},
			result: JSON{
				"field1": "changed value",
				"field2": "value2",
			},
		},
		"shallow, one remove": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Remove,
					Field: "field1",
				},
			},
			result: JSON{
				"field2": "value2",
			},
		},
		"shallow, multiple changes": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Add,
					Field: "field3",
					Value: "new value",
				},
				{
					Type:  conf.Modify,
					Field: "field1",
					Value: "changed value",
				},
				{
					Type:  conf.Remove,
					Field: "field2",
				},
			},
			result: JSON{
				"field1": "changed value",
				"field3": "new value",
			},
		},
		"deep, one add": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Add,
					Field: "field1.inner2",
					Value: "new value",
				},
			},
			result: JSON{
				"field1": &JSON{
					"inner":  "value",
					"inner2": "new value",
				},
				"field2": "value2",
			},
		},
		"deep, one modify": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Modify,
					Field: "field1.inner",
					Value: "changed value",
				},
			},
			result: JSON{
				"field1": &JSON{
					"inner": "changed value",
				},
				"field2": "value2",
			},
		},
		"deep, one remove": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Remove,
					Field: "field1.inner",
				},
			},
			result: JSON{
				"field1": &JSON{},
				"field2": "value2",
			},
		},
		"deep, multiple changes": {
			object: JSON{
				"field1": &JSON{
					"inner": "value",
				},
				"field2": "value2",
			},
			changes: []conf.Change{
				{
					Type:  conf.Add,
					Field: "field1.inner2",
					Value: "new value",
				},
				{
					Type:  conf.Modify,
					Field: "field1.inner",
					Value: &JSON{
						"inner3": "value3",
						"inner4": "value4",
					},
				},
				{
					Type:  conf.Remove,
					Field: "field1.inner.inner3",
				},
			},
			result: JSON{
				"field1": &JSON{
					"inner": &JSON{
						"inner4": "value4",
					},
					"inner2": "new value",
				},
				"field2": "value2",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := ApplyChanges(tc.object, tc.changes)
			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Expected %v, got %v", tc.result, result)
			}
		})
	}

}
