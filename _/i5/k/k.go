package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ktye/k"
)

func main() {
	k := k.New()
	for _, s := range os.Args[1:] {
		if strings.HasSuffix(s, ".k") {
			load(k, s)
		} else {
			panic("usage: k [file.k..]")
		}
	}
	repl(k)
}

func repl(k *k.K) {
	s := bufio.NewScanner(os.Stdin)
	fmt.Printf("ktye/k\n ")
	for s.Scan() {
		t := s.Text()
		switch t {
		case `\\`:
			os.Exit(0)
		case `\`, `\h`:
			fmt.Println("...")
		case `\v`:
			vars(k)
		case `\s`:
			fmt.Println(string(k.Gostack()))
		case `\d`:
			//
		default:
			var tic time.Time
			if strings.HasPrefix(t, `\t`) {
				t = t[2:]
				tic = time.Now()
			}
			line(k, t)
			if tic.IsZero() == false {
				fmt.Println(time.Since(tic))
			}
		}
		fmt.Printf(" ")
	}
}

func line(k *k.K, s string) {
	r := k.Run([]byte(s))
	if r != "" {
		fmt.Println(r)
	}
}
func load(k *k.K, file string) {
	k.Trap = true
	defer func() { k.Trap = false }()
	b, e := os.ReadFile(file)
	if e != nil {
		panic(e)
	}
	if s := k.Run(b); s != "" {
		fmt.Println(s)
	}
}
func vars(k *k.K) {
	name, typ, n, rc := k.Vars()
	w := tabwriter.NewWriter(os.Stdout, 2, 8, 1, ' ', 0)
	fmt.Fprintf(w, "name\ttype\tlen\trc\n")
	for i, s := range name {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", s, string(typ[i]), n[i], rc[i])
	}
	w.Flush()
}
