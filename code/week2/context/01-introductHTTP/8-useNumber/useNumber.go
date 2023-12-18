package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
)

var (
	jsonBlob = []byte(`{"int_max":9223372036854775807}`)
)

func main() {
	noUseNumber()
	useNumber()
}

// noUseNumber 不使用`UseNumber`选项进行JSON反序列化
func noUseNumber() {
	data := make(map[string]any)

	if err := json.Unmarshal(jsonBlob, &data); err != nil {
		log.Fatal(err)
	}

	// 输出一个浮点数 这个浮点数可能会出现精度丢失的现象
	fmt.Println(data["int_max"])
}

// useNumber 使用`UseNumber`选项进行JSON反序列化
func useNumber() {
	data := make(map[string]json.Number)

	decoder := json.NewDecoder(bytes.NewReader(jsonBlob))
	decoder.UseNumber()

	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
	}

	fmt.Println(data["int_max"])

	// 将这个json.Number类型的值安全的转换为更精确的数字类型
	intValue, ok := new(big.Int).SetString(data["int_max"].String(), 10)
	if !ok {
		log.Fatal("Big int conversion failed")
	}

	// 输出一个精确的大整数
	fmt.Printf("The big int is: %d\n", intValue)
}
