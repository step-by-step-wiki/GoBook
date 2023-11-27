package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	f := func(needle int) bool {
		return needle%2 == 0
	}
	result := genericityFilter(s, f)
	fmt.Printf("result: %v\n", result)
}

func genericityFilter(s []int, f func(int) bool) []int {
	var result []int
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
