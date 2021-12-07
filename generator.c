#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
    int i, k;
    unsigned long long x, y;
    if (argc == 1) {
        fprintf(stderr, "Usage: %s <k>\n", argv[0]);
        return 1;
    }
    k = atoi(argv[1]);
    for (x = 0; x < 1ULL<<(2*k); ++x) {
        for (i = 0, y = x; i < k; ++i, y >>= 2)
            putchar("ACGT"[y&3]);
        putchar('\n');
    }
    return 0;
}
