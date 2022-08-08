package heartbeat

import (
	"golang-object-storage/inernal/dataserver/global"
	"golang-object-storage/inernal/pkg/rabbitmq"
	"os"
	"time"
)

func StartHeartbeat() {
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	for {
		mq.Publish("apiServers", global.ListenAddr)
		time.Sleep(5 * time.Second)
	}
}
