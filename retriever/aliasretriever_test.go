package retriever

import (
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
	"testing"
)

func TestAliasgetter_GetAliasNotFailed(t *testing.T) {
	getter := func(string, string, string) ([]byte, error) {
		return []byte(`{
	"_shards": {
		"total": 1,
		"succesful": 1,
		"failed": 0
	},
	"_all": {
		"primaries": {
			"docs": {
				"count": 27178086
			}
		}
	}
}`), nil
	}

	a := NewAliasGetter("foo", "bar", "asd", getter)

	alias, err := a.GetAlias("foo-1", "foo")

	if err != nil {
		t.Error("Error while getting test alias.", err.Error())
	}

	if alias.Count != 27178086 {
		t.Error("Error not the correct count got: ", alias.Count, " expected 27178086")
	}

	if alias.Failed {
		t.Error("Error. Did not return expected Failed value. Returned: ", alias.Failed, " expected false")
	}
}

func TestAliasgetter_GetAliasFailed(t *testing.T) {
	getter := func(string, string, string) ([]byte, error) {
		return []byte(`{
	"_shards": {
		"total": 1,
		"succesful": 1,
		"failed": 1
	},
	"_all": {
		"primaries": {
			"docs": {
				"count": 27178086
			}
		}
	}
}`), nil
	}

	a := NewAliasGetter("foo", "bar", "asd", getter)

	alias, err := a.GetAlias("foo-1", "foo")

	if err != nil {
		t.Error("Error while getting test alias.", err.Error())
	}

	if alias.Count != 27178086 {
		t.Error("Error not the correct count got: ", alias.Count, " expected 27178086")
	}

	if !alias.Failed {
		t.Error("Error. Did not return expected Failed value. Returned: ", alias.Failed, " expected true")
	}
}

func TestAliasStatus_Diff_Same_Index(t *testing.T) {

	getter := func(string, string, string) ([]byte, error) {
		return []byte(`{}`), nil
	}

	ag := NewAliasGetter("foo", "bar", "asd", getter)

	newStatus := models.AliasStatus{Name: "foo", Index: "bar", Count: 2, Getter: ag}
	oldStatus := models.AliasStatus{Name: "foo", Index: "bar", Count: 1, Getter: ag}

	count, err := newStatus.Diff(oldStatus)
	if err != nil {
		t.Error("Error diffing: ", err)
	}

	if count != 1 {
		t.Error("Incorrect count diff, expected ", 1, "received: ", count)
	}
}

func TestAliasStatus_Diff_New_Index(t *testing.T) {

	getter := func(string, string, string) ([]byte, error) {
		return []byte(`{
	"_all": {
		"primaries": {
			"docs": {
				"count": 2
			}
		}
	}
}`), nil
	}

	ag := NewAliasGetter("foo", "bar", "asd", getter)

	newStatus := models.AliasStatus{Name: "foo", Index: "bar-1", Count: 1, Getter: ag}
	oldStatus := models.AliasStatus{Name: "foo", Index: "bar-2", Count: 2, Getter: ag}

	count, err := newStatus.Diff(oldStatus)
	if err != nil {
		t.Error("Error diffing: ", err)
	}

	if count != 1 {
		t.Error("Incorrect count diff, expected ", 1, "received: ", count)
	}
}

func TestAliasStatus_Diff_Different_Aliases(t *testing.T) {
	getter := func(string, string, string) ([]byte, error) {
		return []byte(`{}`), nil
	}

	ag := NewAliasGetter("foo", "bar", "asd", getter)

	a := models.AliasStatus{Name: "anything", Index: "bar", Count: 1, Getter: ag}
	b := models.AliasStatus{Name: "foo", Index: "bar", Count: 2, Getter: ag}

	_, err := a.Diff(b)

	if err == nil {
		t.Error("Expected error, but not received during diff of different alias names")
	}
}
