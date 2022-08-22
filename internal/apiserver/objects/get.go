package objects

import (
	"fmt"
	"golang-object-storage/internal/apiserver/datalocate"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/apiserver/restream"
	"io"
	"log"
	"net/http"
)

func LoadObject(w http.ResponseWriter, hash string, size int64, offset, end int64) (statusCode int, err error) {
	stream, err := getStream(hash, size)
	if err != nil {
		log.Printf("get download stream fail : %s", err)
		return http.StatusBadRequest, err
	}
	contentLength := size
	if offset != 0 {
		contentLength = end - offset + 1
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, end, size))
		w.WriteHeader(http.StatusPartialContent)
	}
	writeen, err := io.CopyN(w, stream, contentLength)
	if err != nil {
		log.Printf("get file stream fail %s", err)
	}
	log.Println("wrote to response length:", writeen)
	stream.Close()
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}
	io.Copy(w, stream)
	// 保证能够正常解码数据后，在将修复的数据保存到数据节点中
	stream.Close()
	return http.StatusOK, nil
}

func getStream(objectName string, size int64) (*restream.RSGetStream, error) {
	locateInfo := datalocate.Locate(objectName)
	if len(locateInfo) < global.RsConfig.DataShards {
		global.Logger.Infof("Error: object %s locate failed, the data shards located is not enough: %v\n",
			objectName, locateInfo)
		return nil, fmt.Errorf("Error: object %s locate failed, the data shards located is not enough: %v\n",
			objectName, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) < global.RsConfig.AllShards {
		log.Printf("INFO: some of shards need to repair\n")
		dataServers = datalocate.ChooseServers(global.RsConfig.AllShards-len(locateInfo), locateInfo)
	}
	return restream.NewRSGetStream(locateInfo, dataServers, objectName, size)
}
