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

	prometheusAliasRateUpdater := updater.NewPrometheusUpdater("opensearch", "alias_rate", "rate of change of alias count")
	prometheusRolloverHealthUpdater := updater.NewPrometheusUpdater("opensearch", "alias_rollover_health", "tracks index rollover failures/health")
	prometheusShardHealthUpdater := updater.NewPrometheusUpdater("opensearch", "alias_shard_health", "tracks unhealthy alias shards")
	prometheusIndexOperationHealthUpdater := updater.NewPrometheusUpdater("opensearch", "alias_index_operation_health", "tracks failed index operations")

	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		newAliasStatuses, err := getCurrentAliasStatuses()
		if err != nil {
			panic(err)
		}
		
		prometheusRolloverHealthUpdater.UpdateRolloverHealth(newAliasStatuses)
		prometheusShardHealthUpdater.UpdateShardHealth(newAliasStatuses)
		prometheusIndexOperationHealthUpdater.UpdateIndexOperationHealth(newAliasStatuses)

		countRates, err := models.GetCountChanges(oldAliasStatuses, newAliasStatuses)
		if err != nil {
			log.Println("Error getting count rates", err)
		}
		prometheusAliasRateUpdater.UpdateCountRates(countRates)

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
