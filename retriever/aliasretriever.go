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
	// Gets docs and index operation failures from _stats
	statsUrl := fmt.Sprintf("%s/%s/_stats", a.host, index)
	statsBody, statsErr := a.getter(statsUrl, a.username, a.password) // https://docs.opensearch.org/2.19/api-reference/index-apis/stats/#example-response
	if statsErr != nil {
		return models.AliasStatus{}, statsErr
	}

	var stats Stats
	statsErr = json.Unmarshal(statsBody, &stats)
	if statsErr != nil {
		return models.AliasStatus{}, statsErr
	}

	// Gets rolloever attempt failures from _ism
	ismUrl := fmt.Sprintf("%s/_plugins/_ism/explain/%s", a.host, index)
	ismBody, err := a.getter(ismUrl, a.username, a.password) // https://docs.opensearch.org/2.19/im-plugin/ism/api/#example-response-12
	if err != nil {
		return models.AliasStatus{}, err
	}

	var ismResponse map[string]json.RawMessage
	if err := json.Unmarshal(ismBody, &ismResponse); err != nil {
		return models.AliasStatus{}, err
	}

	var ism ISM
	if raw, ok := ismResponse[index]; ok {
		if err := json.Unmarshal(raw, &ism); err != nil {
			return models.AliasStatus{}, err
		}
	}

	// Gets index health from _cluster/health/
	healthUrl := fmt.Sprintf("%s/_cluster/health/%s", a.host, index)
	healthBody, err := a.getter(healthUrl, a.username, a.password) // https://docs.opensearch.org/2.19/api-reference/cluster-api/cluster-health/#example-response
	if err != nil {
		return models.AliasStatus{}, err
	}

	var healthResponse map[string]json.RawMessage
	if err := json.Unmarshal(healthBody, &healthResponse); err != nil {
		return models.AliasStatus{}, err
	}

	var health string
	raw, ok := healthResponse["status"]
	if !ok {
		return models.AliasStatus{}, fmt.Errorf("status field not found in health response")
	}
	if err := json.Unmarshal(raw, &health); err != nil {
		return models.AliasStatus{}, err
	}

	alias := models.AliasStatus{
		DocCount:              stats.All.Primaries.Docs.Count,
		FailedIndexOperations: stats.All.Primaries.Indexing.IndexFailed,
		RolloverAttemptFailed: ism.Action.Failed,
		Index:                 index,
		Name:                  name,
		Getter:                a,
		Health:                health,
	}

	return alias, nil
}

func NewAliasGetter(host, username, password string, getter HttpGetter) models.AliasGetter {
	return &aliasgetter{host: host, username: username, password: password, getter: getter}
}
