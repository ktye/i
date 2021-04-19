cat *.w | awk -f awk      > k.wat
/c/local/wabt/wat2wasm      k.wat -o k.wasm
/c/local/wabt/wasm-validate          k.wasm

