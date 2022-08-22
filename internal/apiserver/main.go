package main

import (
	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	v1 "golang-object-storage/internal/apiserver/api/v1"
	"golang-object-storage/internal/apiserver/datalocate"
	"golang-object-storage/internal/apiserver/global"
	"log"
	"net/http"
)

func init() {
	// 加载配置文件
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("godotenv Error: env files load failed")
	}
	// 日志
	sugar := zap.NewExample().Sugar()
	//logger, _ := zap.NewProduction(zap.AddCaller())
	global.Logger = sugar
	defer sugar.Sync()

}

func main() {
	flag.StringVar(&global.ListenAddr, "listenAddr", ":8089", "")
	flag.Parse()

	go datalocate.ListenHeartbeat()
	http.HandleFunc("/objects", v1.ObjectHandler)
	http.HandleFunc("/locate/", v1.VersionHandler)
	http.HandleFunc("/versions/", v1.LocateHandler)
	log.Fatalln(http.ListenAndServe(global.ListenAddr, nil))
}
