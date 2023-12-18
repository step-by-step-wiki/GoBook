package main

import "fmt"

func main() {
	gIntSlice := GSlice[int]{1, 2, 3, 4, 5}

	reduceIntFunc := func(previous int, current int) (result int) {
		return previous + current
	}

	mapIntFunc := func(needle int) int {
		return needle * 2
	}

	resultInt := gIntSlice.Map(mapIntFunc).Reduce(reduceIntFunc)
	fmt.Println(resultInt)
}

type GSlice[T any] []T

func (g GSlice[T]) Map(f func(T) T) (result GSlice[T]) {
	result = make(GSlice[T], len(g))
	for index, element := range g {
		result[index] = f(element)
	}
	return result
}

func (g GSlice[T]) Reduce(f func(previous, current T) (result T)) (result T) {
	for _, element := range g {
		result = f(result, element)
	}
	return result
}
