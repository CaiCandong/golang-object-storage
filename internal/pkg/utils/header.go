package utils

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetHashFromHeader(h http.Header) string {
	const prefix = "SHA-256="
	digest := h.Get("digest")
	// 存放hash值的参数名设为SHA-256，因此若是hash值为空或者参数名对应不上，则直接返回空串
	if !strings.HasPrefix(digest, prefix) {
		return ""
	}
	return digest[len(prefix):]
}

func GetOffsetFromHeader(h http.Header) (offset, end int64) {
	byteRange := h.Get("range")
	const prefix = "bytes="
	log.Printf("GetOffsetFromHeader():range content[%s]\n", byteRange)
	if !strings.HasPrefix(byteRange, prefix) {
		return 0, 0
	}
	bytesPositions := strings.Split(byteRange[len(prefix):], "-")
	log.Printf("%v\n", bytesPositions)
	offset, _ = strconv.ParseInt(bytesPositions[0], 0, 64)
	end, _ = strconv.ParseInt(bytesPositions[1], 0, 64)
	log.Printf("offset[%d]-end[%d]\n", offset, end)

	return offset, end
}
