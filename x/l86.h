// -nostdlib -ffreestanding

//gcc -nostdlib -static

asm(
        ".global _start\n"
        "_start:\n"
        "   xorl %ebp,%ebp\n"       // mark outermost stack frame
        "   movq 0(%rsp),%rdi\n"    // get argc
        "   lea 8(%rsp),%rsi\n"     // the arguments are pushed just below, so argv = %rbp + 8
        "   call bare_main\n"       // call our bare_main
        "   movq %rax,%rdi\n"       // take the main return code and use it as first argument for...
        "   movl $60,%eax\n"        // ... the exit syscall
        "   syscall\n"
        "   int3\n");               // just in case

asm(
        "bare_write:\n"             // write syscall wrapper; the calling convention is pretty much ok as is
        "   movq $1,%rax\n"         // 1 = write syscall on x86_64
        "   syscall\n"
        "   ret\n");


int bare_write(int fd, const void *buf, unsigned count);

unsigned my_strlen(const char *ch) {
    const char *ptr;
    for(ptr = ch; *ptr; ++ptr);
    return ptr-ch;
}

int bare_main(int argc, char *argv[]) {
    for(int i = 0; i < argc; ++i) {
        int len = my_strlen(argv[i]);
        bare_write(1, argv[i], len);
        bare_write(1, "\n", 1);
    }
    return argc;
}
