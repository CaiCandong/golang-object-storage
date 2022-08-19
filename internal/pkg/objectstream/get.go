package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

// NewGetStream
// @author: caicandong
// @date: 2022-08-16 20:15:47
// @Description:
// @param url 为请求文件的网络地址
// @return *GetStream
// @return error
func NewGetStream(download_url string) (*GetStream, error) {
	response, err := http.Get(download_url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: [dataServer] statusCode: %d\n", response.StatusCode)
	}
	return &GetStream{reader: response.Body}, nil
}

func (getStream *GetStream) Read(p []byte) (n int, err error) {
	return getStream.reader.Read(p)
}
