package locate

import (
	"golang-object-storage/internal/dataserver/global"
	"path"
	"path/filepath"
	"sync"
)

type FileHashRecords struct {
	records map[string]bool
	mutex   sync.Mutex
}

var defaultRecord *FileHashRecords

func NewFileHashRecords(pattern string) *FileHashRecords {
	r := &FileHashRecords{}
	r.records = make(map[string]bool)
	files, _ := filepath.Glob(pattern)
	for i := 0; i < len(files); i++ {
		hash := filepath.Base(files[i])
		r.records[hash] = true
	}
	return r
}

func (r *FileHashRecords) ObjectExists(hash string) bool {
	r.mutex.Lock()
	_, ok := r.records[hash]
	r.mutex.Unlock()
	return ok
}

func (r *FileHashRecords) AddNewObject(hash string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.records[hash] = true
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

func AddNewObject(hash string) {
	defaultRecord.AddNewObject(hash)
}

func Delete(hash string) {
	defaultRecord.Delete(hash)
}
