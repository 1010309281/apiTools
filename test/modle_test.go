package test

import (
	modles "apiTools/modle"
	"testing"
)

func init() {
	//config.InitConfig()
	//modles.InitRedis()
}

func TestWhoisQuery(t *testing.T) {
	form := &modles.WhoisForm{
		Domain:  "http://www.baidu.io",
		OutType: "json",
	}
	whoisInfo, err := modles.QueryWhoisInfoToJson(form)
	if err != nil {
		t.Error("err", err)
	}
	t.Logf("%v\n", whoisInfo.TextInfo)
	t.Logf("%#v\n", whoisInfo.JsonInfo)
}

func TestToShortUrl(t *testing.T) {
	shortForm := &modles.ShortForm{
		Url:        "https://www.runoob.com/python3/python-find-url-string.html",
		Domain:     "http://www.baidu.cn",
		ExpireTime: 7,
	}

	shortInfo, err := modles.ToShortUrl(shortForm)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(shortInfo)
}

func TestParseShortUrl(t *testing.T) {
	shortUrl := "http://www.baidu.com/DF5m1YsVSf"
	shortInfo, err := modles.ParseShort(shortUrl)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(shortInfo)
}

func TestInitJsonData(t *testing.T) {
	err := modles.InitJsonData()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
