package modles

import (
	"apiTools/libs/config"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

func init() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = InitRedis()
	if err != nil {
		panic(err)
	}
	err = InitMysql()
	if err != nil {
		panic(err)
	}
}

func TestWhoisQuery(t *testing.T) {
	err := InitWhoisServers()
	if err != nil {
		t.Error(err)
	}
	form := &WhoisForm{
		Domain:  "http://www.baidu.io",
		OutType: "json",
	}
	whoisInfo, err := QueryWhoisInfoToJson(form)
	if err != nil {
		t.Error("err", err)
	}
	t.Logf("%v\n", whoisInfo.TextInfo)
	t.Logf("%#v\n", whoisInfo.JsonInfo)
}

func TestToShortUrl(t *testing.T) {
	shortForm := &ShortForm{
		Url:        "https://www.runoob.com/python3/python-find-url-string.html",
		Domain:     "http://www.baidu.cn",
		ExpireTime: 7,
	}

	shortInfo, err := ToShortUrl(shortForm)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(shortInfo)
}

func TestParseShortUrl(t *testing.T) {
	shortUrl := "http://www.baidu.com/DF5m1YsVSf"
	shortInfo, err := ParseShort(shortUrl)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(shortInfo)
}

func TestInsertProxyInfo(t *testing.T) {
	proxyInfo := &ProxyPool{
		IP:         "1.1.1.2",
		Port:       "8089",
		Anonymity:  "高匿",
		Protocol:   "https",
		Speed:      sql.NullInt64{Int64: 1992, Valid: true},
		VerifyTime: time.Now(),
	}
	err := InsertProxyInfo(proxyInfo, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("insert success")
}

func TestCreateProxyInfo(t *testing.T) {
	for i := 1; i < 255; i++ {
		proxyInfo := &ProxyPool{
			IP:         fmt.Sprintf("%d:%d:%d:%d", i, i, i, i),
			Port:       fmt.Sprintf("8%d", i),
			Anonymity:  "透明",
			Protocol:   "http",
			VerifyTime: time.Now(),
		}
		err := InsertProxyInfo(proxyInfo, false)
		if err != nil {
			t.Error(err)
		}
		log.Printf("insert %s:%s success", proxyInfo.IP, proxyInfo.Port)
		time.Sleep(1 * time.Second)
	}
	t.Log("insert success")

}

func TestExtractProxyInfo(t *testing.T) {
	pools, err := ExtractProxyInfo(10)
	if err != nil {
		t.Error(err)
		return
	}
	for index, info := range pools {
		t.Logf("%d --> %#v\n", index, info)
	}
}

func TestGetLatestProxyInfo(t *testing.T) {
	proxyArray, err := GetLatestProxyInfo(10)
	if err != nil {
		t.Error(err)
		return
	}
	for index, info := range proxyArray {
		t.Logf("%d --> %s\n", index, info)
	}
}

func TestSetProxyInfoToRedis(t *testing.T) {
	proxyArray := []string{
		"254:254:254:254:8254",
		"253:253:253:253:8253",
		"252:252:252:252:8252",
		"251:251:251:251:8251",
		"250:250:250:250:8250",
		"249:249:249:249:8249",
		"248:248:248:248:8248",
		"247:247:247:247:8247",
		"246:246:246:246:8246",
		"245:245:245:245:8245",
	}
	err := SetProxyInfoToRedis("proxyPoolArray", proxyArray)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("set success")

}

func TestReadProxyInfoFromRedis(t *testing.T) {
	proxyArray, err := ReadProxyInfoFromRedis("proxyPoolArray")
	if err != nil {
		t.Error(err)
		return
	}
	for index, info := range proxyArray {
		t.Logf("%d --> %s\n", index, info)
	}
}

func TestDelOneProxyFromDB(t *testing.T) {
	ip := "183.166.133.146"
	err := DelOneProxyFromDB(ip)
	if err != nil {
		t.Error(err)
	}
	t.Log("del proxy ip success")
}

func TestQueryProxyPoolInfo(t *testing.T) {
	proxyPoolForm := &ProxyPoolForm{
		Page:      1,
		Country:   "中国",
		Protocol:  "https",
		Address:   "南京",
		OrderBy:   "speed",
		OrderRule: "desc",
	}
	proxyPools, err := QueryProxyPoolInfo(proxyPoolForm)
	if err != nil {
		t.Error(err)
		return
	}
	println("pages is: ", proxyPools.Pages)
	for index, info := range proxyPools.ProxyPools {
		fmt.Printf("%d --> %v\n", index, info)
	}
}
