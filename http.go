package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// httpRequest 类是用于发送 HTTP 请求的类
type httpRequest struct {
}

func NewHttpRequest() *httpRequest {
	r := new(httpRequest)
	return r
}

func (p *httpRequest) Get(url string, response interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}

	return nil
}

func (p *httpRequest) Request(url string, method string, request interface{}, response interface{}) (err error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("转换JSON失败:", err)
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Close = true
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求失败:", resp.StatusCode)
		err = fmt.Errorf(resp.Status)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 解析JSON响应到结构体
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return
	}
	return
}

func (p *httpRequest) POST(url string, request interface{}, response interface{}) (err error) {
	return p.Request(url, "POST", request, response)
}

func (p *httpRequest) PUT(url string, request interface{}, response interface{}) (err error) {
	return p.Request(url, "PUT", request, response)
}
