package locate

import (
	"encoding/json"
	"golang-object-storage/internal/apiserver/reedso"
	"golang-object-storage/internal/pkg/rabbitmq"
	"os"
	"time"
)

type LocateMessage struct {
	Addr string
	ID   int
}

// Locate 通过向dataServers发送请求查询文件位置
// 并通过dataServers消息队列等待相应
func Locate(objectName string) (locateInfo map[int]string) {
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))

	mq.Publish("dataServers", objectName)
	channel := mq.Consume()
	// Publish()后，设置超时关闭连接，以判断资源是否存在
	go func() {
		time.Sleep(1 * time.Second)
		mq.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < reedso.ALL_SHARDS; i++ {
		msg := <-channel
		if len(msg.Body) == 0 {
			return
		}
		var info LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.ID] = info.Addr
	}
	return
}

func Exist(objectName string) bool {
	return len(Locate(objectName)) >= reedso.DATA_SHARDS
}
