package v1

import (
	"encoding/json"
	"golang-object-storage/internal/apiserver/datalocate"
	"net/http"
	"net/url"
)

func LocateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("filename")
	name = url.PathEscape(name)
	location := datalocate.Locate(name)
	if len(location) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	locationJson, _ := json.Marshal(location)
	w.Write(locationJson)
}
