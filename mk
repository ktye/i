set -x
set -e

# embed z.k in z.go
sed -e '/^\//d' -e 's,  */.*,,' -e 's/^ *//' -e '/^$/d' z.k > k.k
zn=`wc -c k.k | sed 's/ .*//'`
zk=`sed -e 's/\\\/\\\\\\\/g' -e 's/"/\\\"/g' -e 's/$/\\\/g' k.k | tr '\n' 'n'`
cat << EOF > z.go
package main

import . "github.com/ktye/wg/module"

func zk() {
	Data(600, "$zk")
	zn := int32($zn) // should end before 4k
	x := mk(Ct, zn)
	Memorycopy(int32(x), 600, zn)
	dx(Val(x))
}
EOF

if [ "$1" = "kc" ]; then
	wg -c -prefix ktye_ . > k.c
	gcc -O2 k.c
	exit 0
fi

go install

if [ "$1" = "cover" ]; then
	go test -coverprofile=cov.out
	go tool cover -html=cov.out -o cov.html
fi

wg        -nomain . > k.wat


/c/local/wabt/wat2wasm -o k.wasm k.wat
wc -c k.wasm
wasm-opt -Oz --enable-bulk-memory k.wasm -o - | wc -c

if [ "$1" = "web" ]; then
	#cp k.wasm /c/k/ktye.github.io/k.wasm
	#go run _/kdoc.html >/c/k/ktye.github.io/kdoc.html
	ls *.go|grep -v _test.go|xargs cat >/c/k/ktye.github.io/k.go
fi

#rm k.k out k.wat
