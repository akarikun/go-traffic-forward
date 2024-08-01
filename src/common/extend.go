package common

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

func IsNullOrEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

var _uuid = uuid.New()

func UUID() string {
	return _uuid.String()
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func Random(length int) string {
	// 定义可用的字符
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 创建一个字符串构建器
	var b strings.Builder

	// 随机生成指定长度的字符串
	for i := 0; i < length; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}

	return b.String()
}
