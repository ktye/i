package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", "127.0.0.1:2019", "")
	flag.Parse()

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		println(err)
		os.Exit(1)
	}
}
func red(x k) (r k) { // 1:x
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		panic("type")
	}
	if n == 0 {
		panic("stdin")
	}
	xp := ptr(x, C)
	name := string(m.c[xp : xp+n])
	if b, o := files[name]; o == false {
		panic("∄ " + name)
	} else {
		return decr(x, mkb(b))
	}
}

var stdout *bytes.Buffer

func wrt(x, y k) (r k) { // x 1:y
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		panic("type")
	}
	if n != 0 {
		panic("nyi: write to virtual file")
	}
	t, n = typ(y)
	if t != C || n == atom {
		panic("type")
	}
	yp := 8 + y<<2
	stdout.Write(m.c[yp : yp+n])
	return decr(y, x)
}

var files map[string][]byte // uploaded (dropped) files

func handler(w http.ResponseWriter, r *http.Request) {
	if method := r.Method; method == "GET" { // send a new front-end
		files = make(map[string][]byte)
		ini()
		table[21] = red
		table[21+dyad] = wrt
		w.Write([]byte(h))
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		return
	}
	if name := r.Header.Get("file"); name != "" {
		files[name] = body
		w.Write([]byte("w " + name + "\n"))
		return
	}
	stdout = bytes.NewBuffer(nil)
	evp(mkb(body))
	w.Write(stdout.Bytes())
}
