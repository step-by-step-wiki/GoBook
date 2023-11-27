package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/queryValue", queryValue)
	http.ListenAndServe(":8091", nil)
}

func queryValue(w http.ResponseWriter, r *http.Request) {
	// 获取查询字符串中的name参数
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "name %s", name)
}
