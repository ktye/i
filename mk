set -x
set -e
go install
wg . > k.wat
/c/local/wabt/wat2wasm k.wat
