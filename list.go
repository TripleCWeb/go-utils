package main

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m < v {
			m = v
		}
	}
	return m
}

func Min[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m > v {
			m = v
		}
	}
	return m
}

func Has(list interface{}, element interface{}) (isIn bool) {
	listValue := reflect.ValueOf(list)
	listKind := listValue.Kind()
	if listKind == reflect.Slice || listKind == reflect.Array {
		for i := 0; i < listValue.Len(); i++ {
			originValue := listValue.Index(i)
			if reflect.TypeOf(originValue).Kind() == reflect.Ptr {
				originValue = originValue.Elem()
			}
			targetValue := reflect.ValueOf(element)
			if reflect.TypeOf(element).Kind() == reflect.Ptr {
				targetValue = targetValue.Elem()
			}

			originObj := originValue.Interface()
			targetObj := targetValue.Interface()
			if originObj == targetObj {
				return true
			}
		}
	}
	return false
}
