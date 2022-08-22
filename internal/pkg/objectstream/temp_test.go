package objectstream

import (
	"fmt"
	"testing"
)

func TestNewTempPutStream(t *testing.T) {
	stream, err := NewTempPutStream("http://127.0.0.1:8080/temp", "abc.1", 7)
	if err != nil {
		fmt.Println(err)
		return
	}
	//_ = stream
	stream.Write([]byte("1232131"))
	stream.Commit(true)
}
