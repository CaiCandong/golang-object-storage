package datalocate

import (
	"encoding/json"
	"fmt"
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/pkg/rabbitmq"
	"log"
	"os"
	"time"
)

type LocateMessage struct {
	Addr string
	ID   int
}

// Locate 通过向dataServers发送请求查询文件位置
// 并通过dataServers消息队列等待相应
func Locate(objectName string) map[int]string {
	rsconfig := global.RsConfig
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	global.Logger.Info(fmt.Sprintf("Locate File [%s]", objectName))
	mq.Publish("dataServers", objectName)
	channel := mq.Consume()
	// Publish()后，设置超时关闭连接，以判断资源是否存在
	go func() {
		time.Sleep(1 * time.Second)
		mq.Close()
	}()
	locateInfo := make(map[int]string)
	for i := 0; i < rsconfig.AllShards; i++ {
		msg := <-channel
		if len(msg.Body) == 0 {
			return locateInfo
		}
		var info LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.ID] = info.Addr
	}
	return locateInfo
}

func Exist(objectName string) bool {
	rsconfig := global.RsConfig
	exist := len(Locate(objectName)) >= rsconfig.DataShards
	if !exist {
		log.Printf("locate file[%s] fail", objectName)
	}
	return exist
}
