package main

import (
	"runtime/debug"
	"syscall"
)

// To read command line arguments without importing package os, see go/src/os/{proc.go,exec_windows.go}
// Windows needs special handling, see: exec_windows.go (init)
// 	p := syscall.GetCommandLine()
// 	cmd := syscall.UTF16ToString((*[0xffff]uint16)(unsafe.Pointer(p))[:])
//	cmd is a single string and might need special splitting

func main() {
	ini()
	var buf [1024]byte
	p := buf[:]
	for {
		// Read from stdin without os.Stdin
		// syscall.Stdin is not 0 on windows, but a call to GetStdHandle(-10)
		n, err := syscall.Read(syscall.Stdin, p)
		if err != nil {
			panic(err)
		}
		if n > 0 {
			do(p[:n])
		}
	}
}
func do(s []byte) {
	s = cmd(s)
	defer func() {
		if c := recover(); c != nil {
			println(string(debug.Stack()))
			if s, o := c.(string); o {
				println(s)
			} else if e, o := c.(error); o {
				println(e.Error())
			}
		}
	}()
	ns := k(len(s))
	c := mk(C, ns)
	cc := 8 + c<<2
	copy(m.c[cc:cc+ns], s)
	r := kst(evl(prs(c)))
	rc, nr := 8+r<<2, m.k[r]&atom
	dec(r)
	println(string(m.c[rc : rc+nr]))
}
func cmd(b []byte) []byte {
	if !(len(b) == 3 && b[0] == '\\' && b[2] == '\n') {
		return b
	}
	switch b[1] {
	case 'v':
		return []c("lsv 0\n")
	case 'c':
		return []c("clv 0\n")
	case 'h':
		return []c("help 0\n")
	case '\\':
		panic("bye")
	default:
		return b
	}
}
