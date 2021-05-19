package k

import (
	"fmt"
	"testing"

	. "github.com/ktye/wg/module"
)

//

func TestAlloc(t *testing.T) {
	n := 64 * 1024
	if len(Bytes) != 64*1024 {
		t.Fatal("memory size")
	}

	copy(Bytes, make([]byte, n))
	minit(10, 16)
	fmt.Println("g20", I32(20))
	fmt.Println(alloc(10))
	fmt.Println("g20", I32(20))
}
func TestBucket(t *testing.T) {
	tc := []struct{ in, exp int32 }{
		{0, 4},
		{4, 4},
		{8, 4},
		{9, 5},
		{24, 5},
		{25, 6},
	}
	for _, tc := range tc {
		if got := bucket(tc.in); got != tc.exp {
			t.Fatalf("bucket %d => %d (exp %d)\n", tc.in, got, tc.exp)
		}
	}
}
