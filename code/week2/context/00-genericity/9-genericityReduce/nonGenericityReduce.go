package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5}
	handleFunc := func(previous, needle int) (result int) {
		return previous + needle
	}
	result := nonGenericityReduce(s, handleFunc)
	fmt.Println(result)
}

// nonGenericityReduce reduce()函数的非泛型实现
func nonGenericityReduce(s []int, f func(int, int) int) int {
	var result int
	for _, element := range s {
		result = f(result, element)
	}
	return result
}
