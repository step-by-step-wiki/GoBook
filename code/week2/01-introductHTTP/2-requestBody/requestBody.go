package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/readBodyOnce", readBodyOnce)
	http.ListenAndServe(":8091", nil)
}

func readBodyOnce(w http.ResponseWriter, r *http.Request) {
	// 先读取一次body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		return
	}
	fmt.Fprintf(w, "read body: %s\n", string(body))

	// 再读取一次body
	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body one more time got error: %v", err)
		return
	}
	fmt.Fprintf(w, "read body one more time: [%s] and read body length %d \n", string(body), len(body))
}
