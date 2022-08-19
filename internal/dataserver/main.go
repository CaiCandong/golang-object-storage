package main

import (
	"fmt"
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/dataserver/heartbeat"
	"golang-object-storage/internal/dataserver/locate"
	"golang-object-storage/internal/dataserver/objects"
	"golang-object-storage/internal/dataserver/temp"
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

var help bool

func main() {
	flag.Usage = func() {
		fmt.Println(`Usage: main [OPTIONS] `)
		flag.PrintDefaults()
	}
	flag.StringVarP(&global.ListenAddr, "port", "p", ":8080", "listen address ")
	flag.BoolVarP(&help, "help", "h", false, "Print this help message")
	flag.StringVarP(&global.StoragePath, "storageRoot", "s", "static", "storage root directory")

	// 初始化默认的locate字典
	locate.DefaultFileHashRecord()
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	go heartbeat.StartHeartbeat()
	go locate.ListenLocate()

	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	fmt.Printf("listen port:%s ,storage directory: %s\n", global.ListenAddr, global.StoragePath)
	log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
