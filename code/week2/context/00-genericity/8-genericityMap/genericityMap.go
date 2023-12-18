package main

import "fmt"

func main() {
	gIntSlice := GSlice[int]{1, 2, 3, 4, 5}
	handleIntFunc := func(needle int) int {
		return needle * 2
	}
	resultInt := gIntSlice.Map(handleIntFunc)
	fmt.Println(resultInt)

	gStringSlice := GSlice[string]{"1", "2", "3", "4", "5"}
	handleStringFunc := func(needle string) string {
		return needle + " roach"
	}
	fmt.Println(gStringSlice.Map(handleStringFunc))
}

type GSlice[T any] []T

func (g GSlice[T]) Map(f func(T) T) (result GSlice[T]) {
	result = make(GSlice[T], len(g))
	for index, element := range g {
		result[index] = f(element)
	}
	return result
}
