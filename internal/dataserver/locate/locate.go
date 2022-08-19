package locate

import (
	"fmt"
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/pkg/rabbitmq"
	"os"
	"strconv"
)

type LocateReplyMessage struct {
	Addr string
	ID   int
}

// ListenLocate
// 所有的数据服务节点绑定这个exchange并接收来自接口服务的定位消息。
// *拥有该对象的数据服务节点则使用消息单发通知该接口服务节点*
func ListenLocate() {
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	mq.BindExchange("dataServers")
	channel := mq.Consume()
	for msg := range channel {
		fmt.Println(string(msg.Body))
		hash, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		// err_utils.PanicNonNilError(err)
		if ObjectExists(hash) {
			mq.Send(msg.ReplyTo, LocateReplyMessage{
				Addr: global.ListenAddr,
				ID:   defaultRecord.records[hash],
			})
		}
	}
}
