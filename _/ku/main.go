package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/zserge/lorca"
)

func main() {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 480, 320, args...)
	fatal(err)
	defer ui.Close()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	fatal(err)

	defer ln.Close()

	go http.Serve(ln, http.FileServer(http.Dir("c:/k/ktye.github.io")))
	ui.Load(fmt.Sprintf("http://%s/index.html", ln.Addr()))

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
