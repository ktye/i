//+build ignore

package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println(http.ListenAndServe(":3000", nil))
}
