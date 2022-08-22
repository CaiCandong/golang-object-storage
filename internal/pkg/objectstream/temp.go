package objectstream

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string // 数据服务接口地址
	Name   string // 文件名
	Size   int64  // 文件大小
	UUID   string // 临时文件标识
}

// 临时文件上传握手
func (t *TempPutStream) post() error {
	url := t.Server + "/" + t.Name
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("size", fmt.Sprintf("%d", t.Size))
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	log.Printf("post file to %s\n", url)
	if err != nil {
		return err
	}
	uuidBytes, err := ioutil.ReadAll(response.Body)
	uuidBytes = bytes.ReplaceAll(uuidBytes, []byte("\n"), []byte(""))
	t.UUID = string(uuidBytes)
	return nil
}

// 临时文件上传
func (t *TempPutStream) patch(p []byte) error {
	url := t.Server + "/" + t.UUID
	request, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(string(p)))
	if err != nil {
		return err
	}
	httpClient := http.Client{}
	log.Printf("patch file to %s\n", url)
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("dataServer Error: STATUSCODE[%d]\n", response.StatusCode)
	}
	return nil
}

// 临时文件上传完成
func (t *TempPutStream) put() {
	url := t.Server + "/" + t.UUID
	request, _ := http.NewRequest(http.MethodPut, url, nil)
	httpClient := http.Client{}
	httpClient.Do(request)
}

// 临时文件上传撤销
func (t *TempPutStream) del() {
	url := t.Server + "/" + t.UUID
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	httpClient := http.Client{}
	httpClient.Do(request)
}

// 临时文件下载
func (t *TempPutStream) get() (*GetStream, error) {
	url := t.Server + "/" + t.UUID
	return NewGetStream(url)
}
