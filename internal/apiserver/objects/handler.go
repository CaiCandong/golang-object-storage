package objects

import (
	"fmt"
	"golang-object-storage/internal/pkg/elasticsearch"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m == http.MethodPut {
		put(w, r)
	} else if m == http.MethodGet {
		get(w, r)
	} else if m == http.MethodDelete {
		del(w, r)
	} else {
		// 如果不是以上请求方法的任一种，则返回405
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// put 文件上传服务
func put(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	//hashVal 哈希值
	//name 文件名
	//size 文件大小
	parseParam := func(r *http.Request) (md5, name string, size int64, err error) {

		hashVal := GetHashFromHeader(r.Header)
		if hashVal == "" {
			log.Println("API-Server HTTP Error: missing object hash in request header")
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("API-Server HTTP Error: missing object hash in request header")
			return
		}
		// 获取size
		size, err = strconv.ParseInt(r.Header.Get("content-length"), 0, 64)
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
		// 获取文件名
		components := strings.Split(strings.TrimSpace(r.URL.EscapedPath()), "/")
		name = components[len(components)-1]
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("API-Server HTTP Error: missing object filename in request url")
			err = fmt.Errorf("API-Server HTTP Error: missing object filename in request url")
		}
		return
	}
	hashVal, name, size, err := parseParam(r)
	if err != nil {
		return
	}
	//交由下层负责文件存储
	statusCode, err := StoreObject(r.Body, hashVal, size)

	if err != nil {
		log.Println(err)
		w.WriteHeader(statusCode)
	}
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}
<<<<<<< HEAD
=======

>>>>>>> f301f56 (feat✨:  chapter3)
	err = elasticsearch.AddVersion(name, size, hashVal)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	objectName := GetObjectName(r.URL.EscapedPath())
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
	metadata, err := elasticsearch.GetMetadata(objectName, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if metadata.Hash == "" {
		log.Printf("ES INFO: object [%s] not found", objectName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 存储对象数据
	name := url.PathEscape(metadata.Hash)
	stream, err := getStream(name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}

func del(w http.ResponseWriter, r *http.Request) {
	name := GetObjectName(r.URL.EscapedPath())
	latestMetadata, err := elasticsearch.SearchLatestVersion(name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 逻辑删除：将size和hash置空即可
	err = elasticsearch.PutMetadata(name, latestMetadata.Version+1, 0, "")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetObjectName(url string) string {
	url = strings.TrimSpace(url)
	components := strings.Split(url, "/")

	return components[len(components)-1]
}

func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	// 存放hash值的参数名设为SHA-256，因此若是hash值为空或者参数名对应不上，则直接返回空串
	if len(digest) < 9 || digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}
