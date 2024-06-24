package main

import (
	"encoding/json"
	"reflect"
)

func Interface2String(i interface{}) string {
	if reflect.ValueOf(i).Kind() == reflect.String {
		return i.(string)
	}

	jsonData, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

func String2Interface(jStr string, i interface{}) {
	Byte2Interface([]byte(jStr), i)
}

func Byte2Interface(jByte []byte, i interface{}) {
	err := json.Unmarshal(jByte, i)
	if err != nil {
		panic(err)
	}
}
