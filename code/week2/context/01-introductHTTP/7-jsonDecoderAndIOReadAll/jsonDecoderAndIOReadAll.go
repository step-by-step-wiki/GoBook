package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	Id int `json:"id"`
}

func main() {
	http.HandleFunc("/unmarshal", unmarshalHandle)
	http.HandleFunc("/decoder", decoderHandle)
	http.ListenAndServe(":8091", nil)
}

func unmarshalHandle(w http.ResponseWriter, r *http.Request) {
	byteSlice, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(byteSlice, &User{})
	if err != nil {
		fmt.Fprintf(w, "decode failed: %v", err)
		return
	}

	afterRead, _ := io.ReadAll(r.Body)
	fmt.Fprintf(w, "after read: %s", string(afterRead))
}

func decoderHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&User{})
	if err != nil {
		fmt.Fprintf(w, "decode failed: %v", err)
		return
	}

	afterRead, _ := io.ReadAll(r.Body)
	fmt.Fprintf(w, "after read: %s", string(afterRead))
}
