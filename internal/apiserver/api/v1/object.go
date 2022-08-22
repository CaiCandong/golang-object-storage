package v1

import (
	"encoding/json"
	"fmt"
	"golang-object-storage/internal/apiserver/datalocate"
	"golang-object-storage/internal/apiserver/metadata"
	"golang-object-storage/internal/apiserver/objects"
	"golang-object-storage/internal/pkg/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPost {
		post(w, r)
	} else if m == http.MethodPut {
		put(w, r)
	} else if m == http.MethodGet {
		get(w, r)
	} else if m == http.MethodDelete {
		del(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func del(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("filename")
	m := &metadata.Metadata{}
	err := m.SearchLatestVersion(name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 逻辑删除：将size和hash置空即可
	m.Size = 0
	m.Version += 1
	m.Hash = ""
	err = m.Put()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Query().Get("filename")
	versionID := r.URL.Query().Get("version")
	var version int
	var err error
	// 如果有version参数则查找指定版本的对象，否则查找最新版本对象
	if len(versionID) != 0 {
		version, err = strconv.Atoi(versionID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	// 从ES获取对象元数据信息，进而通过元数据新信息的hash值向数据服务层请求对象内容
	var m metadata.Metadata
	err = m.Get(objectName, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if m.Hash == "" {
		log.Printf("ES INFO: object [%s] not found", objectName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 获取对象数据
	// 解析头部的range字段
	offset, end := utils.GetOffsetFromHeader(r.Header)
	log.Printf("apiServer INFO: in get(), get object data range [%d, %d]\n", offset, end)

	statusCode, err := objects.LoadObject(w, url.PathEscape(m.Hash), m.Size, offset, end)
	if statusCode != http.StatusOK || err != nil {
		w.WriteHeader(statusCode)
		return
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Query().Get("filename")
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Printf("apiServer Error: missing object [%s] hash\n", objectName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 文件已经存在
	if datalocate.Exist(url.PathEscape(hash)) {
		m := &metadata.Metadata{}
		err = m.AddVersion(objectName, size, hash)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

}

func put(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	// hashVal 哈希值
	hashVal := utils.GetHashFromHeader(r.Header)
	if hashVal == "" {
		log.Println("API-Server HTTP Error: missing object hash in request header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// size 上传文件大小
	size, err := strconv.ParseInt(r.Header.Get("content-length"), 0, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	if size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("API-Server HTTP Error: missing object file in request body")
		err = fmt.Errorf("API-Server HTTP Error: missing object file in request body")
		return
	}
	//name 文件名
	name := r.URL.Query().Get("filename")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("API-Server HTTP Error: missing object filename in request url")
		err = fmt.Errorf("API-Server HTTP Error: missing object filename in request url")
	}
	// 交由下层负责文件存储
	statusCode, err := objects.StoreObject(r.Body, hashVal, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(statusCode)
	}
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}
	// 保存文件元信息
	m := metadata.Metadata{}
	err = m.AddVersion(name, size, hashVal)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, _ := json.Marshal(m)
	w.Write(data)
}
