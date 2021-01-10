package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const GS = "c:/local/gs/bin/gswin64c.exe"

const T = `1 2 add           == %3
1 2 sub                      == %-1
2 2.0 eq                     == %true
2 3 eq                       == %false
[1]                          == %[1]
1 2 lt                       == %true
(beta) (alpha) lt            == %false
5 6 and                      == %4
3 4 bitshift                 == %48
49 -4 bitshift               == %3
1 %GS<1>`

func TestGS(t *testing.T) {
	if e := ioutil.WriteFile("t.in", []byte(T), 0644); e != nil {
		t.Fatal(e)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	out, e := exec.CommandContext(ctx, GS, "-q", "-f", "t.in").Output()
	if e != nil {
		t.Fatal(e)
	}

	compare(t, string(out), expected(t))
}
func TestPS(t *testing.T) {
	var b strings.Builder
	i := New(&b)
	s := bufio.NewScanner(strings.NewReader(T))
	for s.Scan() {
		i.Run(s.Text())
	}
	fmt.Fprintf(&b, "GS<%d>", len(i.v.stack))
	compare(t, b.String(), expected(t))
}
func compare(t *testing.T, got, exp string) {
	gt, ex := strings.Split(got, "\n"), strings.Split(exp, "\n")
	n := len(gt)
	if n > len(ex) {
		n = len(ex)
	}
	for i := 0; i < n; i++ {
		if gt[i] != ex[i] {
			t.Fatalf("line %d\ngot: %s\nexp: %s\n", i+1, gt[i], ex[i])
		}
	}
	if len(gt) != len(ex) {
		t.Fatalf("got %d lines, expected %d\nout:\n%s", len(gt), len(ex), got)
	}
}
func expected(t *testing.T) string {
	var b strings.Builder
	v := strings.Split(T, "\n")
	for _, s := range v {
		if s == "" {
			fmt.Fprintln(&b, "")
		} else if strings.HasPrefix(s, "%") {
			fmt.Fprintln(&b, s[1:])
		} else {
			i := strings.Index(s, " %")
			if i == -1 {
				t.Fatal("test case")
			}
			fmt.Fprintln(&b, s[i+2:])
		}
	}
	r := b.String()
	return r[:len(r)-1] // final newline
}
