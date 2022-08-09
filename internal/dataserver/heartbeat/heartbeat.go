package heartbeat

import (
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/pkg/rabbitmq"
	"os"
	"time"
)

// StartHeartbeat 给apiServers发送心跳信息
// 心跳信息为当前节点的网络地址
func StartHeartbeat() {
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	for {
		mq.Publish("apiServers", global.ListenAddr)
		time.Sleep(5 * time.Second)
	}
}
