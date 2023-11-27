package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5}
	handleFunc := func(needle int) int {
		return needle * 2
	}
	result := nonGenericityMap(s, handleFunc)
	fmt.Println(result)
}

// nonGenericityMap map()函数的非泛型实现
func nonGenericityMap(s []int, f func(int) int) (result []int) {
	result = make([]int, len(s))
	for index, element := range s {
		result[index] = f(element)
	}
	return result
}
