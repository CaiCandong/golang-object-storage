package objects

import (
	"fmt"
	"golang-object-storage/internal/apiserver/heartbeat"
	"golang-object-storage/internal/apiserver/locate"
	"golang-object-storage/internal/apiserver/temp"
	"golang-object-storage/internal/pkg/fileutils"
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
	// 根据心跳信息选择存储节点
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return http.StatusInternalServerError, fmt.Errorf("Error: no alive data server\n")
	}
	log.Println("Choose random data server:", server)

	stream, err := temp.NewPutStream(server, hash, size)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r := io.TeeReader(reader, stream)
	actualHash := fileutils.CalculateHash(r)
	// 进行文件存储
	if actualHash != hash {
		stream.Commit(false)
		err = fmt.Errorf("Error: object hash value is not match, actualHash=[%s], expectedHash=[%s]\n", actualHash, hash)
		return http.StatusBadRequest, err
	}
	stream.Commit(true)

	return http.StatusOK, nil
}
