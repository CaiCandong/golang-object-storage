package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
)

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 获取文件的md5码
func GetFileHash(filepath string) string {
	// 文件全路径名
	pFile, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer pFile.Close()
	hash := CalculateHash(pFile)
	return hash
}
