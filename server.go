package main

import (
	"golang-object-storage/objects"
	"log"
	"net/http"
	"os"
)

const (
	objectPattern = "/objects/"
	listenAddress = "LISTEN_ADDress"
)

func main() {
	http.HandleFunc(objectPattern, objects.Handler)
	log.Println(http.ListenAndServe(os.Getenv(listenAddress), nil))
}
