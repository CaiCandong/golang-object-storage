package objects

import (
	"fmt"
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/dataserver/locate"
	"golang-object-storage/internal/pkg/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	get(w, r)
}

// 根据文件的hash值 获取文件
func get(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("filename")
	check := r.URL.Query().Get("check")

	filePath := GetFileAndCheckHash(hash, check != "")
	if filePath == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file, _ := os.Open(filePath)
	defer file.Close()
	io.Copy(w, file)
}

// GetFileAndCheckHash
// 验证文件是否完整,若完整返回路径,不完整直接删除返回空串
func GetFileAndCheckHash(hash string, check bool) string {
	pattern := path.Join(global.StoragePath, "objects", fmt.Sprintf("%s.*", url.PathEscape(hash)))
	//filename : <total_hash>.<shard_idx>.<shard_hash>
	files, _ := filepath.Glob(pattern)
	if len(files) != 1 {
		return ""

	}
	shardFileName := files[0]
	if check {
		shardExceptHash := strings.Split(shardFileName, ".")[2]
		shardActualHash := url.PathEscape(utils.GetFileHash(shardFileName))
		if shardActualHash != shardExceptHash {
			log.Printf("Shard file content hash: %s, expected hash: %s\n", shardActualHash, shardExceptHash)
			// 删除已损坏的分片的定位信息
			locate.Delete(hash)
			// 删除文件
			os.Remove(shardFileName)
			return ""
		}
	}
	return shardFileName
}
