//+build ignore

package main

import (
        "net/http"
	"fmt"
)

func main() {
        http.Handle("/", http.FileServer(http.Dir(".")))
        fmt.Println(http.ListenAndServe(":3000", nil))
}
