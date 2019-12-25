// +/!
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type I = int64

func main() {
	ns := parse(os.Args[1])           // chunck size per go-routine
	np := parse(os.Args[2])           // number of go-routines running in parallel
	parallel := parseBool(os.Args[3]) // true/false
	a := make([]I, ns*np)             // array slice(view), only the header is passed, array is shared

	// warmup
	if parallel {
		psum(ns, np, ptil(ns, np, a))
	} else {
		sum(til(a, 0))
	}

	for i := 0; i < 8; i++ {
		if parallel {
			t := time.Now()
			a = ptil(ns, np, a)
			d := time.Since(t)
			fmt.Printf("%v %v\n", d, psum(ns, np, a))
		} else {
			t := time.Now()
			a = til(a, 0)
			d := time.Since(t)
			fmt.Printf("%v %v\n", d, sum(a))
		}
	}
}

func til(a []I, o I) []I {
	for i := range a {
		a[i] = I(i) + o
	}
	return a
}
func sum(a []I) (r I) {
	for _, v := range a {
		r += v
	}
	return r
}
func ptil(ns, np I, a []I) []I {
	c := make(chan bool)
	nw := 0 // number of running goroutines
	for i := I(0); i < np; i++ {
		go func(i I) { til(a[i*ns:i*ns+ns], i*ns); c <- true }(i) // send ok(done) over channel c
		nw++
	}
	wait(c, nw)
	return a
}
func psum(ns, np I, a []I) (r I) {
	c := make(chan int64) // channels are typed. this one carries the result
	nw := 0
	for i := I(0); i < np; i++ {
		go func(i I) { c <- sum(a[i*ns : i*ns+ns]) }(i) // send the result over the channel
		nw++
	}
	n := 0
	for { // wait and accumulate
		select {
		case s := <-c:
			r += s // sum partial results
			n++
			if n == nw {
				return r
			}
		}
	}
}
func wait(c chan bool, n int) { // wait for n go-routines to send their ok(done).
	i := 0
	for {
		select {
		case <-c: // this blocks until a new message is sent over the channel
			i++
		}
		if i == n {
			return
		}
	}
}
func parse(s string) I {
	n, e := strconv.ParseInt(s, 10, 64)
	f(e)
	return n
}
func parseBool(s string) bool {
	b, e := strconv.ParseBool(s)
	f(e)
	return b
}
func f(err error) {
	if err != nil {
		panic(err)
	}
}
