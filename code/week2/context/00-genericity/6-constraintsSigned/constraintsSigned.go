package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
)

// MyInt int的派生类型
type MyInt int

func main() {
	a := MyInt(1)
	b := MyInt(2)
	value := minValue(a, b)
	fmt.Printf("minValue(%v, %v) = %v type = %v \n", a, b, value, reflect.TypeOf(value))
}

// minValue 返回两个值中较小的值 这个值的类型可能为 int/float32/float64
func minValue[Value constraints.Signed](a, b Value) Value {
	if a < b {
		return a
	}
	return b
}

func a(a any) {

}
