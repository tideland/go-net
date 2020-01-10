// Tideland Go Network - HTTP Extensions
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx // import "tideland.dev/go/net/httpx"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
)

//--------------------
// CONSTANTS
//--------------------

// httpMethods contains all valid HTTP methods.
var httpMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodConnect: {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

//--------------------
// METHOD TOOLS
//--------------------

// IsValidMethod returns true if the passed method is valid.
func IsValidMethod(method string) bool {
	_, valid := httpMethods[method]
	return valid
}

// EOF
