package temp

import (
	"encoding/json"
	"golang-object-storage/internal/dataserver/global"
	"golang-object-storage/internal/dataserver/locate"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// 临时文件信息结构体
type tempInfo struct {
	UUID     string `json:"uuid"` //文件的临时uuid
	Hash     string `json:"hash"` //文件的唯一表示hash code
	SharpIdx int    `json:"sharp_idx"`
	Size     int64  `json:"size"` //文件的大小信息
}

func NewTempInfoFromFile(uuid string) (*tempInfo, error) {
	info := &tempInfo{UUID: uuid}
	err := info.readFromFile()
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (t *tempInfo) hash() string {
	return strings.Split(t.Hash, ".")[0]
}

func (t *tempInfo) id() int {
	id, _ := strconv.Atoi(strings.Split(t.Hash, ".")[1])
	return id
}

func (t *tempInfo) getInfoFilePath() string {
	return path.Join(global.StoragePath, "temp", t.UUID)
}

func (t *tempInfo) getTempFilePath() string {
	return t.getInfoFilePath() + ".dat"
}

// 写入存放对象元数据的临时文件
func (t *tempInfo) writeToFile() error {
	file, err := os.Create(path.Join(global.StoragePath, "temp", t.UUID))
	if err != nil {
		return err
	}
	defer file.Close()
	bytesData, _ := json.Marshal(t)
	_, err = file.Write(bytesData)
	if err != nil {
		return err
	}
	return nil
}

// 读取存放对象元数据的临时文件
func (t *tempInfo) readFromFile() error {
	file, err := os.Open(t.getInfoFilePath())
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	tempInfoBytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(tempInfoBytes, t)
	if err != nil {
		return err
	}
	return nil
}

// 删除存放对象元数据的临时文件
func (t *tempInfo) RemoveInfoFromFile() error {
	err := os.Remove(t.getInfoFilePath())
	return err
}

// 删除存放对象数据的临时文件
func (t *tempInfo) removeTempFromFile() error {
	err := os.Remove(t.getTempFilePath())
	return err
}

// 删除存放对象数据的临时文件(包括元数据和数据)
func (t *tempInfo) removeAllTempFromFile() error {
	err := t.RemoveInfoFromFile()
	if err != nil {
		return err
	}
	err = t.removeTempFromFile()
	if err != nil {
		return err
	}
	return nil
}

// 将临时文件移动至节点内部
func (t *tempInfo) commitTempObject(tempFilePath string) {
	err := os.Rename(tempFilePath, path.Join(global.StoragePath, "objects", t.Hash+"."+strconv.Itoa(t.SharpIdx)))
	if err != nil {
		panic(err)
	}
	locate.AddNewObject(t.Hash, t.SharpIdx)
}

// CleanTemp 每隔12小时清理一次临时文件
func CleanTemp() {
	time.Sleep(12 * time.Hour)
	tmpDir := path.Join(global.StoragePath, "temp")
	files, _ := ioutil.ReadDir(tmpDir)
	for i := 0; i < len(files); i++ {
		dif := int(files[i].ModTime().Sub(time.Now()).Minutes())
		if dif >= 30 {
			os.Remove(path.Join(global.StoragePath, "temp", files[i].Name()))
		}
	}
}
