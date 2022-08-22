package datalocate

import (
	"golang-object-storage/internal/apiserver/global"
	"golang-object-storage/internal/pkg/rabbitmq"
	"math/rand"
	"os"
	"strconv"
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

func GetAliveDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()

	dataServers := make([]string, 0)
	for server, _ := range dataServerMap {
		dataServers = append(dataServers, server)
	}

	return dataServers
}

// 获取dataServersNum个服务节点,且不能获得idx2server内的服务节点
func ChooseServers(dataServersNum int, idx2server map[int]string) []string {
	// 所需的用于存储的分片的节点数与已存放正常分片数据的节点数之和应等于一个对象的分片数之和，否则应直接中断程序执行
	if dataServersNum+len(idx2server) != global.RsConfig.AllShards {
		panic("apiServer Error: the sum of brokenShards number and unbrokenShards number is not equal to ALL_SHARDS\n")
	}
	candidateServers := make([]string, 0, dataServersNum)
	server2idx := make(map[string]int)
	for id, serverAddr := range idx2server {
		server2idx[serverAddr] = id
	}

	aliveServers := GetAliveDataServers()
	global.Logger.Infof("alive dataservers: %v", aliveServers)
	for i := 0; i < len(aliveServers); i++ {
		if _, in := server2idx[aliveServers[i]]; !in {
			candidateServers = append(candidateServers, aliveServers[i])
		}
	}
	if len(candidateServers) < dataServersNum {
		return nil
	}
	// 顺序打乱
	randomIds := rand.Perm(len(candidateServers))
	dataServers := make([]string, 0)
	for i := 0; i < dataServersNum; i++ {
		dataServers = append(dataServers, candidateServers[randomIds[i]])
	}
	return dataServers
}
