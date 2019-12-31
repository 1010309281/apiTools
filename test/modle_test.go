package test

import (
	"apiTools/libs/config"
	modles "apiTools/modle"
	"testing"
)

func init() {
	config.InitConfig()
	modles.InitRedis()
}

func TestWhoisQuery(T *testing.T) {
	form := &modles.WhoisForm{
		Domain:  "http://www.baidu.io",
		OutType: "json",
	}
	whoisInfo, err := modles.QueryWhoisInfoToJson(form)
	if err != nil {
		T.Error("err", err)
	}
	T.Logf("%v\n", whoisInfo.TextInfo)
	T.Logf("%#v\n", whoisInfo.JsonInfo)
}

func TestToShortUrl(T *testing.T) {
	shortForm := &modles.ShortForm{
		Url:        "https://www.runoob.com/python3/python-find-url-string.html",
		Domain:     "http://www.baidu.cn",
		ExpireTime: 7,
	}

	shortInfo, err := modles.ToShortUrl(shortForm)
	if err != nil {
		T.Error(err)
		return
	}
	T.Log(shortInfo)
}

func TestParseShortUrl(T *testing.T) {
	shortUrl := "http://www.baidu.com/DF5m1YsVSf"
	shortInfo, err := modles.ParseShort(shortUrl)
	if err != nil {
		T.Error(err)
		return
	}
	T.Log(shortInfo)
}
