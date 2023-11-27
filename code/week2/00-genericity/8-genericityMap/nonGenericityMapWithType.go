package main

import "fmt"

func main() {
	mySlice := MySlice{1, 2, 3, 4, 5}
	handleFunc := func(needle int) int {
		return needle * 2
	}
	result := mySlice.Map(handleFunc)
	fmt.Println(result)
}

type MySlice []int

func (m MySlice) Map(f func(int) int) (result MySlice) {
	result = make(MySlice, len(m))
	for index, element := range m {
		result[index] = f(element)
	}
	return result
}
