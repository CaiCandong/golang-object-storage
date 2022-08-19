package hashutils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 获取文件的md5码
func GetFileHash(filename string) (string, error) {
	// 文件全路径名
	path := fmt.Sprintf("./%s", filename)
	pFile, err := os.Open(path)
	if err != nil {
		fmt.Errorf("打开文件失败，filename=%v, err=%v", filename, err)
		return "", err
	}
	defer pFile.Close()
	hash := CalculateHash(pFile)
	return hash, err
}
