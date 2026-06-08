package retriever

import (
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
	"testing"
)

func TestAliasgetter_GetAliasWithFailures(t *testing.T) {
	callCount := 0
	getter := func(string, string, string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return []byte(`{
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
			"foo-1": {
				"action" : {
					"name" : "rollover",
					"start_time" : 1667411127803,
					"index" : 0,
					"failed" : false,
					"consumed_retries" : 0,
					"last_retry_time" : 0
				}
			}
		}`), nil
	}

	a := NewAliasGetter("foo", "bar", "asd", getter)

	alias, err := a.GetAlias("foo-1", "foo")

	if err != nil {
		t.Error("Error while getting test alias.", err.Error())
	}

	if alias.DocCount != 27178086 {
		t.Error("Error. Did not return expected DocCount value. Returned:", alias.DocCount, " expected 27178086")
	}

	if alias.FailedIndexOperations != 0{
		t.Error("Error. Did not return expected FailedIndexOperations value. Returned: ", alias.FailedIndexOperations, " expected 0")
	}

	if alias.RolloverAttemptFailed {
		t.Error("Error. Did not return expected RolloverAttemptFailed value. Returned: ", alias.RolloverAttemptFailed, " expected false")
	}
}

func TestAliasgetter_GetAliasWithoutFailures(t *testing.T) {
	callCount := 0
	getter := func(string, string, string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return []byte(`{
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
			"foo-1": {
				"action" : {
					"name" : "rollover",
					"start_time" : 1667411127803,
					"index" : 0,
					"failed" : true,
					"consumed_retries" : 0,
					"last_retry_time" : 0
				}
			}
		}`), nil
	}

	a := NewAliasGetter("foo", "bar", "asd", getter)

	alias, err := a.GetAlias("foo-1", "foo")

	if err != nil {
		t.Error("Error while getting test alias.", err.Error())
	}

	if alias.DocCount != 27178086 {
		t.Error("Error not the correct DocCount got: ", alias.DocCount, " expected 27178086")
	}

	if alias.FailedIndexOperations != 1 {
		t.Error("Error. Did not return expected FailedIndexOperations value. Returned: ", alias.FailedIndexOperations, " expected 1")
	}

	if !alias.RolloverAttemptFailed {
		t.Error("Error. Did not return expected RolloverAttemptFailed value. Returned: ", alias.RolloverAttemptFailed, " expected true")
	}
}

func TestAliasStatus_GetDifference_Same_Index(t *testing.T) {

	getter := func(string, string, string) ([]byte, error) {
		return []byte(`{}`), nil
	}

	ag := NewAliasGetter("foo", "bar", "asd", getter)

	newStatus := models.AliasStatus{Name: "foo", Index: "bar", DocCount: 2, Getter: ag, FailedIndexOperations: 3}
	oldStatus := models.AliasStatus{Name: "foo", Index: "bar", DocCount: 1, Getter: ag, FailedIndexOperations: 2}

	difference, err := newStatus.GetDifference(oldStatus)
	if err != nil {
		t.Error("Error during GetDifference: ", err)
	}

	if difference.Docs != 1 {
		t.Error("Incorrect Docs difference, expected ", 1, "received: ", difference.Docs)
	}

	if difference.FailedIndexOperations != 1 {
		t.Error("Incorrect FailedIndexOperations Difference, expected ", 1, "received: ", difference.FailedIndexOperations)
	}
}

func TestAliasStatus_GetDifference_New_Index(t *testing.T) {

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

	newStatus := models.AliasStatus{Name: "foo", Index: "bar-1", DocCount: 1, Getter: ag}
	oldStatus := models.AliasStatus{Name: "foo", Index: "bar-2", DocCount: 2, Getter: ag}

	difference, err := newStatus.GetDifference(oldStatus)
	if err != nil {
		t.Error("Error GetDifferenceing: ", err)
	}

	if difference.Docs != 1 {
		t.Error("Incorrect count GetDifference, expected ", 1, "received: ", difference.Docs)
	}
}

func TestAliasStatus_GetDifference_Different_Aliases(t *testing.T) {
	getter := func(string, string, string) ([]byte, error) {
		return []byte(`{}`), nil
	}

	ag := NewAliasGetter("foo", "bar", "asd", getter)

	a := models.AliasStatus{Name: "anything", Index: "bar", DocCount: 1, Getter: ag}
	b := models.AliasStatus{Name: "foo", Index: "bar", DocCount: 2, Getter: ag}

	_, err := a.GetDifference(b)

	if err == nil {
		t.Error("Expected error, but not received during GetDifference of GetDifferenceerent alias names")
	}
}
