/*
	使用io.Pipe()同步管道，完成文件读写的同步
	通过chan完成错误信息的传递
*/
package objects

import (
	"fmt"
	"golang-object-storage/internal/apiserver/locate"
	"golang-object-storage/internal/pkg/objectstream"
	"io"
)

func getStream(objectName string) (io.Reader, error) {
	server := locate.Locate(objectName)
	if server != "" {
		return nil, fmt.Errorf("ERROR: object '%s' not found", objectName)
	}
	return NewGetStream(server, objectName)
}

func NewGetStream(server, objectName string) (*objectstream.GetStream, error) {
	if server == "" || objectName == "" {
		return nil, fmt.Errorf("Value Error: server='%s', objectName='%s'\n", server, objectName)
	}
	return objectstream.NewGetStream(fmt.Sprintf("http://%s/objects/%s", server, objectName))
}

type PutStream struct {
	writer    *io.PipeWriter
	errorChan chan error
}

func (putStream *PutStream) Write(p []byte) (n int, err error) {
	return putStream.writer.Write(p)
}

func (putStream *PutStream) Close() error {
	putStream.writer.Close()
	return <-putStream.errorChan
}
