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
	go test -tags=small -coverprofile=cov.out
	go tool cover -html=cov.out -o cov.html
fi

wg        -nomain . > k.wat
wg -small -nomain . > s.wat

/c/local/wabt/wat2wasm -o /c/k/ktye.github.io/k.wasm k.wat

cp k.t          /c/k/ktye.github.io/k.t
#cp apl/apl.*    /c/k/ktye.github.io/

#rm k.k out k.wat
