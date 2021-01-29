package main

import (
	"k"
	"fmt"
	"strings"
	"os"
)

func main() {
	i := k.New()
	fmt.Println(i.Run(k.Parse([]byte(strings.Join(os.Args[1:], " ")))))
}
