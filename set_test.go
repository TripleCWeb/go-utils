package main

import (
	"testing"
)

func TestSet(t *testing.T) {
	// 创建一个新的 Set
	s := NewSet[string]()
	if s.Size() != 0 {
		t.Error("New Set should be empty")
	}

	// 向 Set 中添加元素
	s.Add("foo")
	s.Add("bar")
	if s.Size() != 2 {
		t.Error("Set size should be 2")
	}

	// 尝试重复添加元素
	s.Add("foo")
	if s.Size() != 2 {
		t.Error("Set size should still be 2")
	}

	// 从 Set 中删除元素
	s.Remove("foo")
	if s.Size() != 1 {
		t.Error("Set size should be 1 after removing an element")
	}

	// 尝试删除不存在的元素
	s.Remove("baz")
	if s.Size() != 1 {
		t.Error("Set size should still be 1 after removing a non-existent element")
	}

	// 检查元素是否存在于 Set 中
	if !s.Contains("bar") {
		t.Error("Set should contain 'bar'")
	}
	if s.Contains("foo") {
		t.Error("Set should not contain 'foo' after removing it")
	}

	// 获取 Set 中所有元素的键
	keys := s.Keys()
	if len(keys) != 1 {
		t.Error("Set should have 1 key after removing an element")
	}
	if keys[0] != "bar" {
		t.Error("Set should have key 'bar'")
	}
}
