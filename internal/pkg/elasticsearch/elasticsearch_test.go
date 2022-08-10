package elasticsearch

import (
	"testing"
)

func TestXxx(t *testing.T) {
	esClient := getEsClient("localhost:9200")
	print(esClient)
}
