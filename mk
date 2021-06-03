set -x
set -e
go install
wg . > k.wat
/c/local/wabt/wat2wasm --enable-bulk-memory --enable-simd k.wat

cp k.html /c/k/ktye.github.io/
cp k.wasm /c/k/ktye.github.io/
