package rcache

import (
	"github.com/TripleCWeb/go-utils"
)

func Interface2String(i interface{}) string {
	return utils.Interface2String(i)
}

func String2Interface(jStr string, i interface{}) {
	utils.String2Interface(jStr, i)
}

// FieldSerializer 类型
type FieldSerializer struct{}

// 序列化 Field
func (fs FieldSerializer) Serialize(i interface{}) string {
	return Interface2String(i)
}

// 反序列化 Field
func (fs FieldSerializer) Deserialize(data string, i interface{}) error {
	String2Interface(data, i)
	return nil
}

// ValueSerializer 类型
type ValueSerializer struct{}

// 序列化 Value
func (fs ValueSerializer) Serialize(i interface{}) string {
	return Interface2String(i)
}

// 反序列化 Value
func (fs ValueSerializer) Deserialize(data string, i interface{}) error {
	String2Interface(data, i)
	return nil
}
