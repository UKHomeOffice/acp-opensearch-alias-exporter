package retriever

import (
	"encoding/json"
	"fmt"
)

type WriteIndex struct {
	Alias string
	Name  string
}

type IndexGetter interface {
	GetWriteIndexes() ([]WriteIndex, error)
}

type indexGetter struct {
	host     string
	username string
	password string
	getter   HttpGetter
}

func NewIndexGetter(host, username, password string, getter HttpGetter) IndexGetter {
	return &indexGetter{
		host:     host,
		username: username,
		password: password,
		getter:   getter,
	}
}

func (i *indexGetter) GetWriteIndexes() ([]WriteIndex, error) {
	url := fmt.Sprintf("%s/_alias", i.host)
	body, err := i.getter(url, i.username, i.password)
	if err != nil {
		return nil, err
	}
	var data map[string]map[string]map[string]map[string]bool
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	aliasIndexes := make([]WriteIndex, 0)
	for index := range data {
		if len(data[index]["aliases"]) > 0 {
			for alias, aliasIndexSettings := range data[index]["aliases"] {
				if aliasIndexSettings["is_write_index"] == true {
					a := WriteIndex{
						Alias: alias,
						Name:  index,
					}
					aliasIndexes = append(aliasIndexes, a)
				}
			}
		}
	}
	return aliasIndexes, nil
}
