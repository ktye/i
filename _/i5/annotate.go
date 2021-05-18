package k

import (
	"bytes"
	"fmt"
	godebug "runtime/debug"
	"strings"
)

type ctx struct {
	src     [][]byte
	pos     []*int
	gostack []byte
}

func (c *ctx) push(s []byte, p *int) {
	c.src = append(c.src, s)
	c.pos = append(c.pos, p)
}
func (c *ctx) drop() {
	c.src = c.src[:len(c.src)-1]
	c.pos = c.pos[:len(c.pos)-1]
}
func (c *ctx) restore(f string) {
	if r := recover(); r != nil {
		c.gostack = godebug.Stack()
		c.annotate()
		c.drop()
		panic(r)
	}
	c.drop()
}
func (c ctx) annotate() {
	if n := len(c.src); n > 0 {
		b, p := c.src[n-1], c.pos[n-1]
		l, i := trimLine(b, *p)
		fmt.Println(l)
		sp := ""
		if i > 0 && i < len(l) {
			sp = strings.Repeat(" ", i)
		}
		fmt.Println(sp + "^")
	}
}
func trimLine(b []byte, p int) (string, int) {
	b, p = trimStart(b, p)
	i := bytes.IndexByte(b, 10)
	if i < 0 {
		return string(b), p
	}
	return string(b[:i]), p
}
func trimStart(b []byte, p int) ([]byte, int) {
	for {
		i := bytes.IndexByte(b, 10)
		if i < 0 || p < i {
			return b, p
		}
		b = b[1+i:]
		p -= 1 + i
	}
}
func (c ctx) Gostack() []byte { return c.gostack }
