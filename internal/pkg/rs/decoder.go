package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

// 完成数据解码 以及可能的分片修复
type Decoder struct {
	*Config
	// 为什么不仅需要readers，还需要writers?因为在读取数据的同时需要进行可能的数据修复
	Readers   []io.Reader         // 可正常读且数据完好取对象分片的数据节点的文件读对象
	Writers   []io.Writer         // 不可正常读取或数据缺失的对象分片的数据节点的文件写对象
	rsEnc     reedsolomon.Encoder // reedsolomon中对对象进行分片、编码、解码及数据恢复都需要依靠该对象进行
	size      int64               // 对象数据的大小，也就是数据分片中的实际数据量
	cache     []byte              // 用于缓存读取的数据
	cacheSize int                 // 用于计算缓存了多少数据
	total     int64               // 用于读取数据时进行计数
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64, config *Config) *Decoder {
	if config == nil {
		config = DefaultConfig()
	}
	enc, _ := reedsolomon.New(config.DataShards, config.ParityShards)
	return &Decoder{
		Config:    config,
		Readers:   readers,
		Writers:   writers,
		rsEnc:     enc,
		size:      size,
		cache:     make([]byte, 0, config.BlockPerShard),
		cacheSize: 0,
		total:     0,
	}
}

// 将解码后的数据cache中的数据 写入 p []byte
func (d *Decoder) Read(p []byte) (n int, err error) {
	// 当缓存中没有数据时，会通过调用getData()获取数据
	if d.cacheSize == 0 {
		err := d.getData()
		if err != nil {
			return 0, err
		}
	}
	dataLength := len(p)
	if d.cacheSize < dataLength {
		dataLength = d.cacheSize
	}
	d.cacheSize -= dataLength
	copy(p, d.cache[:dataLength])
	d.cache = d.cache[dataLength:]
	return dataLength, nil
}

// 对数据进行解码和恢复 每次处理一个 BLOCK_PER_SHARD写入缓存cache中
func (d *Decoder) getData() error {
	// 如果当前rsDecoder读取的数据总量total已达到对象数据大小，则直接返回已读完
	if d.total == d.size {
		return io.EOF
	}

	// 读取rsDecoder中的readers序列的文件读对象，将其数据写入到对应的字节序列中
	// 若是某一文件读对象为空则说明该编号的对象分片缺失，需要进行修复，并将其放入repairIds中
	shards := make([][]byte, d.AllShards)
	repairIds := make([]int, 0)
	for i := 0; i < len(shards); i++ {
		if d.Readers[i] == nil {
			repairIds = append(repairIds, i)
		} else {
			shards[i] = make([]byte, d.BlockPerShard)
			readCount, err := io.ReadFull(d.Readers[i], shards[i])
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if int64(readCount) != d.BlockPerShard {
				shards[i] = shards[i][:readCount]
			}
		}
	}
	// 如果存在需要修复的分片，则进行分片修复
	if len(repairIds) > 0 {
		err := d.rsEnc.Reconstruct(shards)
		if err != nil {
			return err
		}
		// 将恢复的分片写入对应的服务节点
		for i := 0; i < len(repairIds); i++ {
			id := repairIds[i]
			d.Writers[id].Write(shards[id])
		}
	}
	// 解码数据分片，还原数据
	for i := 0; i < d.DataShards; i++ {
		shardSize := int64(len(shards[i]))
		// 如果处理到最后一块数据分片时，存在数据填充，则只取实际数据
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		// 将数据分片的数据存入缓存中，同时计算缓存数据总量
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}
	return nil
}
