package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
)

func main() {
	value := minValue(-1, 2.4)
	fmt.Printf("minValue(%v, %v) = %v type = %v \n", -1, 2.4, value, reflect.TypeOf(value))
}

// minValue 返回两个值中较小的值 这个值的类型可能为 int/float32/float64
func minValue[Value constraints.Ordered](a, b Value) Value {
	if a < b {
		return a
	}
	return b
}
