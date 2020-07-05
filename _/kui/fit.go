package main

import (
	"bytes"
	"fmt"

	"github.com/tormoder/fit"
)

func decodeFit(b []byte) {
	fit, err := fit.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	fmt.Println(fit)
}
