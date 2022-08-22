package main

import (
	"fmt"
	"go.uber.org/zap"
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/dataserver/heartbeat"
	"golang-object-storage/internal/dataserver/locate"
	"golang-object-storage/internal/dataserver/objects"
	"golang-object-storage/internal/dataserver/temp"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
)

// 加载配置文件
func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}

	// 日志
	logger, _ := zap.NewProduction(zap.AddCaller())
	global.Logger = logger
	defer logger.Sync()
}

var help bool

func main() {
	flag.Usage = func() {
		fmt.Println(`Usage: main [OPTIONS] `)
		flag.PrintDefaults()
	}
	flag.StringVarP(&global.ListenAddr, "port", "p", ":8080", "listen address ")
	flag.BoolVarP(&help, "help", "h", false, "Print this help message")
	flag.StringVarP(&global.StoragePath, "storageRoot", "s", "static/8080", "storage root directory")

	flag.Parse()
	// 初始化默认的locate字典
	locate.DefaultFileHashRecord()
	port := global.ListenAddr
	global.ListenAddr = "http://127.0.0.1" + global.ListenAddr
	if help {
		flag.Usage()
		return
	}

	go heartbeat.StartHeartbeat()
	go locate.ListenLocate()

	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)

	fmt.Printf("listen port:%s ,storage directory: %s\n", global.ListenAddr, global.StoragePath)
	err := ensureDir(global.StoragePath + "/temp")
	if err != nil {
		panic(err)
	}
	ensureDir(global.StoragePath + "/objects")
	log.Fatalln(http.ListenAndServe(port, nil))
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory")
		}
		return nil
	}
	return err
}
