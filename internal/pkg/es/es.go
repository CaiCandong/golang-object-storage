package es

import (
	"github.com/olivere/elastic/v7"
	"log"
)

// 连接池管理
var serverEsClientMap = make(map[string]*elastic.Client)

func GetEsClient(esServer string) *elastic.Client {
	if serverEsClientMap[esServer] != nil {
		return serverEsClientMap[esServer]
	}
	newClient, err := elastic.NewClient(
		elastic.SetURL(esServer),
		//在Docker容器中,需要设置sniff为false
		//ref: https://github.com/olivere/elastic/wiki/Connection-Problems#how-to-figure-out-connection-problems
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Printf("ES Error: failed to new es client for server [%s]\n, details: %s\n", esServer, err.Error())
		return nil
	}
	serverEsClientMap[esServer] = newClient
	return newClient
}
