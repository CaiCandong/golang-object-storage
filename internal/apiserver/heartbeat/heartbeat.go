package heartbeat

import (
	"golang-object-storage/internal/pkg/rabbitmq"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	dataServerMap = make(map[string]time.Time)
	mutex         sync.Mutex
)

func ListenHeartbeat() {
	// 连接到rabbitmq
	mq := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	// 打开apiServicers交换机
	mq.BindExchange("apiServers")
	//新建一个通道
	channel := mq.Consume()
	defer mq.Close()

	go removeExpiredServerNode()
	// 监听消息队列chan,根据消息体更新dataServerMap
	for msg := range channel {
		server, err := strconv.Unquote(string(msg.Body))
		// err_utils.PanicNonNilError(err)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		dataServerMap[server] = time.Now()
		mutex.Unlock()
	}
}

// 每隔5s扫描一遍dataServers，并清除其中超过10s没收到心跳消息的数据服务节点
func removeExpiredServerNode() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for server, heartbeatTime := range dataServerMap {
			if heartbeatTime.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServerMap, server)
			}
		}
		mutex.Unlock()
	}
}

// GetAliveDataServers 返回可用的DataServer
func GetAliveDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()

	dataServers := make([]string, 0)
	for server, _ := range dataServerMap {
		dataServers = append(dataServers, server)
	}

	return dataServers
}

func ChooseRandomDataServer() string {
	dataServers := GetAliveDataServers()
	serverCount := len(dataServers)

	if serverCount == 0 {
		return ""
	}

	log.Println("Alive data servers:", strings.Join(dataServers, ", "))
	return dataServers[rand.Intn(serverCount)]
}
