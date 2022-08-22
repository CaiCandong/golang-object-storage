package objects

import (
	"fmt"
	"strings"
	"testing"
)

func TestStoreObject(t *testing.T) {
	//go dataserver.ListenHeartbeat()
	r := strings.NewReader("你吃饭了嘛?")

	statusCode, err := StoreObject(r, "saf", 12)
	if err != nil {
		return
	}
	fmt.Println(statusCode)
}
