set -x
set -e
go install
wg . > k.wat
/c/local/wabt/wat2wasm --enable-bulk-memory --enable-simd k.wat

cp k.html /c/k/ktye.github.io/index.html
cp k.wasm /c/k/ktye.github.io/k.wasm
cp k.t    /c/k/ktye.github.io/k.t
cp readme /c/k/ktye.github.io/readme


# wavm:
# wk='rlwrap -H /dev/null /c/local/wavm/bin/wavm run --mount-root . /c/k/i/k.wat'

