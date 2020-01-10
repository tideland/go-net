// Tideland Go Network - HTTP Extensions - Unit Tests
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx_test // import "tideland.dev/go/net/httpx_test"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/net/httpx"
)

//--------------------
// TESTS
//--------------------

// TestAcceptsContentType tests if the checking for accepted content
// types works correctly.
func TestAcceptsContentType(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	h := http.Header{}

	h.Set("Accept", "text/plain; q=0.5, text/html")

	assert.True(httpx.AcceptsContentType(h, httpx.ContentTypePlain))
	assert.True(httpx.AcceptsContentType(h, httpx.ContentTypeHTML))
	assert.False(httpx.AcceptsContentType(h, httpx.ContentTypeJSON))
}

// EOF
