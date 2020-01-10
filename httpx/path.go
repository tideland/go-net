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
	"strings"
)

//--------------------
// PATH TOOLS
//--------------------

// PathParts splits the request path into its parts.
func PathParts(r *http.Request) []string {
	rawParts := strings.Split(r.URL.Path, "/")
	parts := []string{}
	for _, part := range rawParts {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

// PathAt returns the nth part of the request path and true
// if it exists. Otherwise an empty string and false.
func PathAt(r *http.Request, n int) (string, bool) {
	if n < 0 {
		panic("illegal path index")
	}
	parts := PathParts(r)
	if len(parts) < n+1 {
		return "", false
	}
	return parts[n], true
}

// EOF
