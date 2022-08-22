package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"golang-object-storage/internal/pkg/es"
	"log"
	"os"
)

// Metadata 主键 Name + Version
type Metadata struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	Size    int64  `json:"size"`
	Hash    string `json:"hash"`
}

// 增
func (m *Metadata) Put() error {
	if m.Exists() {
		m.Version += 1
		return m.Put()
	}

	esClient := es.GetEsClient(os.Getenv("ES_SERVER"))
	searchResult, err := esClient.Index().
		Index("metadata").
		Id(fmt.Sprintf("%s_%d", m.Name, m.Version)).
		BodyJson(*m).
		Refresh("wait_for").
		Do(context.Background())
	_ = searchResult
	//fmt.Println(searchResult)
	return err
}

// Get 查
func (m *Metadata) Get(name string, version int) error {
	esClient := es.GetEsClient(os.Getenv("ES_SERVER"))
	nameQuery, versionQuery := elastic.NewTermQuery("name", name), elastic.NewTermQuery("version", version)
	searchResult, err := esClient.Search().
		Index("metadata").
		Query(nameQuery).Query(versionQuery).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return err
	}
	if searchResult.TotalHits() == 0 {
		return fmt.Errorf("name = %s version = %d metadata no found", name, version)
	}
	dataBytes, err := searchResult.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataBytes, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Metadata) Exists() bool {
	esClient := es.GetEsClient(os.Getenv("ES_SERVER"))
	nameQuery := elastic.NewTermQuery("name", m.Name)
	versionQuery := elastic.NewTermQuery("version", m.Version)
	//sizeQuery := elastic.NewTermQuery("size", m.Size)
	//hashQuery := elastic.NewTermQuery("hash", m.Hash)
	searchResult, err := esClient.Search().
		Index("metadata").
		Query(nameQuery).Query(versionQuery).
		//Query(sizeQuery).Query(hashQuery).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		log.Println(err)
		return false
	}
	return searchResult.TotalHits() != 0
}
