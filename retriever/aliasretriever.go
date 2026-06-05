package retriever

import (
	"encoding/json"
	"fmt"
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
)

type Stats struct {
	Shards struct {
        Failed int `json:"failed"`
    } `json:"_shards"`
	All struct {
		Primaries struct {
			Docs struct {
				Count int `json:"count"`
			} `json:"docs"`
			Indexing struct {
                IndexFailed int `json:"index_failed"`
            } `json:"indexing"`
		} `json:"primaries"`
	} `json:"_all"`
}

type ISM struct {
	Action struct {
		Name string `json:"name"`
		Failed bool `json:"failed"`
	} `json:"action"`
}

type aliasgetter struct {
	host     string
	username string
	password string
	getter   HttpGetter
}

func (a *aliasgetter) GetAlias(index string, name string) (models.AliasStatus, error) {
	stats_url := fmt.Sprintf("%s/%s/_stats", a.host, index)
	stats_body, stats_err := a.getter(stats_url, a.username, a.password)

	var stats Stats
	stats_err = json.Unmarshal(stats_body, &stats)
	if stats_err != nil {
		return models.AliasStatus{}, stats_err
	}

	ism_url := fmt.Sprintf("%s/_plugins/_ism/explain/%s", a.host, index)
	ism_body, ism_err := a.getter(ism_url, a.username, a.password)

	var ism ISM
	ism_err = json.Unmarshal(ism_body, &ism)
	if ism_err != nil {
		return models.AliasStatus{}, ism_err
	}

	alias := models.AliasStatus{
		Count:  stats.All.Primaries.Docs.Count,
		FailedShards: stats.Shards.Failed,
		FailedIndexOperations: stats.All.Primaries.Indexing.IndexFailed,
		RolloverAttemptFailed: ism.Action.Failed,
		Index:  index,
		Name:   name,
		Getter: a,
	}
	return alias, nil
}

func NewAliasGetter(host, username, password string, getter HttpGetter) models.AliasGetter {
	return &aliasgetter{host: host, username: username, password: password, getter: getter}
}
