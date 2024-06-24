package utils

import (
	"math/rand"
	"time"
)

const (
	// 定义字符串的字符集
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	// 生成随机字符串
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomString)
}

func RandomTransactionHash() string {
	return RandomHexStr(128)
}

func RandomSign() string {
	return RandomHexStr(132)
}

func RandomHexStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	hexChars := "0123456789abcdef"
	hexString := "0x"
	for i := 0; i < length-2; i++ {
		hexString += string(hexChars[rand.Intn(len(hexChars))])
	}
	return hexString
}
