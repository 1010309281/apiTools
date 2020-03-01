package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// 获取项目根目录
func GetRootPath() (rootPath string) {
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("get project root path faild: %s", err))
	}
	return
}

// 字符串计算md5值
func GetMD5(text string) string {
	h := md5.New()
	salt := "apiTools"
	io.WriteString(h, text+salt)
	urlmd5 := fmt.Sprintf("%x", h.Sum(nil))
	return urlmd5
}

// 获取随机的唯一短串
func GetShortStr() (tiny string) {
	// 时间戳随机加盐避免重复
	rand.Seed(time.Now().UnixNano() >> 3)
	num := rand.Int63n(time.Now().UnixNano() >> 3)
	alpha := merge(getRange(48, 57), getRange(65, 90))
	alpha = merge(alpha, getRange(97, 122))
	if num < 62 {
		tiny = string(alpha[num])
		return tiny
	} else {
		var runes []rune
		runes = append(runes, alpha[num%62])
		num = num / 62
		for num >= 1 {
			if num < 62 {
				runes = append(runes, alpha[num-1])
			} else {
				runes = append(runes, alpha[num%62])
			}
			num = num / 62
		}
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		tiny = string(runes)
		return
	}
}

func getRange(start, end rune) (ran []rune) {
	for i := start; i <= end; i++ {
		ran = append(ran, i)
	}
	return ran
}

func merge(a, b []rune) []rune {
	c := make([]rune, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}

// 重新定义cron定时任务初始化
func NewWithCron() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

// 定时器，用于堵塞进程
// second 定时时间 秒
func TimerUtil(second int64) {
	timer := time.NewTimer(time.Second * time.Duration(second))
	<-timer.C
}

// 判断值是否在一个切片中存在
func IsInSelic(data string, slice []string) bool {
	for _, s := range slice {
		if s == data {
			return true
		}
	}
	return false
}
