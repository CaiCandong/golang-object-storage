package objects

import (
	"io"
	"log"
	"net/http"
	"os"
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
		put(w, r)
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

func get(w http.ResponseWriter, r *http.Request) {
	objectName := strings.Split(r.RequestURI, uriSep)[2]
	file, err := os.Open(filepath.Join(os.Getenv(storageRootEnvName), objectParentDirName, objectName))
	if err != nil {
		log.Println("GET FAILED:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	io.Copy(w, file)
	log.Printf("GET SUCCESS: object '%s'\n", file.Name())
}
