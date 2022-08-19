package reedso

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

const (
	// 数据分片
	DATA_SHARDS = 4
	// 校验分片
	PARITY_SHARDS = 2
	// 校验分片与数据分片的总和，该总和不应超出可用的服务节点
	ALL_SHARDS = DATA_SHARDS + PARITY_SHARDS
	// 每个分片一次性可写入的最大数据量
	BLOCK_PER_SHARD = 8000
	// 每个对象的所有分片一次性可写入的最大数据量，超过该数据量则要分批写入
	BLOCK_SIZE = BLOCK_PER_SHARD * DATA_SHARDS
)

// rsEncoder
// @Description: 封装reedsolomon库,实现容错码的业务要求，增加一个写入缓冲区
type rsEncoder struct {
	writers  []io.Writer         // 写入指针
	cache    []byte              // 写入缓冲区
	rsEncode reedsolomon.Encoder //算法逻辑
}

func NewRSEncoder(writers []io.Writer) *rsEncoder {
	rsEnc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &rsEncoder{
		writers:  writers,
		rsEncode: rsEnc,
		cache:    make([]byte, 0),
	}
}

func (rsEnc *rsEncoder) Write(p []byte) (n int, err error) {
	// 维护大小为Block_size的滑窗
	start := 0
	// 每次循环写入一个Block_size大小
	for {
		if start+BLOCK_SIZE > len(p) {
			// TODO:不足一个BLOCK_SIZE的数据块保存在缓冲区中未进行flush
			rsEnc.cache = append(rsEnc.cache, p[start:]...)
			break
		}
		rsEnc.cache = append(rsEnc.cache, p[start:start+BLOCK_SIZE]...)
		start += BLOCK_SIZE
		rsEnc.Flush()
	}
	return len(p), nil
}

func (rsEnc *rsEncoder) Flush() {
	if len(rsEnc.cache) == 0 {
		return
	}
	// 划分为4个数据片
	shards, _ := rsEnc.rsEncode.Split(rsEnc.cache)
	// 调用算法生产2个校验片
	rsEnc.rsEncode.Encode(shards)
	for i := 0; i < len(shards); i++ {
		rsEnc.writers[i].Write(shards[i])
	}
	rsEnc.cache = []byte{}
}

type rsDecoder struct {
	// 为什么不仅需要readers，还需要writers?因为在读取数据的同时需要进行可能的数据修复
	readers   []io.Reader         // 可正常读且数据完好取对象分片的数据节点的文件读对象
	writers   []io.Writer         // 不可正常读取或数据缺失的对象分片的数据节点的文件写对象
	rsEnc     reedsolomon.Encoder // reedsolomon中对对象进行分片、编码、解码及数据恢复都需要依靠该对象进行
	size      int64               // 对象数据的大小，也就是数据分片中的实际数据量
	cache     []byte              // 用于缓存读取的数据
	cacheSize int                 // 用于计算缓存了多少数据
	total     int64               // 用于读取数据时进行计数
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *rsDecoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &rsDecoder{
		readers:   readers,
		writers:   writers,
		rsEnc:     enc,
		size:      size,
		cache:     make([]byte, 0, BLOCK_PER_SHARD),
		cacheSize: 0,
		total:     0,
	}
}

// 【读取并解码数据】从cache中读取数据
func (rsDec *rsDecoder) Read(p []byte) (n int, err error) {
	// 当缓存中没有数据时，会通过调用getData()获取数据
	if rsDec.cacheSize == 0 {
		err := rsDec.getData()
		if err != nil {
			return 0, err
		}
	}
	dataLength := len(p)
	if rsDec.cacheSize < dataLength {
		dataLength = rsDec.cacheSize
	}
	rsDec.cacheSize -= dataLength
	copy(p, rsDec.cache[:dataLength])
	rsDec.cache = rsDec.cache[dataLength:]
	return dataLength, nil
}

func (rsDec *rsDecoder) getData() error {
	// 如果当前rsDecoder读取的数据总量total已达到对象数据大小，则直接返回已读完
	if rsDec.total == rsDec.size {
		return io.EOF
	}

	shards := make([][]byte, ALL_SHARDS)
	repairIds := make([]int, 0)
	for shardIdx := 0; shardIdx < len(shards); shardIdx++ {
		if rsDec.readers[shardIdx] == nil {
			repairIds = append(repairIds, shardIdx)
		} else {
			shards[shardIdx] = make([]byte, BLOCK_PER_SHARD)
			readCount, err := io.ReadFull(rsDec.readers[shardIdx], shards[shardIdx])
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				shards = nil
			} else if readCount != BLOCK_PER_SHARD {
				shards[shardIdx] = shards[shardIdx][:readCount]
			}
		}
	}

	if len(repairIds) > 0 {
		err := rsDec.rsEnc.Reconstruct(shards)
		if err != nil {
			return err
		}
		for i := 0; i < len(repairIds); i++ {
			id := repairIds[i]
			rsDec.writers[id].Write(shards[id])
		}
	}

	for i := 0; i < DATA_SHARDS; i++ {
		shardSize := int64(len(shards[i]))
		// 去除数据填充
		if rsDec.total+shardSize > rsDec.size {
			shardSize -= rsDec.total + shardSize - rsDec.size
		}
		rsDec.cache = append(rsDec.cache, shards[i][:shardSize]...)
		rsDec.total += shardSize
	}
	return nil
}
