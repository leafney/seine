package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestIconURLToBase64(t *testing.T) {
	// 示例用法，设置超时为 10 秒
	url := "https://example.com/icon.png"
	timeout := 10 * time.Second
	base64Str, err := IconURLToBase64(url, timeout)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}
	fmt.Println("Base64 编码:", base64Str)
}
