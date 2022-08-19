package main

import (
	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/apiserver/heartbeat"
	"golang-object-storage/internal/apiserver/index"
	"golang-object-storage/internal/apiserver/locate"
	"golang-object-storage/internal/apiserver/objects"
	"golang-object-storage/internal/apiserver/versions"
	"log"
	"net/http"
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
	http.HandleFunc("/index/", index.Handler)
	http.HandleFunc("/objects", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
