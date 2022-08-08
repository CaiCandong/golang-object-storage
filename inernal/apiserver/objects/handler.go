package objects

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m == http.MethodPut {
		put(w, r)
	} else if m == http.MethodGet {
		get(w, r)
	} else {
		// 如果不是以上请求方法的任一种，则返回405
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	objectName := GetObjectName(r.URL.EscapedPath())
	statusCode, err := storeObject(r.Body, objectName)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(statusCode)
}

func get(w http.ResponseWriter, r *http.Request) {
	objectName := GetObjectName(r.URL.EscapedPath())
	stream, err := getStream(objectName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, stream)
}

func GetObjectName(url string) string {
	url = strings.TrimSpace(url)
	components := strings.Split(url, "/")

	return components[len(components)-1]
}
