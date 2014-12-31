package charybdis

/*
#include <ctype.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define HOSTLEN 100

void do_host_cloak_host(const char *inbuf, char *outbuf);
void do_host_cloak_ip(const char *inbuf, char *outbuf);

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
*/
import "C"
import "unsafe"

// CloakHost will apply the charybdis cloaking function to a given string
// and return its result. The limit of the string given in is 100.
func CloakHost(host string) (result string) {
	cstring := C.CString(host)
	defer C.free(unsafe.Pointer(cstring))

	cresult := C.CString(host)
	defer C.free(unsafe.Pointer(cresult))

	C.do_host_cloak_host(cstring, cresult)

	return C.GoString(cresult)
}

// CloakIP will apply the charybdis ip address cloaking function to a
// given string and return its result. This should be an IP address.
func CloakIP(host string) (result string) {
	cstring := C.CString(host)
	//defer C.free(unsafe.Pointer(cstring))

	cresult := C.CString(host)
	//defer C.free(unsafe.Pointer(cresult))

	C.do_host_cloak_ip(cstring, cresult)

	return C.GoString(cresult)
}
