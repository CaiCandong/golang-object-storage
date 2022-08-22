package metadata

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"golang-object-storage/internal/pkg/es"
	"os"
	"reflect"
)

func (m *Metadata) SearchLatestVersion(name string) error {
	esClient := es.GetEsClient(os.Getenv("ES_SERVER"))
	searchResult, err := esClient.Search().
		Index("metadata").
		Query(elastic.NewTermQuery("name", name)).
		Sort("version", false).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return err
	}
	if searchResult.TotalHits() > 0 {
		dataBytes, err := searchResult.Hits.Hits[0].Source.MarshalJSON()
		if err != nil {
			return err
		}
		_ = json.Unmarshal(dataBytes, m)
	}
	return nil
}

func (*Metadata) SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	esClient := es.GetEsClient(os.Getenv("ES_SERVER"))
	var searchResult *elastic.SearchResult
	var err error
	searchService := esClient.Search().
		Index("metadata").
		From(from).Size(size).
		Pretty(true)
	if name == "" {
		searchResult, err = searchService.Do(context.Background())
	} else {
		nameQuery := elastic.NewTermQuery("name", name)
		searchResult, err = searchService.
			Query(nameQuery).
			Do(context.Background())
	}
	if err != nil {
		return nil, err
	}

	var metadata Metadata
	var metadatas []Metadata
	for _, item := range searchResult.Each(reflect.TypeOf(metadata)) {
		if t, ok := item.(Metadata); ok {
			metadatas = append(metadatas, Metadata{
				Name:    t.Name,
				Version: t.Version,
				Size:    t.Size,
				Hash:    t.Hash,
			})
		}
	}
	return metadatas, nil
}

func (m *Metadata) AddVersion(name string, size int64, hash string) error {
	// 查找版本号
	err := m.SearchLatestVersion(name)
	if err != nil {
		return err
	}
	m.Name = name
	m.Version += 1
	m.Size = size
	m.Hash = hash
	err = m.Put()
	if err != nil {
		return err
	}
	return nil
}
