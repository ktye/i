package main

import (
	"fmt"
	"os"
)

func c(s string, a map[v]v) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] != '/' && s[0] != '\\' {
		return true
	}
	switch s[1:] {
	case "":
		fmt.Print(a["doc"].(string))
	case `\`, "/":
		os.Exit(0)
	}
	fmt.Printf("\n ")
	return false
}
