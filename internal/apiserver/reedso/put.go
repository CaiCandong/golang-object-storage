package reedso

import (
	"fmt"
	"golang-object-storage/internal/apiserver/temp"
	"io"
)

// RSPutStream
// @Description: 文件分块上传
type RSPutStream struct {
	*rsEncoder
}

// NewRSPutStream
// @author: caicandong
// @date: 2022-08-16 20:33:21
// @Description:
// @param dataServers 数据服务的地址数组
// @param objectHash
// @param objectSize
// @return *RSPutStream
// @return error
func NewRSPutStream(dataServers []string, objectHash string, objectSize int64) (*RSPutStream, error) {
	var err error
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("Error: dataServer number is not enough\n")
	}
	// 上取整
	perShardSize := (objectSize + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	for i := 0; i < len(writers); i++ {
		// 拼接文件上传地址
		uploadUrl := fmt.Sprintf("http://%v/temp/%s.%d", dataServers[i], objectHash)
		writers[i], err = temp.NewPutStream(uploadUrl, perShardSize)
		if err != nil {
			return nil, err
		}
	}
	rsEnc := NewRSEncoder(writers)
	rsPutStream := &RSPutStream{rsEnc}
	return rsPutStream, nil
}

// Commit
// @author: caicandong
// @date: 2022-08-16 20:40:00
// @Description:
// @receiver rsPutStream
// @param positive
func (rsPutStream *RSPutStream) Commit(positive bool) {
	rsPutStream.Flush()
	for i := 0; i < len(rsPutStream.writers); i++ {
		rsPutStream.writers[i].(*temp.PutStream).Commit(positive)
	}
}
