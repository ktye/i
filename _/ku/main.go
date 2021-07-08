package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/zserge/lorca"
)

//go:embed www
var fsys embed.FS

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
	sub, err := fs.Sub(fsys, "www")
	fatal(err)

	go http.Serve(ln, http.FileServer(http.FS(sub)))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

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
