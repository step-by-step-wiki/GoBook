package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/formAndPostForm", formAndPostFormHandle)
	log.Fatal(http.ListenAndServe(":8091", nil))
}

func formAndPostFormHandle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("ParseForm failed: ", err)
	}

	fmt.Fprintf(w, "Form: %v\n", r.Form)
	fmt.Fprintf(w, "PostForm: %v\n", r.PostForm)
}
