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

// RFC 1459 maximum hostname length
#define HOSTLEN 100

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

void
do_host_cloak_host(const char *inbuf, char *outbuf)
{
	char b26_alphabet[] = "abcdefghijklmnopqrstuvwxyz";
	char *tptr;
	uint32_t accum = fnv_hash((const unsigned char*) inbuf, 32);

	strncpy(outbuf, inbuf, HOSTLEN + 1);

	for (tptr = outbuf; *tptr != '\0'; tptr++)
	{
		if (*tptr == '.')
			break;

		if (isdigit(*tptr) || *tptr == '-')
			continue;

		*tptr = b26_alphabet[(*tptr + accum) % 26];

		accum = (accum << 1) | (accum >> 31);
	}

	for (tptr = outbuf; *tptr != '\0'; tptr++)
	{
		if (isdigit(*tptr))
			*tptr = '0' + (*tptr + accum) % 10;

		accum = (accum << 1) | (accum >> 31);
	}
}

void
do_host_cloak_ip(const char *inbuf, char *outbuf)
{
	char chartable[] = "ghijklmnopqrstuvwxyz";
	char *tptr;
	uint32_t accum = fnv_hash((const unsigned char*) inbuf, 32);
	int sepcount = 0;
	int totalcount = 0;
	int ipv6 = 0;

	strncpy(outbuf, inbuf, 100 + 1);

	if (strchr(outbuf, ':'))
	{
		ipv6 = 1;

		for (tptr = outbuf; *tptr != '\0'; tptr++)
			if (*tptr == ':')
				totalcount++;
	}
	else if (!strchr(outbuf, '.'))
		return;

	for (tptr = outbuf; *tptr != '\0'; tptr++)
	{
		if (*tptr == ':' || *tptr == '.')
		{
			sepcount++;
			continue;
		}

		if (ipv6 && sepcount < totalcount / 2)
			continue;

		if (!ipv6 && sepcount < 2)
			continue;

		*tptr = chartable[(*tptr + accum) % 20];
		accum = (accum << 1) | (accum >> 31);
	}
}

