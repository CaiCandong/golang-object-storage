package elasticsearch

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetEsClient(t *testing.T) {
	esClient := getEsClient("http://127.0.0.1:9200")
	fmt.Print(esClient)
}

func TestMetadataExists(t *testing.T) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}

	metadataExists("abc.txt", 1, int64(1234), "bdfa")
}

func TestPutMetadata(t *testing.T) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}

	PutMetadata("abc.txt", 2, int64(1234), "bdfa")
}
func TestGetMetadata(t *testing.T) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}

	metadata, err := GetMetadata("abc.txt", 1)
	if err != nil {
		return
	}
	fmt.Println(metadata)
}

func TestSearchAllVersions(t *testing.T) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}

	metadata, err := SearchAllVersions("abc.txt", 0, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(metadata))
}
