package restream

import (
	"fmt"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/pkg/objectstream"
	"golang-object-storage/internal/pkg/rs"
	"io"
	"strconv"
)

type RSGetStream struct {
	*rs.Decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	AllShards := global.RsConfig.AllShards
	DataShards := global.RsConfig.DataShards

	if len(locateInfo)+len(dataServers) != AllShards {
		return nil, fmt.Errorf("Error: dataServer number is not equal to %d\n", AllShards)
	}
	readers := make([]io.Reader, AllShards)
	writers := make([]io.Writer, AllShards)
	perShard := (size + int64(DataShards) - 1) / int64(DataShards)
	for shardIdx := 0; shardIdx < AllShards; shardIdx++ {
		server, ok := locateInfo[shardIdx]
		if !ok {
			writer, err := objectstream.NewTempPutStream(
				dataServers[0]+"/temp",
				hash+"."+strconv.Itoa(shardIdx),
				perShard)
			dataServers = dataServers[1:]
			if err == nil {
				writers[shardIdx] = writer
			}
		} else {
			reader, err := objectstream.NewGetStream(server + "/objects?filename=" + hash)
			if err == nil {
				readers[shardIdx] = reader
			}
		}
	}
	dec := rs.NewDecoder(readers, writers, size, global.RsConfig)
	return &RSGetStream{dec}, nil
}

// 虚假seek 读入服务服务的数据并丢弃
func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("Error: Only support SeekCurrent")
	}
	if offset < 0 {
		panic("Error: offset should not be lower than 0")
	}
	// 每次读取BlockSize字节并丢弃
	for offset != 0 {
		length := s.BlockSize
		if length > offset {
			length = offset
		}
		offset -= length
		buff := make([]byte, length)
		io.ReadFull(s, buff)
	}
	return offset, nil
}

// 修复文件上传提交
func (s *RSGetStream) Close() {
	for i := 0; i < len(s.Writers); i++ {
		if s.Writers[i] != nil {
			s.Writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}
