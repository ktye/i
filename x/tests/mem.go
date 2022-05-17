package mem

import . "github.com/ktye/wg/module"

func init() {
	Memory(1)
	Data(0, "\x61\x62\x63\x64\x01\x02\x03")
}
