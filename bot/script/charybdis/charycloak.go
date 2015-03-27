package charybdis

/*
#include <ctype.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void do_host_cloak_host(const char *inbuf, char *outbuf);
void do_host_cloak_ip(const char *inbuf, char *outbuf);
*/
import "C"

import "unsafe"

// CloakHost will apply the charybdis cloaking function to a given string
// and return its result. The limit of the string given in is 100.
func CloakHost(host string) (result string) {
	cstring := C.CString(host)
	cresult := C.CString(host)

	defer C.free(unsafe.Pointer(cstring))
	defer C.free(unsafe.Pointer(cresult))

	C.do_host_cloak_host(cstring, cresult)

	return C.GoString(cresult)
}

// CloakIP will apply the charybdis ip address cloaking function to a
// given string and return its result. This should be an IP address.
func CloakIP(host string) (result string) {
	cstring := C.CString(host)
	cresult := C.CString(host)

	defer C.free(unsafe.Pointer(cstring))
	defer C.free(unsafe.Pointer(cresult))

	C.do_host_cloak_ip(cstring, cresult)

	return C.GoString(cresult)
}
