/*
使用io.Pipe()同步管道，完成文件读写的同步
通过chan完成错误信息的传递
*/
package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer    *io.PipeWriter
	errorChan chan error
}

func NewPutStream(upload_url string) *PutStream {
	r, w := io.Pipe()
	errorChan := make(chan error)

	go func() {
		// 根据 io.Pipe()的规则,此处会发生阻塞，等待w完成写入。
		request, _ := http.NewRequest(http.MethodPut, upload_url, r)
		httpClient := http.Client{}
		response, err := httpClient.Do(request)
		if err == nil && response.StatusCode != http.StatusOK {
			err = fmt.Errorf("Error: [dataServer] statusCode: %d\n", response.StatusCode)
		}
		errorChan <- err
	}()

	return &PutStream{
		writer:    w,
		errorChan: errorChan,
	}
}

func (putStream *PutStream) Write(p []byte) (n int, err error) {
	return putStream.writer.Write(p)
}

func (putStream *PutStream) Close() error {
	putStream.writer.Close()
	return <-putStream.errorChan
}
