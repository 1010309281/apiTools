package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	apiUrl := "https://api.devopsclub.cn/api/whoisquery"
	//apiUrl := "http://127.0.0.1:8081/api/whoisquery"
	// 组装请求数据
	reqParam := make(map[string]interface{})
	reqParam["domain"] = "devopsclub.cn"
	reqParam["type"] = "json"
	reqBytes, _ := json.Marshal(reqParam)
	reader := bytes.NewReader([]byte(reqBytes))

	// 发送post请求
	reponse, err := http.Post(apiUrl, "application/json", reader)
	if err != nil {
		fmt.Printf("接口请求失败, err: %v\n", err)
		return
	}
	var content []byte
	for {
		buf := make([]byte, 1024)
		n, err := reponse.Body.Read(buf)
		if err == io.EOF || err != nil {
			content = append(content, buf[:n]...)
			break
		} else {
			content = append(content, buf[:n]...)
		}
	}
	fmt.Println(string(content))
}
