package modles

import (
	"apiTools/utils"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jonsen/gotld"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

// 接收转换为短链接 api form表单
type ShortForm struct {
	Url        string `form:"url" json:"url" xml:"url" binding:"required"`   // 要进行转换的长链接
	Domain     string `form:"domain" json:"domain" xml:"domain"`             // 短链接域名绑定自己的域名
	ExpireTime int    `form:"expireTime" json:"expireTime" xml:"expireTime"` // 设置过期时间, (以分钟为单位, -1代表用不过期)
}

// 返回的生成短链接信息数据结构
type ShortInfo struct {
	LongUrlMd5 string // 长链接转换后的md5值
	LongUrl    string // 原来长链接
	Domain     string // 短链接域名绑定的域名
	ShortStr   string // 短链接串
}

// 检测转换的域名是否存在，如果存在则返回数据
func checkShortUrl(shortForm *ShortForm) (*ShortInfo, bool, error) {
	shortInfo := &ShortInfo{
		LongUrlMd5: utils.GetMD5(shortForm.Url + shortForm.Domain),
	}
	// 获取redis连接
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	exists, err := redis.Bool(redisClient.Do("EXISTS",
		fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5)))
	if err != nil || !exists {
		_, _ = redisClient.Do("DEL", fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5))
		return shortInfo, false, fmt.Errorf("check long exist, err: %v", err)
	}
	// 根据长串获取存储的短域名信息
	LongInfoBytes, err := redis.ByteSlices(redisClient.Do("HMGET",
		fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5),
		"longUrl", "domain", "shortStr"))
	if err != nil || len(LongInfoBytes) != 3 {
		return shortInfo, false, fmt.Errorf("get long info fail, err: %v", err)
	}
	// 赋值
	shortInfo.LongUrl = string(LongInfoBytes[0])
	shortInfo.Domain = string(LongInfoBytes[1])
	shortInfo.ShortStr = string(LongInfoBytes[2])

	return shortInfo, true, nil
}

// 转换成短链接
func ToShortUrl(shortForm *ShortForm) (*ShortInfo, error) {
	// 校验长连接是否存在
	tmpShortInfo, status, err := checkShortUrl(shortForm)
	if status && err == nil {
		return tmpShortInfo, nil
	}
	// 创建一个shortInfo数据对象
	shortInfo := &ShortInfo{
		LongUrl:    shortForm.Url,
		LongUrlMd5: utils.GetMD5(shortForm.Url + shortForm.Domain),
		ShortStr:   utils.GetShortStr(),
	}
	redisClient := RedisPool.Get()
	defer redisClient.Close()

	// domain
	_, _, err = gotld.GetTld(shortForm.Domain)
	if err != nil {
		return shortInfo, fmt.Errorf("short form domain [%s] not normal domain name", shortForm.Domain)
	}
	if ok := strings.HasPrefix(shortForm.Domain, "http"); !ok {
		shortForm.Domain = fmt.Sprintf("http://%s", shortForm.Domain)
	}
	parse, err := url.Parse(shortForm.Domain)
	if err != nil {
		return shortInfo, fmt.Errorf("short domain [%s] name parse fail", shortForm.Domain)
	}
	shortInfo.Domain = parse.Hostname()
	// 事务开始创建数据到redis
	_ = redisClient.Send("MULTI")
	// 设置hash   长地址md5值对应相关数据
	_ = redisClient.Send(
		"HMSET",
		fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5),
		"longUrl", shortInfo.LongUrl,
		"domain", shortInfo.Domain,
		"shortStr", shortInfo.ShortStr,
	)
	// 设置string 短串对应长地址md5值
	_ = redisClient.Send("SET", fmt.Sprintf("short_%s", shortInfo.ShortStr),
		shortInfo.LongUrlMd5)
	// 设置过期时间, 默认为用不过期
	if shortForm.ExpireTime != -1 {
		expireTime := 60 * shortForm.ExpireTime

		// 设置long md5过期
		_ = redisClient.Send("EXPIRE", fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5), expireTime)
		// 设置short str过期
		_ = redisClient.Send("EXPIRE", fmt.Sprintf("short_%s", shortInfo.ShortStr), expireTime)
	}
	_, err = redisClient.Do("EXEC")
	if err != nil {
		_, _ = redisClient.Do("DISCARD")
		return shortInfo, fmt.Errorf("set short data to redis fail, err: %v", err)
	}
	return shortInfo, nil
}

// 解析短链接
func ParseShort(shortUrl string) (*ShortInfo, error) {
	// 创建一个数据返回对象
	shortInfo := &ShortInfo{}
	// 解析短地址url
	urlParse, err := url.Parse(shortUrl)
	if err != nil {
		return shortInfo, fmt.Errorf("parse short url fail, err: %v", err)
	}
	// 短链接域名
	shortInfo.Domain = urlParse.Hostname()
	// 短链接短串
	if ok := strings.HasPrefix(urlParse.Path, "/"); !ok {
		return shortInfo, errors.New("short url not has param")
	}
	shortPathSlice := strings.Split(urlParse.Path, "/")
	if len(shortPathSlice) > 3 || len(shortPathSlice) < 2 {
		return shortInfo, errors.New("short url param malformed")
	}
	shortStr := shortPathSlice[1]
	if shortUrl == "" {
		return shortInfo, errors.New("short url param malformed")
	}
	shortInfo.ShortStr = shortStr
	// 获取redis连接
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	// 根据短串获取长串md5
	longUrlMd5, err := redis.String(redisClient.Do("GET",
		fmt.Sprintf("short_%s", shortInfo.ShortStr)))
	if err != nil {
		return shortInfo, fmt.Errorf("get long url md5 fail from short str, err: %v", err)
	}
	if longUrlMd5 == "" {
		_, _ = redisClient.Do("DEL", fmt.Sprintf("short_%s", shortInfo.ShortStr))
		return shortInfo, errors.New("get long url md5 fail, because is empty")
	}
	shortInfo.LongUrlMd5 = longUrlMd5

	// 根据长串获取存储的短域名信息
	LongInfoBytes, err := redis.ByteSlices(redisClient.Do("HMGET", fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5),
		"longUrl", "domain", "shortStr"))
	if err != nil || len(LongInfoBytes) != 3 {
		return shortInfo, fmt.Errorf("get long info fail, err: %v", err)
	}
	// 赋值
	shortInfo.LongUrl = string(LongInfoBytes[0])
	shortInfo.Domain = string(LongInfoBytes[1])
	shortInfo.ShortStr = string(LongInfoBytes[2])

	return shortInfo, nil
}
