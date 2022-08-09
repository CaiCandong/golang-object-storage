package objects

import (
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/dataserver/locate"
	"golang-object-storage/internal/pkg/hash"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	storageRootEnvName  = "STORAGE_ROOT"
	objectParentDirName = "objects"
	uriSep              = "/"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m == http.MethodPut {
		//不允许直接上传
		//put(w, r)
	} else if m == http.MethodGet {
		get(w, r)
	} else {
		// 如果不是以上请求方法的任一种，则返回405
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	objectName := strings.Split(r.RequestURI, uriSep)[2]
	file, err := os.Create(filepath.Join(os.Getenv(storageRootEnvName), objectParentDirName, objectName))
	if err != nil {
		log.Println("PUT FAILED:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	io.Copy(file, r.Body)
	log.Printf("PUT SUCCESS: object '%s'\n", file.Name())
}

// 根据文件的hash值 获取文件
func get(w http.ResponseWriter, r *http.Request) {
	hash := strings.Split(r.RequestURI, uriSep)[2]
	file, err := os.Open(filepath.Join(os.Getenv(storageRootEnvName), objectParentDirName, hash))
	if err != nil {
		log.Println("GET FAILED:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	io.Copy(w, file)
	log.Printf("GET SUCCESS: object '%s'\n", file.Name())
}

func getFilePath(name string) string {
	filePath := path.Join(global.StoragePath, "objects", name)
	file, _ := os.Open(filePath)
	// 计算实际存储的文件的哈希值
	storedObjectHash := url.PathEscape(hash.CalculateHash(file))
	file.Close()
	// 校验：校验接口层中ES存储的哈希值与实际存储的内容的哈希值是否一致，若是发生了变化则不一致，并且删除该对象数据
	// 数据存放久了可能会发生数据降解等问题，因此有必要做一致性校验
	if storedObjectHash != name {
		log.Printf("dataServer INFO: the object`s stored in node %s is broken, we just have removed it from dataServer node.",
			global.ListenAddr)
		locate.Delete(name)
		os.Remove(filePath)
		return ""
	}
	return filePath
}
