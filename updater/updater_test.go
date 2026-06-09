package updater

import (
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"strings"
	"testing"
)

func testProvideUpdater(metricName string) *PrometheusUpdater {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Namespace: "opensearch", Name: metricName, Help: "test"},
		[]string{"namespace"},
	)
	return &PrometheusUpdater{gaugeVec: gaugeVec}
}

var TestASs = models.AliasStatuses{
	"healthy_foo": models.AliasStatus{
		Name:                  "foo",
		Health:                "green",
		RolloverAttemptFailed: false,
	},
	"poorly_bar": models.AliasStatus{
		Name:                  "bar",
		Health:                "yellow",
		RolloverAttemptFailed: false,
	},
	"very_poorly_boo": models.AliasStatus{
		Name:                  "boo",
		Health:                "red",
		RolloverAttemptFailed: true,
	},
}

func TestUpdateHealth(t *testing.T) {
	testUpdater := testProvideUpdater("index_health")
	testUpdater.UpdateHealth(TestASs)
	expected := `
        # HELP opensearch_index_health test
        # TYPE opensearch_index_health gauge
        opensearch_index_health{namespace="bar"} 0.5
		opensearch_index_health{namespace="boo"} 1
		opensearch_index_health{namespace="foo"} 0
    `
	if err := testutil.CollectAndCompare(testUpdater.gaugeVec, strings.NewReader(expected)); err != nil {
		t.Error(err)
	}
}

func TestUpdateRolloverAttemptFailures(t *testing.T) {
	testUpdater := testProvideUpdater("rollover_attempt_health")
	testUpdater.UpdateRolloverAttemptFailures(TestASs)
	expected := `
        # HELP opensearch_rollover_attempt_health test
        # TYPE opensearch_rollover_attempt_health gauge
        opensearch_rollover_attempt_health{namespace="bar"} 0
		opensearch_rollover_attempt_health{namespace="boo"} 1
		opensearch_rollover_attempt_health{namespace="foo"} 0
    `
	if err := testutil.CollectAndCompare(testUpdater.gaugeVec, strings.NewReader(expected)); err != nil {
		t.Error(err)
	}
}

var TestSCs = []models.StatusChange{
	{
		Alias:                     "foo",
		DocsAdded:                 3,
		NewIndexOperationFailures: 0,
	},
	{
		Alias:                     "bar",
		DocsAdded:                 2,
		NewIndexOperationFailures: 1,
	},
	{
		Alias:                     "boo",
		DocsAdded:                 1,
		NewIndexOperationFailures: 2,
	},
}

func TestUpdateDocsAddedRate(t *testing.T) {
	testUpdater := testProvideUpdater("alias_rate")
	testUpdater.UpdateDocsAddedRate(TestSCs)
	expected := `
        # HELP opensearch_alias_rate test
        # TYPE opensearch_alias_rate gauge
        opensearch_alias_rate{namespace="bar"} 2
		opensearch_alias_rate{namespace="boo"} 1
		opensearch_alias_rate{namespace="foo"} 3
    `
	if err := testutil.CollectAndCompare(testUpdater.gaugeVec, strings.NewReader(expected)); err != nil {
		t.Error(err)
	}
}

func TestUpdateIndexOperationFailures(t *testing.T) {
	testUpdater := testProvideUpdater("index_operation_failures")
	testUpdater.UpdateIndexOperationFailures(TestSCs)
	expected := `
        # HELP opensearch_index_operation_failures test
        # TYPE opensearch_index_operation_failures gauge
        opensearch_index_operation_failures{namespace="bar"} 1
		opensearch_index_operation_failures{namespace="boo"} 2
		opensearch_index_operation_failures{namespace="foo"} 0
    `
	if err := testutil.CollectAndCompare(testUpdater.gaugeVec, strings.NewReader(expected)); err != nil {
		t.Error(err)
	}
}
