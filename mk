set -x
set -e

# embed z.k in z.go
zk(){
cat $1 $2 | sed -e '/^\//d' -e 's,  */.*,,' -e 's/^ *//' -e '/^$/d' > k.k
zn=`wc -c k.k | sed 's/ .*//'`
zk=`sed -e 's/\\\/\\\\\\\/g' -e 's/"/\\\"/g' -e 's/$/\\\/g' k.k | tr '\n' 'n'`
cat << EOF > $4
$3

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
}
zk z.k ""  "//go:build !small" z.go
zk z.k s.k "//go:build small"  s.go

if [ "$1" = "kc" ]; then
	wg -c -prefix ktye_ . > k.c
	gcc -O2 k.c
	exit 0
fi

go install

#small: go build -tags small

if [ "$1" = "cover" ]; then
	go test -coverprofile=cov.out
	go tool cover -html=cov.out -o cov.html
fi

wg             . > k.wat
wg -small      . > s.wat
/c/local/wabt/wat2wasm -o web/k.wasm k.wat
/c/local/wabt/wat2wasm -o web/s.wasm s.wat

cp web/k.wasm   /c/k/ktye.github.io/
cp web/*.js     /c/k/ktye.github.io/
cp k.t          /c/k/ktye.github.io/k.t
cp apl/apl.html /c/k/ktye.github.io/
cp apl/apl.k    /c/k/ktye.github.io/

