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
		} `json:"primaries"`
	} `json:"_all"`
}

type aliasgetter struct {
	host     string
	username string
	password string
	getter   HttpGetter
}

func (a *aliasgetter) GetAlias(index string, name string) (models.AliasStatus, error) {
	url := fmt.Sprintf("%s/%s/_stats", a.host, index)
	body, err := a.getter(url, a.username, a.password)

	var stats Stats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return models.AliasStatus{}, err
	}

	alias := models.AliasStatus{
		Count:  stats.All.Primaries.Docs.Count,
		Index:  index,
		Name:   name,
		Getter: a,
	}
	return alias, nil
}

func NewAliasGetter(host, username, password string, getter HttpGetter) models.AliasGetter {
	return &aliasgetter{host: host, username: username, password: password, getter: getter}
}
