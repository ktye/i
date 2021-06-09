set -x
set -e
go install
wg . > k.wat
/c/local/wabt/wat2wasm --enable-bulk-memory --enable-simd k.wat

cp k.html /c/k/ktye.github.io/index.html
cp k.wasm /c/k/ktye.github.io/k.wasm


# wavm:
# wk='rlwrap -H /dev/null /c/local/wavm/bin/wavm run --mount-root . /c/k/i/k.wat'

