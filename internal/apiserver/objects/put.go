package objects

import (
	"fmt"
	"golang-object-storage/internal/apiserver/datalocate"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/apiserver/restream"
	"golang-object-storage/internal/pkg/utils"
	"io"
	"log"
	"net/http"
	"net/url"
)

// StoreObject 将文件进行存储
func StoreObject(reader io.Reader, hash string, size int64) (statusCode int, err error) {
	escapedHash := url.PathEscape(hash)
	// 若是对象的内容数据已存在，则不用重复上传，否则将对象数据保存到临时缓存中等待校验
	if datalocate.Exist(escapedHash) {
		log.Printf("locate file [%s] success", escapedHash)
		return http.StatusOK, nil
	}

	stream, err := putStream(escapedHash, size)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	r := io.TeeReader(reader, stream)
	actualHash := utils.CalculateHash(r)
	// 进行文件存储
	if actualHash != hash {
		stream.Commit(false)
		err = fmt.Errorf("Error: object hashutils value is not match, actualHash=[%s], expectedHash=[%s]\n", actualHash, hash)
		return http.StatusBadRequest, err
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*restream.RSPutStream, error) {
	servers := datalocate.ChooseServers(global.RsConfig.AllShards, nil)
	if len(servers) != global.RsConfig.AllShards {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	log.Printf("apiServer INFO: Choose random data servers to save object %s: %v\n", hash, servers)
	return restream.NewRSPutStream(servers, hash, size)
}
