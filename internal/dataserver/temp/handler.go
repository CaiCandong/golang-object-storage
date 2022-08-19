package temp

import (
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/pkg/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodPost {
		// 1.发post 获取uuid
		post(w, r)
		return
	}
	if method == http.MethodPatch {
		// 2. 发patch 完成文件上传
		patch(w, r)
		return
	}
	if method == http.MethodPut {
		// 3.临时文件转正
		put(w, r)
		return
	}
	if method == http.MethodDelete {
		// 4.md5校验不通过时,删除临时文件
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 创建存储对象元数据的临时文件，同时返回临时文件的uuid
func post(w http.ResponseWriter, r *http.Request) {
	// 注意：产生的uuid值末尾会携带一个换行符，因此必须去除换行符
	uUid := uuid.GenUUid()
	name := GetObjectName(r.URL.EscapedPath())
	components := strings.Split(name, ".")
	hash := components[0]
	sharpIdx, err := strconv.Atoi(components[1])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temp := tempInfo{
		UUID:     uUid,
		Hash:     hash,
		SharpIdx: sharpIdx,
		Size:     size,
	}
	// 缓存对象的临时元数据
	err = temp.writeToFile()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 创建用于存放对象数据的临时文件，这个与对象的临时元数据文件的作用不一样，前者用于标志临时对象所在的服务节点，后者用于存放对象的内容数据
	file, err := os.Create(path.Join(global.StoragePath, "temp", temp.UUID+".dat"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	w.Write([]byte(uUid))
}

// 缓存对象数据，初步校验数据的大小是否匹配
func patch(w http.ResponseWriter, r *http.Request) {
	uuid := GetObjectName(r.URL.EscapedPath())
	// 获取uuid对应的临时文件存放的对象元数据信息，用于校验实际上传的数据的信息是否正确
	tempinfo, err := NewTempInfoFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 将实际从http中获取的数据写入data文件中
	// infoFile和dataFile是在post请求时创建的，infoFile存放了对象的元数据，dataFile用于存放对象内容数据
	// 读取对象数据文件的内容
	dataFile := tempinfo.getTempFilePath()

	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// 写入存储对象数据的临时文件中
	io.Copy(file, r.Body)
	// 比较实际获取的对象数据文件的大小与期望的对象数据大小是否相同
	// 若不相同，则删除创建的临时文件：对象元数据文件、对象数据文件
	actualInfo, err := os.Stat(dataFile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ac_size := actualInfo.Size(); ac_size != tempinfo.Size {
		tempinfo.removeAllTempFromFile()
		if err != nil {
			panic(err)
		}

		log.Printf("Error: the actual uploaded file`s size [%d] is dismatched with expected size [%d]\n",
			actualInfo.Size(), tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// 确认完成文件上传,临时文件转为正式文件(移动文件位置,并删除文件元信息)
func put(w http.ResponseWriter, r *http.Request) {
	name := GetObjectName(r.URL.EscapedPath())
	var tempinfo tempInfo
	tempinfo.UUID = name
	// 从临时缓存区获取存储对象元数据的临时文件
	err := tempinfo.readFromFile()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 读取对象数据文件的内容,检查文件是否存在
	dataFile := tempinfo.getTempFilePath()
	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// 校验临时缓存区的对象元数据与对象数据的大小是否匹配（一般是匹配的，并且哈希值校验在数据接口层就已完成）
	actualInfo, err := os.Stat(dataFile)
	tempinfo.RemoveInfoFromFile()
	if actualInfo.Size() != tempinfo.Size {
		file.Close()
		tempinfo.removeTempFromFile()
		log.Printf("Error: the actual uploaded file`s size [%d] is dismatched with expected size [%d]\n",
			actualInfo.Size(), tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file.Close()
	// 数据转正
	tempinfo.commitTempObject(dataFile)
}

// 若是对象数据校验未通过则删除缓存区的临时文件
func del(w http.ResponseWriter, r *http.Request) {
	// 获取文件uuid
	uuid := GetObjectName(r.URL.EscapedPath())
	// 新建文件信息结构体
	tempinfo, err := NewTempInfoFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//删除临时文件信息
	tempinfo.removeAllTempFromFile()
}
