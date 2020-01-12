extern unsigned char __heap_base;

unsigned int b = (unsigned int)&__heap_base;
void* grw(int n) {
	unsigned int r = b;
	b += n;
	return (void *)r;
}


// https://dassur.ma/things/c-to-webassembly/
// scroll way down to "LLVM's memory model"

// https://stackoverflow.com/questions/58252467/heap-base-seems-to-be-missing-in-clang-9-0-0-is-there-a-replacement
