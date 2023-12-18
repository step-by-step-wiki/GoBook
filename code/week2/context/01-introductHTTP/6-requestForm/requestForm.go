package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/form", form)
	http.ListenAndServe(":8091", nil)
}

func form(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Before ParseForm: %v\n", r.Form)
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm err: %v\n", err)
		return
	}
	fmt.Fprintf(w, "After ParseForm: %v\n", r.Form)
}
