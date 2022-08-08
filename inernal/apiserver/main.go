package main

import (
	"flag"
	"golang-object-storage/inernal/apiserver/global"
	"golang-object-storage/inernal/apiserver/heartbeat"
	"golang-object-storage/inernal/apiserver/locate"
	"golang-object-storage/inernal/apiserver/objects"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// 加载配置文件
func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}
}

func main() {
	flag.StringVar(&global.ListenAddr, "listenAddr", ":8089", "")
	flag.Parse()
	// global.CheckSharedVars()

	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
