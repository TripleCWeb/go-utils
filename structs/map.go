package structs

import "encoding/json"

// Map converts the given struct to a map[string]interface{}. For more info
// refer to Struct types Map() method. It panics if s's kind is not struct.
func Map(s interface{}) map[string]interface{} {
	// 如果struct中的json tag包含omitempty，例如
	// type Student struct {
	// Id 		int 	`json:"id"`
	// Name 	string 	`json:"name,omitempty"`
	// }
	// 那么json序列化的时候将不输出该字段

	b, _ := json.Marshal(s)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}
