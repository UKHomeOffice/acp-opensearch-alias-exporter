package retriever

import (
	"encoding/json"
	"fmt"
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
)

type Stats struct {
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
		Name   string `json:"name"`
		Failed bool   `json:"failed"`
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
	stats_body, stats_err := a.getter(stats_url, a.username, a.password) // example response: https://docs.opensearch.org/2.19/api-reference/index-apis/stats/#example-response

	var stats Stats
	stats_err = json.Unmarshal(stats_body, &stats)
	if stats_err != nil {
		return models.AliasStatus{}, stats_err
	}

	ism_url := fmt.Sprintf("%s/_plugins/_ism/explain/%s", a.host, index)
	ism_body, err := a.getter(ism_url, a.username, a.password) // example response: https://docs.opensearch.org/2.19/im-plugin/ism/api/#example-response-12
	if err != nil {
		return models.AliasStatus{}, err
	}

	var ismResponse map[string]json.RawMessage
	if err := json.Unmarshal(ism_body, &ismResponse); err != nil {
		return models.AliasStatus{}, err
	}

	var ism ISM
	if raw, ok := ismResponse[index]; ok {
		if err := json.Unmarshal(raw, &ism); err != nil {
			return models.AliasStatus{}, err
		}
	}

	alias := models.AliasStatus{
		DocCount:              stats.All.Primaries.Docs.Count,
		FailedIndexOperations: stats.All.Primaries.Indexing.IndexFailed,
		RolloverAttemptFailed: ism.Action.Failed,
		Index:                 index,
		Name:                  name,
		Getter:                a,
	}
	return alias, nil
}

func NewAliasGetter(host, username, password string, getter HttpGetter) models.AliasGetter {
	return &aliasgetter{host: host, username: username, password: password, getter: getter}
}
