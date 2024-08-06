package common

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"regexp"
	"strconv"
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

// 获取并验证端口,暂只开放50000-60000之间的端口
//
// sourcePort 匹配格式: 50001 / :50001 / localhost:40000 / 127.0.0.1:40000 / aaaaa.com:40000 / aa.bb.com:40000
//
// uint16 返回匹配端口
// string 返回"IP:端口"格式,无IP则为127.0.0.1
func GetPort(sourcePort string) (uint16, string, error) {
	regex := regexp.MustCompile(`^(((?:(?:\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|localhost|(?:[a-z0-9]+\.)+[a-z0-9]+)?)?:?)(\d{5}))$`)
	match := regex.FindStringSubmatch(sourcePort)
	if match != nil {
		ip := match[2]
		port := match[3]

		p, err := strconv.Atoi(port)
		if err != nil {
			return 0, "", errors.New("配置异常:" + sourcePort)
		}
		if p <= 50000 || p >= 60000 {
			return 0, "", errors.New("暂只开放50000-60000之间的端口," + sourcePort)
		}

		if ip == ":" {
			return uint16(p), fmt.Sprintf("127.0.0.1:%s", port), nil
		}
		return uint16(p), fmt.Sprintf("%s:%s", ip, port), nil
	} else {
		return 0, "", errors.New("配置异常:" + sourcePort)
	}
}

func ValidatePort(cp string) (bool, error) {
	_, _, err := GetPort(cp)
	if err != nil {
		return false, err
	}
	conn, err := net.DialTimeout("tcp", cp, time.Second)
	if err != nil {
		return true, nil
	} else {
		defer conn.Close()
		return false, nil
	}
}

type TransferredUpdateFunc func(uint64)

func RunTransferred(value uint64, minute int, sourcePort string, destinationAddress string, action TransferredUpdateFunc) {
	var m sync.Mutex
	var use uint64 = 1 //值为0时会重置使用量
	if err := tcp_transferred(value, sourcePort, destinationAddress, func(tm *TransModel) {
		tm.initFunc = func() {
			log.Printf("%s任务执行时间：%s", sourcePort, time.Now())
			c := cron.New()
			c.AddFunc(fmt.Sprintf("@every %dm", minute), func() {
				m.Lock()
				defer m.Unlock()
				if use > 0 {
					log.Printf("任务执行时间：%s,%d,%s", time.Now(), use, FormatUse(use))
					//需要统计流量
					action(use)
				}
				use = 0 //重置状态
			})
			c.Start()
		}
		tm.transfunc = func(_cur int) uint {
			// log.Printf("reset use %d,%d,%s", use, tm.use, FormatUse(tm.use))
			if use == 0 {
				// log.Printf("reset use %d,%s", tm.use, FormatUse(tm.use))
				use = 1
				return 1
			} else {
				use = tm.use
			}
			return 100 //0停止 1重新统计
		}
	}); err != nil {
		fmt.Printf("tcp_transferred error:%s,%s\r\n", sourcePort, destinationAddress)
		return
	}
}
