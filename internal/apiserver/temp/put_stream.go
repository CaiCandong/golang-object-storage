package temp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// PutStream 记录数据服务的临时文件信息
type PutStream struct {
	Server string // 数据服务地址
	UUID   string // 临时文件标识
}

// NewPutStream 建立与数据服务的临时文件
// "http://"+server+"/temp/"+objectName,
func NewPutStream(uploadUrl string, size int64) (*PutStream, error) {
	// 创建临时对象
	request, err := http.NewRequest(
		http.MethodPost,
		uploadUrl,
		nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	httpClient := http.Client{}
	// 执行请求后，正常情况会在在响应中返回临时对象的UUID
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	uuidBytes, err := ioutil.ReadAll(response.Body)
	// 注意：该处读取的uuid的值的末尾会有换行符，因此必须去除，否则会引起http url语法错误
	uuidBytes = bytes.ReplaceAll(uuidBytes, []byte("\n"), []byte(""))
	if err != nil {
		return nil, err
	}

	return &PutStream{
		Server: uploadUrl,
		UUID:   string(uuidBytes),
	}, nil
}

// 实现io.Writer的Write接口
func (t *PutStream) Write(p []byte) (n int, err error) {
	// 通过patch操作，将数据写入临时文件对象，并且校验数据大小，若不符合则删除为该对象创建的临时文件【该步仅校验文件大小】
	request, err := http.NewRequest(http.MethodPatch,
		"http://"+t.Server+"/temp/"+t.UUID,
		strings.NewReader(string(p)))
	if err != nil {
		return 0, err
	}
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer Error: STATUSCODE[%d]\n", response.StatusCode)
	}
	return len(p), nil
}

// Commit 根据positive决定是删除临时对象数据还是转正保存到节点中
func (t *PutStream) Commit(positive bool) {
	method := http.MethodDelete
	if positive {
		method = http.MethodPut
	}
	request, _ := http.NewRequest(method, "http://"+t.Server+"/temp/"+t.UUID, nil)
	httpClient := http.Client{}
	httpClient.Do(request)
}