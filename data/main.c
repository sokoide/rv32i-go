char is_even(int a);
void _out(int);

char is_even(int a) { return a % 2 == 0; }
int main();

void riscv32_boot() {
    // Note:
    // This is only for the simulator.
    // Need to clear bss, setup timer/UART or etc
    // on a real hardware.
    main();
}

int fib(int i) {
    if (i <= 1)
        return i;
    return fib(i - 1) + fib(i - 2);
}

int main() {
    fib(30);

    return 0;

    // int a, b;
    // char c, d;
    // a = 1000000000;
    // b = 1;
    // c = is_even(a);
    // d = is_even(b);
    // _out(a);
    // _out(b);
    // _out(c);
    // _out(d);
    // return 0;
}
