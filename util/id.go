package util

import (
	"fmt"
	"hash/crc32"
	"time"
)

func GenerateId(prefix string, data interface{}) string {
	return fmt.Sprintf("%s%d%d", prefix, crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v", data)))%1000, time.Now().In(time.Local).UnixNano())
}
