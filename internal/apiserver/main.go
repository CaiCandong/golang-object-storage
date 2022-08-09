package main

import (
<<<<<<< HEAD
<<<<<<< HEAD
	flag "github.com/spf13/pflag"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/apiserver/heartbeat"
	"golang-object-storage/internal/apiserver/index"
=======
	"flag"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/apiserver/heartbeat"
>>>>>>> 60c4855 (build(chapter03): 第三章代码)
=======
	flag "github.com/spf13/pflag"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/apiserver/heartbeat"
	"golang-object-storage/internal/apiserver/index"
>>>>>>> f301f56 (feat✨:  chapter3)
	"golang-object-storage/internal/apiserver/locate"
	"golang-object-storage/internal/apiserver/objects"
	"golang-object-storage/internal/apiserver/versions"
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
	http.HandleFunc("/index/", index.Handler)
	http.HandleFunc("/objects", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
