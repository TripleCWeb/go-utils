package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
)

func HttpGet(url string, responseData interface{}, response any) (err error) {
	// 发送请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	return nil
}

func HttpPost(url string, requestData interface{}, responseData interface{}) (err error) {
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := response.Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return fmt.Errorf("expected code = 0 OK, but got %d", response.Code)
	}

	// 解析data
	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid data format in response")
	}
	dataMapBytes, err := json.Marshal(dataMap)
	if err != nil {
		return fmt.Errorf("error decoding login data: %s", err.Error())
	}

	err = json.Unmarshal(dataMapBytes, &responseData)
	if err != nil {
		return fmt.Errorf("error decoding login data: %s", err.Error())
	}
	return
}
