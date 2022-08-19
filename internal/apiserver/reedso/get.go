package reedso

import (
	"fmt"
	"golang-object-storage/internal/apiserver/temp"
	"golang-object-storage/internal/pkg/objectstream"
	"golang-object-storage/internal/pkg/tpfmt"
	"io"
)

// RSGetStream
// @Description:文件分块下载
type RSGetStream struct {
	*rsDecoder
}

// NewRSGetStream
// @author: caicandong
// @date: 2022-08-16 20:47:44
// @Description:
// @param locateInfo 存储数据的节点 shard_idx => data_server_host
// @param dataServers 需要上传修复数据的节点
// @param hash
// @param size
// @return *RSGetStream
// @return error
func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	// 检查总节点个数
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("Error: dataServer number is not equal to %d\n", ALL_SHARDS)
	}
	// 给缺失分块分配存储节点
	// 构造readers
	readers := make([]io.Reader, ALL_SHARDS)
	downloadUrlTpl := tpfmt.Format("http://%{{.host}}/objects/{{.hash}}.{{.shardIdx}}")
	for shardIdx := 0; shardIdx < ALL_SHARDS; shardIdx++ {
		if server, ok := locateInfo[shardIdx]; !ok {
			locateInfo[shardIdx] = dataServers[0]
			dataServers = dataServers[1:]
		} else {
			downloadUrl := downloadUrlTpl.Exec(map[string]interface{}{
				"host":     server,
				"hash":     hash,
				"shardIdx": shardIdx,
			})
			reader, err := objectstream.NewGetStream(downloadUrl)
			if err == nil {
				readers[shardIdx] = reader
			}
		}
	}
	// 上传(分片修复)
	// 构造writer
	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var err error
	uploadUrlTpl := tpfmt.Format("http://%{{.host}}/temp/{{.hash}}.{{.shardIdx}}")
	for shardIdx := 0; shardIdx < ALL_SHARDS; shardIdx++ {
		if readers[shardIdx] == nil {
			uploadUrl := uploadUrlTpl.Exec(map[string]interface{}{
				"host":     locateInfo[shardIdx],
				"hash":     hash,
				"shardIdx": shardIdx,
			})
			writers[shardIdx], err = temp.NewPutStream(uploadUrl, perShard)
			if err != nil {
				return nil, err
			}
		}
	}
	// 创建解码对象
	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil
}

// Close
// @author: caicandong
// @date: 2022-08-16 20:45:17
// @Description: 完成修复文件上传
// @receiver
func (s *RSGetStream) Close() {
	for i := 0; i < len(s.writers); i++ {
		if s.writers[i] != nil {
			s.writers[i].(*temp.PutStream).Commit(true)
		}
	}
}
