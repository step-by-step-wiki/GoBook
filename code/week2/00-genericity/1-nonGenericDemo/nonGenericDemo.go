package main

import "fmt"

func main() {
	fmt.Printf("minInt(42, 32) = %d\n", minInt(42, 32))
}

// minInt 返回2个整型中较小的一个
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
