package main

import (
	"encoding/json"
	"log"
	"strings"
)

var jsonStr = `{"knownField":"value", "unknownField":"should cause error"}`

type MyStruct struct {
	KnownField string `json:"knownField"`
}

func main() {
	noDisallowUnknownFields()
	disallowUnknownFields()
}

// noDisallowUnknownFields 不使用`DisallowUnknownFields`选项进行JSON反序列化
func noDisallowUnknownFields() {
	myStruct := &MyStruct{}
	err := json.Unmarshal([]byte(jsonStr), myStruct)
	if err != nil {
		log.Fatal("Unmarshal error:", err)
	}

	log.Printf("Unmarshal success: %+v\n", myStruct)
}

// disallowUnknownFields 使用`DisallowUnknownFields`选项进行JSON反序列化
func disallowUnknownFields() {
	myStruct := &MyStruct{}
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(myStruct)
	if err != nil {
		// 这里将输出错误 因为JSON中包含了MyStruct没有定义的unknownField
		log.Fatal("Decode error:", err)
	}

	log.Printf("Decode success: %+v\n", myStruct)
}
