package main

import "testing"

var save []byte

func newtest() {
        rand_ = 1592653589
        if save == nil {
                kinit()
                save = make([]byte, len(Bytes))
                copy(save, Bytes)
        } else {
                Bytes = make([]byte, len(save))
                copy(Bytes, save)
                pp, pe, sp = 0, 0, 256
        }
}

func TestCsv(t *testing.T) {
	newtest()	
}

const file1 = `
1,90,2,270
2,90,3,180
`
const k1 = `(1a90 2a90;2a270 3a180)`
