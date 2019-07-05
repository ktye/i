.PHONY: index.html
index.html: s h
	sh mkhtml

deploy: index.html
	cp index.html ../../ktye.github.io

o: s
	⍳ s 2>&1 | xxd -r -p >o

# https://github.com/WebAssembly/wabt/releases
WABT=~/wabt
valid: o
	$(WABT)/wasm-validate o

dump:	o
	$(WABT)/wasm-objdump -sd o

wat:	o
	$(WABT)/wasm2wat o
	
