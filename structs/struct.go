package structs

// import (
// 	"fmt"
// 	"reflect"
// 	"strings"

// 	"github.com/TripleCWeb/go-utils/"
// )

// const (
// 	JsonTag = "json"
// 	DdbTag  = "ddb"
// )

// type StructField struct {
// 	Name        string       // 名字
// 	Type        reflect.Type // 类型
// 	DdbTag      string       // ddb tag名字
// 	JsonTag     string       // json tag名字
// 	Value       interface{}  // 值
// 	StructValue reflect.Value
// 	StructField reflect.StructField // 完整的field信息
// }

// // 判断是否私有
// func IsStructFieldPrivate(sf *reflect.StructField) bool {
// 	firstChar := sf.Name[0]
// 	return firstChar >= 'a' && firstChar <= 'z'
// }

// type Struct struct {
// 	StructObj   interface{} // 可以是指针/对象
// 	Type        reflect.Type
// 	Fields      []*StructField
// 	json2Fields map[string]*StructField
// 	ddb2Fields  map[string]*StructField
// }

// func New(structObj interface{}) *Struct {
// 	s := &Struct{StructObj: structObj, json2Fields: make(map[string]*StructField), ddb2Fields: make(map[string]*StructField)}
// 	s.Parser()
// 	return s
// }

// func (p *Struct) Parser() {
// 	p.Fields = make([]*StructField, 0)

// 	p.Type = reflect.TypeOf(p.StructObj)
// 	v := reflect.ValueOf(p.StructObj)
// 	if p.Type.Kind() == reflect.Ptr {
// 		p.Type = p.Type.Elem()
// 		v = reflect.ValueOf(p.StructObj).Elem()
// 	}

// 	for i := 0; i < v.NumField(); i++ {
// 		typeField := p.Type.Field(i)

// 		// 不妨问私有字段，会出错
// 		if IsStructFieldPrivate(&typeField) {
// 			continue
// 		}

// 		name := typeField.Name
// 		structField, _ := p.Type.FieldByName(name)
// 		structValue := v.FieldByName(name)

// 		// json:"user_id,omitempty"
// 		jsonTag := structField.Tag.Get(JsonTag)
// 		if jsonTag != "" {
// 			jsonTag = strings.Split(jsonTag, ",")[0]
// 		}
// 		ddbTag := structField.Tag.Get(DdbTag)

// 		f := &StructField{
// 			Name:        name,
// 			Type:        structField.Type,
// 			DdbTag:      ddbTag,
// 			JsonTag:     jsonTag,
// 			Value:       structValue.Interface(),
// 			StructValue: structValue,
// 			StructField: structField,
// 		}
// 		p.Fields = append(p.Fields, f)
// 		if jsonTag != "" {
// 			p.json2Fields[jsonTag] = f
// 		}
// 		if ddbTag != "" {
// 			p.ddb2Fields[ddbTag] = f
// 		}

// 	}
// }

// func (p *Struct) GetFieldByJsonTag(json string) *StructField {
// 	if f, ok := p.json2Fields[json]; ok {
// 		return f
// 	}
// 	return nil
// }

// func (p *Struct) GetFieldByDdbTag(ddb string) *StructField {
// 	if f, ok := p.ddb2Fields[ddb]; ok {
// 		return f
// 	}
// 	return nil
// }

// // ForEachField每个字段的循环处理
// func (p *Struct) ForEachField(f func(f *StructField)) {
// 	for _, field := range p.Fields {
// 		f(field)
// 	}
// }

// // Serialize
// func (p *Struct) Serialize(sfunc func(valuePtr interface{}) string) string {
// 	var condLs []string
// 	New(p.StructObj).ForEachField(func(f *StructField) {
// 		if IsObjectEmpty(f.Value) {
// 			return
// 		}
// 		// 成员是结构体先不处理
// 		switch f.Type.Kind() {
// 		case reflect.Struct:
// 			return
// 		}
// 		condLs = append(condLs, fmt.Sprintf("%s=%s", f.JsonTag, sfunc(f.Value)))
// 	})
// 	return strings.Join(condLs, "&")
// }

// // CopyFrom 从目标struct拷贝数据, org:拷贝的参数对象
// func (p *Struct) CopyFrom(org interface{}) {
// 	orgStruct := New(org)
// 	if p.Type != orgStruct.Type {
// 		panic("type not the same")
// 	}

// 	orgStruct.ForEachField(func(orgField *StructField) {
// 		targetValue := reflect.ValueOf(p.StructObj).Elem().FieldByName(orgField.Name)
// 		if targetValue.IsValid() && !IsObjectEmpty(orgField.Value) {
// 			targetValue.Set(reflect.ValueOf(orgField.Value))
// 		}
// 	})
// }

// type FilterFunc func(obj interface{}) bool

// func (p *Struct) IsValid(filters ...FilterFunc) bool {
// 	if len(filters) == 0 {
// 		filters = append(filters, IsDefaultValid)
// 	}

// 	for _, field := range p.Fields {
// 		targetValue := reflect.ValueOf(p.StructObj).Elem().FieldByName(field.Name)
// 		if !targetValue.IsValid() {
// 			continue
// 		}
// 		for _, f := range filters {
// 			isValid := f(field.Value)
// 			if isValid {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// // CopyFields 从目标struct拷贝相同的字段
// // useTag 使用的tag，一个数据
// func (p *Struct) CopyFields(org interface{}, useTag string) {
// 	if !utils.Has([]string{DdbTag, JsonTag}, useTag) {
// 		panic("ilvalid tag")
// 	}

// 	orgStruct := New(org)
// 	p.ForEachField(func(f *StructField) {
// 		tag := ""
// 		if useTag == JsonTag {
// 			tag = f.JsonTag
// 		} else {
// 			tag = f.DdbTag
// 		}
// 		originFiled := orgStruct.GetFieldByJsonTag(tag)
// 		if originFiled == nil {
// 			return
// 		}
// 		if f.StructValue.IsValid() && !IsObjectEmpty(originFiled.Value) {
// 			originValue := originFiled.StructValue
// 			if originFiled.Type.Kind() == reflect.Ptr {
// 				originValue = originValue.Elem()
// 			}
// 			targetValue := f.StructValue
// 			if f.Type.Kind() == reflect.Ptr {
// 				// 先给字段设置指针
// 				targetValue.Set(reflect.New(f.Type.Elem()))
// 				// 再给指针设置值
// 				targetValue = targetValue.Elem()
// 			}
// 			targetValue.Set(originValue)
// 		}
// 	})
// }

// // CopyBYJsonFields 从目标struct拷贝相同的json字段
// func (p *Struct) CopyBYJsonFields(org interface{}) {
// 	p.CopyFields(org, JsonTag)
// }

// // CopyByDdbFields 从目标struct拷贝相同的ddb字段
// func (p *Struct) CopyByDdbFields(org interface{}) {
// 	p.CopyFields(org, JsonTag)
// }

// // IsEmpty 判断规则，
// // 返回True情况：空指针，空slice
// // 空字符串算有值，值为空
// // 其他无法判断，返回False
// func IsObjectEmpty(obj interface{}) bool {
// 	k := reflect.TypeOf(obj).Kind()
// 	switch k {
// 	case reflect.Ptr:
// 		// 如果数据为空，直接返回
// 		if reflect.ValueOf(obj).IsNil() {
// 			return true
// 		}
// 		// 判断指针对象的值
// 		return IsObjectEmpty(reflect.ValueOf(obj).Elem().Interface())
// 	case reflect.Slice:
// 		return reflect.ValueOf(obj).Len() == 0
// 	default:
// 		return false
// 	}
// }

// // IsValid 判断是否数据有效
// // 返回False情况：空指针，字符串长度为0
// // 空字符串算有值，值为空
// // 其他无法判断，返回False
// func IsDefaultValid(obj interface{}) bool {
// 	k := reflect.TypeOf(obj).Kind()
// 	switch k {
// 	case reflect.Ptr:
// 		// 如果数据为空，直接返回
// 		if reflect.ValueOf(obj).IsNil() {
// 			return false
// 		}
// 		// 判断指针对象的值
// 		return !IsObjectEmpty(reflect.ValueOf(obj).Elem().Interface())
// 	case reflect.Slice:
// 		return reflect.ValueOf(obj).Len() > 0
// 	case reflect.String:
// 		return reflect.ValueOf(obj).Len() > 0
// 	default:
// 		return true
// 	}
// }

// // CopyFields 拷贝字段
// // origin： 源结构体指针
// // target: 目标结构体指针
// // tag: 使用的字段
// func CopyFields(origin interface{}, target interface{}, tag string) {
// 	originType := reflect.TypeOf(origin)
// 	targetType := reflect.TypeOf(origin)
// 	if originType.Kind() != reflect.Ptr || targetType.Kind() != reflect.Ptr {
// 		panic("ilvalid type, must be pointer")
// 	}
// 	if originType.Elem().Kind() != reflect.Struct || targetType.Elem().Kind() != reflect.Struct {
// 		panic("ilvalid type, must be struct")
// 	}
// 	New(target).CopyFields(origin, tag)
// }
