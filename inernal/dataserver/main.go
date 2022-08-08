package main

import (
	"golang-object-storage/inernal/dataserver/global"
	"golang-object-storage/inernal/dataserver/heartbeat"
	"golang-object-storage/inernal/dataserver/locate"
	"golang-object-storage/inernal/dataserver/objects"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
)

// 加载配置文件
func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}
}

func main() {
	// flag.IntP("flagname", "f", 1234, "help message")
	flag.StringVar(&global.ListenAddr, "listenAddr", ":8080", "listen address ")
	// flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	flag.StringVar(&global.StoragePath, "storageRoot", "static/objects", "storage root directory")

	flag.Parse()

	go heartbeat.StartHeartbeat()
	go locate.ListenLocate()

	http.HandleFunc("/objects/", objects.Handler)
	log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))

}
