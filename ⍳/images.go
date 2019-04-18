package main

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strings"

	"github.com/mattn/go-sixel"
)

// Values that are url-encoded images "data:image/png;base64,..." are displayed as images
// instead of text, depending on the terminal type.
// By default, images are converted to sixel format, which is supported by xterm and mintty.
//
// On wasm, images are set to the img node below the terminal.
//
// The plot package encodes plots in this format with the default stringer.

const pngPrefix = "data:image/png;base64,"

func sxl(s string) string {
	if !strings.HasPrefix(s, pngPrefix) {
		return s
	}
	b, err := base64.StdEncoding.DecodeString(s[len(pngPrefix):])
	if err != nil {
		return err.Error()
	}
	r := bytes.NewReader(b)
	img, err := png.Decode(r)
	if err != nil {
		return err.Error()
	}
	var buf bytes.Buffer
	sixel.NewEncoder(&buf).Encode(img)
	return string(buf.Bytes())
}
