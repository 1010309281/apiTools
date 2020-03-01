package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// http proxy get方法
func HttpProxyGet(dataUrl, proxyIp string, headers map[string]string) (data []byte, response *http.Response, err error) {
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //ssl证书报错问题
		DisableKeepAlives: false,                                 //关闭连接复用，因为后台连接过多最后会造成端口耗尽
		MaxIdleConns:      -1,                                    //最大空闲连接数量
		IdleConnTimeout:   time.Duration(5 * time.Second),        //空闲连接超时时间
	}
	if proxyIp != "" { // 设置代理
		proxyUrl, _ := url.Parse("http://" + proxyIp)
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	client := &http.Client{
		Timeout:   time.Duration(30 * time.Second),
		Transport: transport,
	}

	request, err := http.NewRequest("GET", dataUrl, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")
	if headers != nil {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}
	if err != nil {
		return
	}
	// 请求数据
	resp, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("request %s, proxyIp: (%s),err: %v", dataUrl, proxyIp, err)
		return
	}
	defer resp.Body.Close()
	// 读取数据
	buf := make([]byte, 128)
	data = make([]byte, 0, 2048)
	for {
		n, err := resp.Body.Read(buf)
		data = append(data, buf[:n]...)

		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
	}
	response = resp
	return
}

// 校验代理 http和http协议协议代理
func CheckProtocolHttp(proxyAddr, checkToUrl string) bool {
	httpClient := &http.Client{
		Timeout: time.Duration(10 * time.Second), //客户端设置10秒超时
	}
	httpClient.Transport = &http.Transport{
		DisableKeepAlives: false,                          //关闭连接复用，因为后台连接过多最后会造成端口耗尽
		MaxIdleConns:      -1,                             //最大空闲连接数量
		IdleConnTimeout:   time.Duration(5 * time.Second), //空闲连接超时时间
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   proxyAddr,
		}),                                                //设置http代理地址
	}
	_, err := httpClient.Get(fmt.Sprintf("http://%s", checkToUrl))
	if err != nil {
		return false
	}
	return true
}

