package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/header", header)
	http.ListenAndServe(":8091", nil)
}

func header(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "header is: %v\n", r.Header)
}
