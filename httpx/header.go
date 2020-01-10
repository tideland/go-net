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
// CONSTANTS
//--------------------

// Headers fields and values.
const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	ContentTypePlain      = "text/plain"
	ContentTypeHTML       = "text/html"
	ContentTypeXML        = "application/xml"
	ContentTypeJSON       = "application/json"
	ContentTypeURLEncoded = "application/x-www-form-urlencoded"

	contentTypesText = "text/"
)

//--------------------
// HEADER TOOLS
//--------------------

// AcceptsContentType checks a request header accepts a given content type.
func AcceptsContentType(h http.Header, contentType string) bool {
	return strings.Contains(h.Get(HeaderAccept), contentType)
}

// ContainsContentType checks if the header contains the given content type.
func ContainsContentType(h http.Header, contentType string) bool {
	return strings.Contains(h.Get(HeaderContentType), contentType)
}

// EOF
