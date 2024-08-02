package common

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron"
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

func FormatUse(use uint64) string {
	const (
		KB = 1 << (10 * (iota + 1))
		MB
		GB
		TB
	)

	switch {
	case use >= TB:
		return fmt.Sprintf("%.2f TB", float64(use)/TB)
	case use >= GB:
		return fmt.Sprintf("%.2f GB", float64(use)/GB)
	case use >= MB:
		return fmt.Sprintf("%.2f MB", float64(use)/MB)
	case use >= KB:
		return fmt.Sprintf("%.2f KB", float64(use)/KB)
	default:
		return fmt.Sprintf("%d B", use)
	}
}

func RunTransferred(value uint64, sourcePort string, destinationAddress string) {
	var m sync.Mutex
	var use uint64 = 1 //值为0时会重置使用量
	go Transferred(value, sourcePort, destinationAddress, func(_use uint64, _cur int) uint {
		//log.Printf("30003: %d, %s", _use, FormatUse(_use))
		if use == 0 {
			// log.Printf("reset use %d,%s", _use, FormatUse(_use))
			use = 1
			return 1
		} else {
			use = _use
		}
		return 100 //0停止 1重新统计
	})
	log.Printf("任务执行时间：%s", time.Now())
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		m.Lock()
		defer m.Unlock()
		if use > 0 {
			//log.Printf("任务执行时间：%s,%d,%s", time.Now(), use, FormatUse(use))
			//需要统计流量
		}
		use = 0 //重置状态
	})
	c.Start()
}
