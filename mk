set -x
set -e

# embed z.k in z.go
zn=`wc -c z.k | sed 's/ .*//'`
zk=`sed -e 's/\\\/\\\\\\\/g' -e 's/"/\\\"/g' -e 's/$/\\\/g' z.k | tr '\n' 'n'`
cat << EOF > z.go
package main

import . "github.com/ktye/wg/module"

func zk() {
	Data(600, "$zk")
	zn := int32($zn) // should end before 8k
	x := mk(Ct, zn)
	Memorycopy(int32(x), 600, zn)
	dx(Val(x))
}
EOF

go install

if [ "$1" = "cover" ]; then
	go test -coverprofile=cov.out
	go tool cover -html=cov.out -o cov.html
fi

wg . > k.wat
/c/local/wabt/wat2wasm --enable-bulk-memory --enable-simd k.wat

cp k.html /c/k/ktye.github.io/index.html
cp k.wasm /c/k/ktye.github.io/k.wasm
cp k.t    /c/k/ktye.github.io/k.t
cp readme /c/k/ktye.github.io/readme
cp ∘.k    /c/k/ktye.github.io/∘.k


# wavm:
# wk='rlwrap -H /dev/null /c/local/wavm/bin/wavm run --mount-root . /c/k/i/k.wat'

