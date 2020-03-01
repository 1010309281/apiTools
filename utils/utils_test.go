package utils

import "testing"

func TestCheckProtocolHttp(t *testing.T) {
	bool := CheckProtocolHttp("218.26.178.22:8080", "www.baidu.com")
	if !bool {
		t.Error("fail")
		return
	}
	t.Log("success")
}