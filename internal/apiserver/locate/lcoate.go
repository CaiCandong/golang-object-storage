package locate

import (
	"golang-object-storage/internal/pkg/rabbitmq"
	"log"
	"os"
	"strconv"
	"time"
)

// Locate 通过向dataServers发送请求查询文件位置
// 并通过dataServers消息队列等待相应
func Locate(objectName string) string {
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))

	mq.Publish("dataServers", objectName)
	channel := mq.Consume()
	// Publish()后，设置超时关闭连接，以判断资源是否存在
	go func() {
		time.Sleep(1 * time.Second)
		mq.Close()
	}()
	//lcoateInfo := make(map[int]string)
	//for i:=0;i<
	// 准备接收消息
	msg := <-channel
	result, err := strconv.Unquote(string(msg.Body))
	// TODO:处理文件不存在情况
	if err != nil {
		panic(err)
	}
	// err_utils.PanicNonNilError(err)
	log.Printf("INFO: object at server '%s'\n", result)

	return result
}

func Exist(objectName string) bool {
	return Locate(objectName) != ""
}
