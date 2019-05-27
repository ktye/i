.PHONY: index.html
index.html: s h
	sh mkhtml

deploy:
	echo TODO

o: s
	⍳ s 2>&1 | xxd -r -p >o

# https://github.com/WebAssembly/wabt/releases
WABT=/d/wabt
valid: o
	$(WABT)/wasm-validate o

dump:	o
	$(WABT)/wasm-objdump -s o

wat:	o
	$(WABT)/wasm2wat o
	