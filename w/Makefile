.PHONY: index.html
index.html: s h
	sh mkhtml

deploy:
	echo TODO

o: s
	⍳ s 2>&1 | xxd -r -p >o

# https://github.com/WebAssembly/wabt/releases
WABT=~/wabt
validate: o
	$(WABT)/wasm-validate o

dump:	o
	$(WABT)/wasm-objdump -s o
