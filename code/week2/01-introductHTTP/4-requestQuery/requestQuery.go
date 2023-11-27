package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/queryParams", queryParams)
	http.ListenAndServe(":8091", nil)
}

func queryParams(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fmt.Fprintf(w, "query is: %v\n", values)
}
