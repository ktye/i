set -x
set -e

#wg -tags simd4 .
#wg -c . > /tmp/k/k.c && cc -O3 -Wall /tmp/k/k.c -Wfatal-errors -lm
#wg -tags simd5 -c . > /tmp/k/v.c && clang-17 -O3 -Wall -mavx2 /tmp/k/v.c -Wfatal-errors -lm
#exit 0

# embed z.k in z.go
sed -e '/^\//d' -e 's,  */.*,,' -e 's/^ *//' -e '/^$/d' z.k > k.k
zn=`wc -c k.k | sed 's/ .*//'`
zk=`sed -e 's/\\\/\\\\\\\/g' -e 's/"/\\\"/g' -e 's/$/\\\/g' k.k | tr '\n' 'n'`
rm k.k
cat << EOF > z.go
package main

import . "github.com/ktye/wg/module"

func zk() {
	Data(280, "$zk")
	zn := int32($zn) // should end before 2k
	x := mk(Ct, zn)
	Memorycopy(int32(x), 280, zn)
	dx(Val(x))
}
EOF

go install

if [ "$1" = "cover" ]; then
	go test -coverprofile=cov.out
	go tool cover -html=cov.out -o cov.html
fi

if [ "$1" = "web" ]; then
	wg -tags small -nomain . | wat2wasm - -o /c/k/ktye.github.io/k.wasm
	wasm-opt -Oz --enable-bulk-memory   /c/k/ktye.github.io/k.wasm -o - | wc -c
	go run ./_/kdoc.go > /c/k/ktye.github.io/kdoc.htm
else
	wg -tags small -nomain . | wat2wasm - -o k.wasm 
	wasm-opt -Oz --enable-bulk-memory   k.wasm -o z.wasm
	zstd -19 -f z.wasm
	wc -c k.wasm z.wasm z.wasm.zst
	rm k.wasm z.wasm z.wasm.zst
fi

