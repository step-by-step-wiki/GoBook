package main

import "fmt"

func main() {
	// 以string类型为key的map
	stringSet := NewSet[string, bool]()
	Add(stringSet, "apple", true)
	Add(stringSet, "orange", false)
	fmt.Printf("stringSet = %v\n", stringSet)

	// 以int类型为key的map
	intSet := NewSet[int, string]()
	Add(intSet, 1, "one")
	Add(intSet, 2, "two")
	fmt.Printf("intSet = %v\n", intSet)
}

// MySet 泛型map类型 要求其key必须是可比较的
type MySet[Key comparable, Value any] map[Key]Value

func (m MySet[Key, Value]) Add(key Key, value Value) {
	m[key] = value
}

// NewSet 创建一个泛型map
func NewSet[Key comparable, Value any]() MySet[Key, Value] {
	return make(MySet[Key, Value])
}

// Add 向泛型map中添加元素
func Add[Key comparable, Value any](set MySet[Key, Value], key Key, value Value) {
	set[key] = value
}
