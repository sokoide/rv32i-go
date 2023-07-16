char is_even(int a);
void _out(int);

char is_even(int a) { return a % 2 == 0; }
int main();

void riscv32_boot() {
    // TODO: need to clear bss, setup timer/UART or etc
    main();
}

int main() {
    int a, b;
    char c, d;
    a = 10;
    b = 1;
    c = is_even(a);
    d = is_even(b);
    _out(a);
    _out(b);
    _out(c);
    _out(d);
    return 0;
}