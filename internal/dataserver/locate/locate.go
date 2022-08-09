package locate

import (
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/pkg/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
)

// ListenLocate
//所有的数据服务节点绑定这个exchange并接收来自接口服务的定位消息。
// *拥有该对象的数据服务节点则使用消息单发通知该接口服务节点*
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
			// 消息单发通知该接口服务节点
			mq.Send(msg.ReplyTo, global.ListenAddr)
		}
	}
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
