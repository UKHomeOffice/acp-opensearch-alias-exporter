package models

import (
	"testing"
)

func TestGetStatusChanges_same_index(t *testing.T) {
	oldStatuses := map[string]AliasStatus{"foo": {DocCount: 1, Name: "foo", Index: "foo-1"}}
	newStatuses := map[string]AliasStatus{"foo": {DocCount: 2, Name: "foo", Index: "foo-1"}}

	countChanges, _ := GetStatusChanges(oldStatuses, newStatuses)
	if len(countChanges) != 1 {
		t.Error("Expected 1 count change, received: ", len(countChanges))
	}
	if countChanges[0].DocsAdded != 1 {
		t.Error("Expected alias count change to be 100, received: ", countChanges[0].DocsAdded)
	}
	if countChanges[0].Alias != "foo" {
		t.Error("Expected first item to be alias name foo, received: ", countChanges[0].Alias)
	}
}

func TestGetStatusChanges_new(t *testing.T) {
	oldStatuses := map[string]AliasStatus{}
	newStatuses := map[string]AliasStatus{"foo": {DocCount: 2, Name: "foo", Index: "foo-1"}}

	countChanges, _ := GetStatusChanges(oldStatuses, newStatuses)
	if len(countChanges) != 1 {
		t.Error("Expected 1 count change, received: ", len(countChanges))
	}
	if countChanges[0].Alias != "foo" {
		t.Error("Expected first item to be alias name foo, received: ", countChanges[0].Alias)
	}

	if countChanges[0].DocsAdded != 2 {
		t.Error("Expected alias count change to be 100, received: ", countChanges[0].DocsAdded)
	}
}

func TestGetStatusChanges_multiple(t *testing.T) {
	oldStatuses := map[string]AliasStatus{
		"foo": {DocCount: 1, Name: "foo", Index: "foo-1"},
	}
	newStatuses := map[string]AliasStatus{
		"foo": {DocCount: 1, Name: "foo", Index: "foo-1"},
		"bar": {DocCount: 200, Name: "bar", Index: "bar-1"},
	}

	countChanges, _ := GetStatusChanges(oldStatuses, newStatuses)
	if len(countChanges) != 2 {
		t.Error("Expected 2 count change aliases, received: ", len(countChanges))
	}

	checker := func(alias string) (*StatusChange, bool) {
		for _, change := range countChanges {
			if change.Alias == alias {
				return &change, true
			}
		}
		return nil, false
	}

	if c, ok := checker("foo"); !ok {
		t.Error("Expected foo key in CountChanges")
	} else {
		if c.DocsAdded != 0 {
			t.Error("Expected alias count change to be 0, received: ", c.DocsAdded)
		}
	}

	if c, ok := checker("bar"); !ok {
		t.Error("Expected bar key in CountChanges")
	} else {
		if c.DocsAdded != 200 {
			t.Error("Expected alias count change to be 0, received: ", c.DocsAdded)
		}
	}
}
