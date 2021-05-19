package k

import . "github.com/ktye/wg/module"

const (
	Err int32 = iota
	Nyi
	Grow
	Unref
)

func trap(x int32) { SetI32(0, x); panic(x) }
