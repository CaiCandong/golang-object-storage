package locate

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	location := Locate(strings.Split(r.RequestURI, "/")[2])
	if len(location) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	locationJson, err := json.Marshal(location)
	if err != nil {
		panic(err)
	}
	w.Write(locationJson)
}
