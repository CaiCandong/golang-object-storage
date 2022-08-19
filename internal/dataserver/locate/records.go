package locate

import (
	"golang-object-storage/internal/dataserver/global"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type FileHashRecords struct {
	records map[string]int // hash => shard_idx (每个数据服务点只能保存一个分片)
	mutex   sync.Mutex
}

var defaultRecord *FileHashRecords

func NewFileHashRecords(pattern string) *FileHashRecords {
	r := &FileHashRecords{}
	r.records = make(map[string]int)
	files, _ := filepath.Glob(pattern)
	for i := 0; i < len(files); i++ {
		// file_name: <file_hash>.<shard_idx>.<shard_hash>
		shardNameComponents := strings.Split(filepath.Base(files[i]), ".")
		fileHash := shardNameComponents[0]
		shardIdx, err := strconv.Atoi(shardNameComponents[1])
		if err != nil {
			log.Fatalf("Error: shard %v name is invalid, it should be 3 compoments [objectHash.ID.shardHash]\n", shardNameComponents)
		}
		r.records[fileHash] = shardIdx
	}
	return r
}

func (r *FileHashRecords) ObjectExists(hash string) bool {
	r.mutex.Lock()
	_, ok := r.records[hash]
	r.mutex.Unlock()
	return ok
}

func (r *FileHashRecords) AddNewObject(hash string, shardIdx int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.records[hash] = shardIdx
}

func (r *FileHashRecords) Delete(hash string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.records, hash)
}

func DefaultFileHashRecord() {
	pattern := path.Join(global.StoragePath, "objects", "*")
	defaultRecord = NewFileHashRecords(pattern)
}

func ObjectExists(hash string) bool {
	return defaultRecord.ObjectExists(hash)
}

func AddNewObject(hash string, shardIdx int) {
	defaultRecord.AddNewObject(hash, shardIdx)
}

func Delete(hash string) {
	defaultRecord.Delete(hash)
}
