package main

import (
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/retriever"
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/updater"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

func getCurrentAliasStatuses() (models.AliasStatuses, error) {
	indexGetter := retriever.NewIndexGetter(os.Getenv("OPENSEARCH_HOST"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), retriever.HttpGet)
	aliasGetter := retriever.NewAliasGetter(os.Getenv("OPENSEARCH_HOST"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), retriever.HttpGet)

	writeIndexes, err := indexGetter.GetWriteIndexes()
	if err != nil {
		return nil, err
	}

	aliasStatuses := make(models.AliasStatuses)
	for _, index := range writeIndexes {
		aliasStatus, err := aliasGetter.GetAlias(index.Name, index.Alias)
		if err != nil {
			log.Println("Error getting alias", err)
		} else {
			aliasStatuses[index.Alias] = aliasStatus
		}
	}
	return aliasStatuses, nil
}

func start() {

	oldAliasStatuses, err := getCurrentAliasStatuses()
	if err != nil {
		panic(err)
	}

	prometheusHealthUpdater := updater.NewPrometheusUpdater("opensearch", "index_health", "tracks health of index based on whether primary shards and their replicas are allocated to nodes.")
	prometheusRolloverAttemptUpdater := updater.NewPrometheusUpdater("opensearch", "rollover_attempt_health", "tracks index rollover attempt failures")
	prometheusIndexOperationFailureUpdater := updater.NewPrometheusUpdater("opensearch", "index_operation_failures", "tracks new failed index operations")
	prometheusAliasRateUpdater := updater.NewPrometheusUpdater("opensearch", "alias_rate", "tracks rate of new documents added to alias")

	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		newAliasStatuses, err := getCurrentAliasStatuses()
		if err != nil {
			panic(err)
		}
		prometheusHealthUpdater.UpdateHealth(newAliasStatuses)
		prometheusRolloverAttemptUpdater.UpdateRolloverAttemptFailures(newAliasStatuses)

		statusChanges, err := models.GetStatusChanges(oldAliasStatuses, newAliasStatuses)
		if err != nil {
			log.Println("Error getting status changes", err)
		}
		prometheusAliasRateUpdater.UpdateDocsAddedRate(statusChanges)
		prometheusIndexOperationFailureUpdater.UpdateIndexOperationFailures(statusChanges)

		oldAliasStatuses = newAliasStatuses
	}
}

func main() {

	go func() {
		start()
	}()

	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":8080", nil)
}
