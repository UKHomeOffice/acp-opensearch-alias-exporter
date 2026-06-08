package models

import (
	"testing"
)

func TestGetStatusChanges_same_index(t *testing.T) {
	oldStatuses := map[string]AliasStatus{"foo": {DocCount: 1, Name: "foo", Index: "foo-1", FailedIndexOperations: 2}}
	newStatuses := map[string]AliasStatus{"foo": {DocCount: 2, Name: "foo", Index: "foo-1", FailedIndexOperations: 3}}

	statusChanges, _ := GetStatusChanges(oldStatuses, newStatuses)
	if len(statusChanges) != 1 {
		t.Error("Expected 1 count change, received: ", len(statusChanges))
	}
	if statusChanges[0].DocsAdded != 1 {
		t.Error("Expected DocsAdded to be 1, received: ", statusChanges[0].DocsAdded)
	}
	if statusChanges[0].Alias != "foo" {
		t.Error("Expected first item to be alias named foo, received: ", statusChanges[0].Alias)
	}
	if statusChanges[0].NewIndexOperationFailures != 1 {
		t.Error("Expected NewIndexOperationFailures to be 1, received: ", statusChanges[0].DocsAdded)
	}
}

type mockAliasGetter struct {
	statusToReturn AliasStatus
}

func (m *mockAliasGetter) GetAlias(string, string) (AliasStatus, error) {
	return m.statusToReturn, nil
}

func TestGetStatusChanges_new_index(t *testing.T) {
	mock := &mockAliasGetter{
		statusToReturn: AliasStatus{
			DocCount:              2,
			FailedIndexOperations: 3,
		},
	}
	
	oldStatuses := map[string]AliasStatus{"foo": {DocCount: 1, Name: "foo", Index: "foo-1", FailedIndexOperations: 2, Getter: mock}}
	newStatuses := map[string]AliasStatus{"foo": {DocCount: 2, Name: "foo", Index: "foo-2", FailedIndexOperations: 3}}

	statusChanges, _ := GetStatusChanges(oldStatuses, newStatuses)

	if len(statusChanges) != 1 {
		t.Error("Expected 1 count change, received: ", len(statusChanges))
	}

	if statusChanges[0].DocsAdded != 3 {
		t.Error("Expected DocsAdded to be 1, received: ", statusChanges[0].DocsAdded)
	}
	if statusChanges[0].Alias != "foo" {
		t.Error("Expected first item to be alias named foo, received: ", statusChanges[0].Alias)
	}
	if statusChanges[0].NewIndexOperationFailures != 4 {
		t.Error("Expected NewIndexOperationFailures to be 1, received: ", statusChanges[0].DocsAdded)
	}
}

func TestGetStatusChanges_new_alias(t *testing.T) {
	oldStatuses := map[string]AliasStatus{}
	newStatuses := map[string]AliasStatus{"foo": {DocCount: 2, Name: "foo", Index: "foo-1"}}

	statusChanges, _ := GetStatusChanges(oldStatuses, newStatuses)
	if len(statusChanges) != 1 {
		t.Error("Expected 1 status change, received: ", len(statusChanges))
	}
	if statusChanges[0].Alias != "foo" {
		t.Error("Expected first item to have alias name foo, received: ", statusChanges[0].Alias)
	}

	if statusChanges[0].DocsAdded != 2 {
		t.Error("Expected alias DocsAdded to be 2, received: ", statusChanges[0].DocsAdded)
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

	statusChanges, _ := GetStatusChanges(oldStatuses, newStatuses)
	if len(statusChanges) != 2 {
		t.Error("Expected 2 count change aliases, received: ", len(statusChanges))
	}

	checker := func(alias string) (*StatusChange, bool) {
		for _, change := range statusChanges {
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
