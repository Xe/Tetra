// Package xmlrpc implements a simple xmlrpc client. This is modified to work
// with the oddities that Atheme's XMLRPC server has. Do not use this fork of
// this library for any reason other than to interface with Atheme.
package xmlrpc

import (
	"fmt"
)

// xmlrpcError represents errors returned on xmlrpc request.
type xmlrpcError struct {
	Code int
	Err  string
}

// Error() method implements Error interface
func (e *xmlrpcError) Error() string {
	return fmt.Sprintf("error: \"%s\" code: %d", e.Err, e.Code)
}

// Base64 represents value in base64 encoding
type Base64 string
