package objects

import (
	"fmt"
	"golang-object-storage/internal/apiserver/heartbeat"
	"golang-object-storage/internal/apiserver/locate"
	"golang-object-storage/internal/apiserver/reedso"
	"golang-object-storage/internal/pkg/hashutils"
	"log"
	"net/url"

	"io"
	"net/http"
)

// StoreObject 将文件进行存储
func StoreObject(reader io.Reader, hash string, size int64) (statusCode int, err error) {
	escapedHash := url.PathEscape(hash)
	// 若是对象的内容数据已存在，则不用重复上传，否则将对象数据保存到临时缓存中等待校验
	if locate.Exist(escapedHash) {
		return http.StatusOK, nil
	}
	stream, err := putStream(escapedHash, size)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	r := io.TeeReader(reader, stream)
	actualHash := hashutils.CalculateHash(r)
	// 进行文件存储
	if actualHash != hash {
		stream.Commit(false)
		err = fmt.Errorf("Error: object hashutils value is not match, actualHash=[%s], expectedHash=[%s]\n", actualHash, hash)
		return http.StatusBadRequest, err
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*reedso.RSPutStream, error) {
	servers := heartbeat.ChooseServers(reedso.ALL_SHARDS, nil)
	if len(servers) != reedso.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	log.Printf("apiServer INFO: Choose random data servers to save object %s: %v\n", hash, servers)
	return reedso.NewRSPutStream(servers, hash, size)
}

func LoadObject(writer http.ResponseWriter, hash string, size int64) (statusCode int, err error) {
	stream, err := getStream(hash, size)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}
	io.Copy(writer, stream)
	// 保证能够正常解码数据后，在将修复的数据保存到数据节点中
	stream.Close()
	return http.StatusOK, nil
}

func getStream(objectName string, size int64) (*reedso.RSGetStream, error) {
	locateInfo := locate.Locate(objectName)
	if len(locateInfo) < reedso.DATA_SHARDS {
		return nil, fmt.Errorf("Error: object %s locate failed, the data shards located is not enough: %v\n",
			objectName, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) < reedso.ALL_SHARDS {
		log.Printf("INFO: some of shards need to repair\n")
		dataServers = heartbeat.ChooseServers(reedso.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return reedso.NewRSGetStream(locateInfo, dataServers, objectName, size)
}
