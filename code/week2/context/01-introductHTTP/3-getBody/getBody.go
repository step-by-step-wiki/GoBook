package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/getBodyIsNil", getBodyIsNil)
	http.ListenAndServe(":8091", nil)
}

func getBodyIsNil(w http.ResponseWriter, r *http.Request) {
	if r.GetBody == nil {
		fmt.Fprintf(w, "GetBody is nil\n")
	} else {
		fmt.Fprintf(w, "GetBody is not nil\n")
	}
}
