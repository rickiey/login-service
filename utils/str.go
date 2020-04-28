package utils

import (
	"math/rand"
	"strings"
	"time"
)

// 随机字符串 小写
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}

	return strings.ToLower(string(bytes))
}
