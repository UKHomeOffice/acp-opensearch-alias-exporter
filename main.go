package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.digital.homeoffice.gov.uk/acp/opensearch-reporter/models"
	"gitlab.digital.homeoffice.gov.uk/acp/opensearch-reporter/retriever"
	"gitlab.digital.homeoffice.gov.uk/acp/opensearch-reporter/updater"
	"log"
	"net/http"
	"os"
	"time"
)

func getCurrentAliasStatuses() (models.AliasStatuses, error) {
	indexGetter := retriever.NewIndexGetter(os.Getenv("host"), os.Getenv("username"), os.Getenv("password"), retriever.HttpGet)
	aliasGetter := retriever.NewAliasGetter(os.Getenv("host"), os.Getenv("username"), os.Getenv("password"), retriever.HttpGet)

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

	prometheusUpdater := updater.NewPrometheusUpdater("opensearch", "alias_rate", "rate of change of alias count")

	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		newAliasStatuses, err := getCurrentAliasStatuses()
		if err != nil {
			panic(err)
		}

		countRates, err := models.GetCountChanges(oldAliasStatuses, newAliasStatuses)
		if err != nil {
			log.Println("Error getting count rates", err)
		}
		prometheusUpdater.Update(countRates)

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
