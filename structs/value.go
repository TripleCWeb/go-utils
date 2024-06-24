package structs

import (
	"reflect"
	"strconv"
	"unsafe"
)

// ValueCopy 对象拷贝
// target必须是CanSet的指针，如果是一个空指针，那么需要传递指针的指针。
// origin不限制类型，有值就行，指针/对象都可以
func ValueCopy(origin interface{}, target interface{}) {
	// 源数据获取到数据，非指针类型。
	originValue := reflect.ValueOf(origin)
	for {
		if reflect.TypeOf(originValue.Interface()).Kind() != reflect.Ptr {
			break
		}
		originValue = originValue.Elem()
	}

	// 源指针如果不能设置，则取地址的地址去设置
	var curTargetValue reflect.Value
	if target, ok := target.(reflect.Value); ok {
		curTargetValue = target
	} else {
		curTargetValue = reflect.ValueOf(target)
	}

	if !curTargetValue.CanSet() {
		panic("target value cannot be set!!!")
	}

	for {
		if curTargetValue.Kind() == reflect.Ptr {
			childType := curTargetValue.Type().Elem()
			childValue := reflect.New(childType)
			curTargetValue.Set(childValue)
			curTargetValue = curTargetValue.Elem()
		} else {
			targetKind := curTargetValue.Type().Kind()

			switch targetKind {
			case reflect.Int32:
				intValue, err := strconv.Atoi(originValue.Interface().(string))
				if err != nil {
					panic(err)
				}
				curTargetValue.SetInt(int64(intValue))
			case reflect.String:
				curTargetValue.SetString(originValue.Interface().(string))
				break
			}
			break
		}
	}
}

func Int64ToInt(n64 interface{}) int {
	n64Value := reflect.ValueOf(n64)
	if n64Value.Kind() == reflect.Ptr {
		n64Value = n64Value.Elem()
	}

	tn64 := n64Value.Interface().(int64)
	return *(*int)(unsafe.Pointer(&tn64))
}
