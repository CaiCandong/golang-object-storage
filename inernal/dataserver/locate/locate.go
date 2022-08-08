package locate

import (
	"golang-object-storage/inernal/dataserver/global"
	"golang-object-storage/inernal/pkg/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
)

func ListenLocate() {
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	mq.BindExchange("dataServers")
	channel := mq.Consume()
	for msg := range channel {
		objectName, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		// err_utils.PanicNonNilError(err)

		filePath := filepath.Join(global.StoragePath, "objects", objectName)
		if pathExist(filePath) {
			mq.Send(msg.ReplyTo, global.ListenAddr)
		}
	}
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
