package heartbeat

import (
	"golang-object-storage/internal/apiserver/reedso"
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

// ChooseRandomDataServer
// @author: caicandong
// @date: 2022-08-16 14:09:18
// @Description:
// @return string
func ChooseRandomDataServer() string {

	dataServers := GetAliveDataServers()
	serverCount := len(dataServers)

	if serverCount == 0 {
		return ""
	}

	log.Println("Alive data servers:", strings.Join(dataServers, ", "))
	return dataServers[rand.Intn(serverCount)]
}

// ChooseServers
// @author: caicandong
// @date: 2022-08-16 14:10:16
// @Description:
// @param dataServersNum  需要的数据服务节点个
// @param exclude 已经存放数据服务的节点
// @return dataServers
func ChooseServers(dataServersNum int, exclude map[int]string) (dataServers []string) {
	// 所需的用于存储的分片的节点数与已存放正常分片数据的节点数之和应等于一个对象的分片数之和，否则应直接中断程序执行
	if dataServersNum+len(exclude) != reedso.ALL_SHARDS {
		panic("apiServer Error: the sum of brokenShards number and unbrokenShards number is not equal to ALL_SHARDS\n")
	}
	candiateServers := make([]string, 0, dataServersNum)
	reverseUnbrokenShardMap := make(map[string]int)
	for id, serverAddr := range exclude {
		reverseUnbrokenShardMap[serverAddr] = id
	}

	aliveServers := GetAliveDataServers()
	for i := 0; i < len(aliveServers); i++ {
		if _, in := reverseUnbrokenShardMap[aliveServers[i]]; !in {
			candiateServers = append(candiateServers, aliveServers[i])
		}
	}
	if len(candiateServers) < dataServersNum {
		return
	}
	randomIds := rand.Perm(len(candiateServers))
	for i := 0; i < dataServersNum; i++ {
		dataServers = append(dataServers, candiateServers[randomIds[i]])
	}
	return
}
