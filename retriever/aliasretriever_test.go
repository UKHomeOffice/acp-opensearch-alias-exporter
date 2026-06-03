package retriever

import (
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
	"testing"
)

func TestAliasgetter_GetAliasNotFailed(t *testing.T) {
	callCount := 0
	getter := func(string, string, string) ([]byte, error) {
		callCount++
		if callCount == 1 {
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
			},
			"indexing": {
				"index_failed": 0
			}	
		}
	}
}`), nil
		}
		return []byte(`{
	"action" : {
      "name" : "rollover",
      "start_time" : 1667411127803,
      "index" : 0,
      "failed" : false,
      "consumed_retries" : 0,
      "last_retry_time" : 0
    }
}`), nil
	}

	a := NewAliasGetter("foo", "bar", "asd", getter)

	alias, err := a.GetAlias("foo-1", "foo")

	if err != nil {
		t.Error("Error while getting test alias.", err.Error())
	}

	if alias.Count != 27178086 {
		t.Error("Error. Did not return expected Count value. Returned:", alias.Count, " expected 27178086")
	}

	if alias.FailedShards != 0 {
		t.Error("Error. Did not return expected FailedShards value. Returned: ", alias.FailedShards, " expected 0")
	}

	if alias.FailedIndexOperations != 0{
		t.Error("Error. Did not return expected FailedIndexOperations value. Returned: ", alias.FailedIndexOperations, " expected 0")
	}

	if alias.RolloverAttemptFailed {
		t.Error("Error. Did not return expected RolloverAttemptFailed value. Returned: ", alias.RolloverAttemptFailed, " expected false")
	}
}

func TestAliasgetter_GetAliasFailed(t *testing.T) {
	callCount := 0
	getter := func(string, string, string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return []byte(`{
	"_shards": {
		"total": 2,
		"succesful": 1,
		"failed": 1
	},
	"_all": {
		"primaries": {
			"docs": {
				"count": 27178086
			},
			"indexing": {
				"index_failed": 1
			}	
		}
	}
}`), nil
		}
		return []byte(`{
	"action" : {
      "name" : "rollover",
      "start_time" : 1667411127803,
      "index" : 0,
      "failed" : true,
      "consumed_retries" : 0,
      "last_retry_time" : 0
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

	if alias.FailedShards != 1 {
		t.Error("Error. Did not return expected FailedShards value. Returned: ", alias.FailedShards, " expected 1")
	}

	if alias.FailedIndexOperations != 1 {
		t.Error("Error. Did not return expected FailedIndexOperations value. Returned: ", alias.FailedIndexOperations, " expected 1")
	}

	if !alias.RolloverAttemptFailed {
		t.Error("Error. Did not return expected RolloverAttemptFailed value. Returned: ", alias.RolloverAttemptFailed, " expected true")
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
