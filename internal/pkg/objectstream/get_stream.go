package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func NewGetStream(download_url string) (*GetStream, error) {
	//global.Logger.Info(fmt.Sprintf("HTTP GET:%s", download_url))
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
