package main

import (
	"math/rand"
	"time"
)

func RandomHexStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	hexChars := "0123456789abcdef"
	hexString := "0x"
	for i := 0; i < length-2; i++ {
		hexString += string(hexChars[rand.Intn(len(hexChars))])
	}
	return hexString
}
