package restream

import (
	"fmt"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/pkg/objectstream"
	"golang-object-storage/internal/pkg/rs"
	"io"
	"strconv"
)

type RSPutStream struct {
	*rs.Encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	AllShards := global.RsConfig.AllShards
	DataShards := global.RsConfig.DataShards

	if len(dataServers) != AllShards {
		return nil, fmt.Errorf("Error: dataServer number is not enough\n")
	}
	perShard := (size + int64(DataShards) - 1) / int64(DataShards)
	writers := make([]io.Writer, AllShards)
	for i := range writers {
		writer, err := objectstream.NewTempPutStream(
			dataServers[i]+"/temp",
			hash+"."+strconv.Itoa(i),
			perShard)
		if err != nil {
			return nil, err
		}
		writers[i] = writer
	}
	encoder := rs.NewRSEncoder(writers, global.RsConfig)
	return &RSPutStream{encoder}, nil
}

func (s *RSPutStream) Commit(positive bool) {
	s.Flush()
	for i := 0; i < len(s.Writers); i++ {
		s.Writers[i].(*objectstream.TempPutStream).Commit(positive)
	}
}
