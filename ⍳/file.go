package main

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/ktye/i"
)

func file(name string, b []byte, a map[v]v) {
	file := name
	defer func() {
		if c := recover(); c != nil {
			stk, err := stack(c)
			println(stk)
			println(file + ":" + err)
		}
	}()
	lines := bytes.Split(b, []byte("\n"))
	for n := range lines {
		file = name + ":" + strconv.Itoa(n+1)
		s := string(lines[n])
		if strings.TrimSpace(s) == "" {
			continue
		}
		i.E(i.P(s), a)
	}
}
