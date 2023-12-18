package main

import (
	"fmt"
	"reflect"
)

// Number 自定义的类型约束 类型可以为int/float32/float64
type Number interface {
	int | float32 | float64
}

func main() {
	value := minValue(-1, 2.4)
	fmt.Printf("minValue(%v, %v) = %v type = %v \n", -1, 2.4, value, reflect.TypeOf(value))
}

// minValue 返回两个值中较小的值 这个值的类型可能为 int/float32/float64
func minValue[Value Number](a, b Value) Value {
	if a < b {
		return a
	}
	return b
}
