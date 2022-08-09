package fileutils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// CalculateHash
//计算文件的hash值
func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
