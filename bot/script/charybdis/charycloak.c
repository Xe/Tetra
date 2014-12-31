/* 
 * Charybdis: an advanced ircd
 * ip_cloaking.c: provide user hostname cloaking
 *
 * Written originally by nenolod, altered to use FNV by Elizabeth in 2008
 */

#include <ctype.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* Magic value for FNV hash functions */
#define FNV1_32_INIT 0x811c9dc5UL

uint32_t
fnv_hash(const unsigned char *s, int bits)
{
        uint32_t h = FNV1_32_INIT;

        while (*s)
        {
                h ^= *s++;
                h += (h<<1) + (h<<4) + (h<<7) + (h << 8) + (h << 24);
        }
        if (bits < 32)
                h = ((h >> bits) ^ h) & ((1<<bits)-1);
        return h;
}
