package main

import (
	"os"
)

// k='rlwrap -H /dev/null /path/to/k.exe -s `stty size`'
func main() {
	version = "k/g"
	Main(os.Args[1:])
}
