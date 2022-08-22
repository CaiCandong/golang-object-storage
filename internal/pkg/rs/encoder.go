package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

// p []byte 中的数据分片写入到 Writers[]io.Writer
type Encoder struct {
	*Config
	Writers  []io.Writer         // io.Writer接口，用于将对象分片数据写入到指定的存储位置
	rsEncode reedsolomon.Encoder // reedsolomon用于编码的调用对象
	cache    []byte              // 缓存待写入的数据，大小一般为BLOCK_SIZE，一次性可写入ALL_SHARDS个切片的数据
}

func NewRSEncoder(writers []io.Writer, config *Config) *Encoder {
	if config == nil {
		config = DefaultConfig()
	}
	rsEnc, _ := reedsolomon.New(config.DataShards, config.ParityShards)
	return &Encoder{
		Config:   config,
		Writers:  writers,
		rsEncode: rsEnc,
		cache:    make([]byte, 0),
	}
}

func (e *Encoder) Write(p []byte) (n int, err error) {
	residueLength := int64(len(p)) //剩余待处理长度
	var current int64 = 0
	for residueLength != 0 {
		//delta := BLOCK_SIZE - len(rsEnc.cache) //当前处理长度
		delta := e.BlockSize - int64(len(e.cache)) //当前处理长度
		if delta > residueLength {
			delta = residueLength
		}
		//delta:= min(BLOCK_SIZE,residueLength)
		e.cache = append(e.cache, p[current:current+delta]...)
		// 若是cache的数据量已达到一次性可写入的最大的数据量，则先将该部分数据写入数据服务节点中，然后再读取下一批
		if int64(len(e.cache)) == e.BlockSize {
			e.Flush()
		}
		current += delta
		residueLength -= delta
	}

	return len(p), nil
}

func (e *Encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	// 分片、编码、写入对应的服务节点的文件中
	shards, _ := e.rsEncode.Split(e.cache)
	e.rsEncode.Encode(shards)
	for i := 0; i < len(shards); i++ {
		e.Writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}
