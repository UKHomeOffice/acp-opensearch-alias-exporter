package retriever

import (
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
	"slices"
	"testing"
)

func TestAliasgetter_GetAliasWithFailures(t *testing.T) {
	callCount := 0
	getter := func(url string, name string, password string) ([]byte, error) {
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
		if callCount == 2 {
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
		return []byte(`{
			"cluster_name" : "658718335966:acp-prod-logging-v2",
			"status" : "green",
			"timed_out" : false,
			"number_of_nodes" : 11,
			"number_of_data_nodes" : 8,
			"discovered_master" : true,
			"discovered_cluster_manager" : true,
			"active_primary_shards" : 2,
			"active_shards" : 4,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 0,
			"delayed_unassigned_shards" : 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch" : 0,
			"task_max_waiting_in_queue_millis" : 0,
			"active_shards_percent_as_number" : 100.0
		}`), nil
	}

	a := NewAliasGetter("foo", "bar", "asd", getter)

	alias, err := a.GetAlias("foo-1", "foo")

	if err != nil {
		t.Error("Error while getting test alias.", err.Error())
	}

	if callCount != 3 {
		t.Error("Error. Did not make expected nuo. of API calls. Returned:", callCount, "; Expected: 3")
	}
	
	if alias.DocCount != 27178086 {
		t.Error("Error. Did not return expected DocCount value. Returned:", alias.DocCount, " expected 27178086")
	}

	if alias.FailedIndexOperations != 0 {
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
							"index_failed": 0
						}	
					}
				}
			}`), nil
		}
		if callCount == 2 {
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
		return []byte(`{
			"cluster_name" : "658718335966:acp-prod-logging-v2",
			"status" : "green",
			"timed_out" : false,
			"number_of_nodes" : 11,
			"number_of_data_nodes" : 8,
			"discovered_master" : true,
			"discovered_cluster_manager" : true,
			"active_primary_shards" : 2,
			"active_shards" : 4,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 0,
			"delayed_unassigned_shards" : 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch" : 0,
			"task_max_waiting_in_queue_millis" : 0,
			"active_shards_percent_as_number" : 100.0
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

	if alias.FailedIndexOperations != 0 {
		t.Error("Error. Did not return expected FailedIndexOperations value. Returned: ", alias.FailedIndexOperations, " expected 0")
	}

	if alias.RolloverAttemptFailed {
		t.Error("Error. Did not return expected RolloverAttemptFailed value. Returned: ", alias.RolloverAttemptFailed, " expected false")
	}
}

func TestAliasgetter_GetAlias_urls_called_and_health(t *testing.T) {
	expected_urls := []string{"foo/foo-1/_stats", "foo/_plugins/_ism/explain/foo-1", "foo/_cluster/health/foo-1"}
	var urls_called []string
	callCount := 0
	getter := func(url string, name string, password string) ([]byte, error) {
		urls_called = append(urls_called, url)
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
		if callCount == 2 {
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
		return []byte(`{
			"cluster_name" : "658718335966:acp-prod-logging-v2",
			"status" : "green",
			"timed_out" : false,
			"number_of_nodes" : 11,
			"number_of_data_nodes" : 8,
			"discovered_master" : true,
			"discovered_cluster_manager" : true,
			"active_primary_shards" : 2,
			"active_shards" : 4,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 0,
			"delayed_unassigned_shards" : 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch" : 0,
			"task_max_waiting_in_queue_millis" : 0,
			"active_shards_percent_as_number" : 100.0
		}`), nil
	}

	a := NewAliasGetter("foo", "bar", "asd", getter)

	alias, err := a.GetAlias("foo-1", "foo")

	if err != nil {
		t.Error("Error while getting test alias.", err.Error())
	}

	if alias.Health != "green" {
		t.Error("Error. Did not return expected Health value. Returned: ", alias.Health, "; Expected: green")
	}

	if !slices.Equal(urls_called, expected_urls) {
		t.Error("Error. Unexpected urls called by getter. Returned: ", urls_called, "; Expected", expected_urls)
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
	callCount := 0
	getter := func(string, string, string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return []byte(`{
				"_all": {
					"primaries": {
						"docs": {
							"count": 2
						},
						"indexing": {
							"index_failed": 0
						}	
					}
				}
			}`), nil
		}
		if callCount == 2 {
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
		return []byte(`{
			"cluster_name" : "658718335966:acp-prod-logging-v2",
			"status" : "green",
			"timed_out" : false,
			"number_of_nodes" : 11,
			"number_of_data_nodes" : 8,
			"discovered_master" : true,
			"discovered_cluster_manager" : true,
			"active_primary_shards" : 2,
			"active_shards" : 4,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 0,
			"delayed_unassigned_shards" : 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch" : 0,
			"task_max_waiting_in_queue_millis" : 0,
			"active_shards_percent_as_number" : 100.0
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
