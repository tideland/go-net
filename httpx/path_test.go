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

// TestPathParts tests the splitting of request paths into its parts.
func TestPathParts(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	assert.NoError(err)
	ps := httpx.PathParts(r)
	assert.Length(ps, 0)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo", nil)
	assert.NoError(err)
	ps = httpx.PathParts(r)
	assert.Length(ps, 1)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo/bar", nil)
	assert.NoError(err)
	ps = httpx.PathParts(r)
	assert.Length(ps, 2)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo//bar", nil)
	assert.NoError(err)
	ps = httpx.PathParts(r)
	assert.Length(ps, 2)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo/bar?yadda=1", nil)
	assert.NoError(err)
	ps = httpx.PathParts(r)
	assert.Length(ps, 2)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo/bar?yadda=1/2", nil)
	assert.NoError(err)
	ps = httpx.PathParts(r)
	assert.Length(ps, 2)
}

// TestPathAt tests the checking and extrecting of a part out of
// a request path.
func TestPathAt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	assert.NoError(err)
	assert.Panics(func() {
		httpx.PathAt(r, -1)
	})

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo", nil)
	assert.NoError(err)
	p, ok := httpx.PathAt(r, 0)
	assert.True(ok)
	assert.Equal(p, "foo")
	p, ok = httpx.PathAt(r, 1)
	assert.False(ok)
	assert.Equal(p, "")

	r, err = http.NewRequest(http.MethodGet, "http://localhost/orders/4711/items/1", nil)
	assert.NoError(err)
	p, ok = httpx.PathAt(r, 1)
	assert.True(ok)
	assert.Equal(p, "4711")
	p, ok = httpx.PathAt(r, 3)
	assert.True(ok)
	assert.Equal(p, "1")
	p, ok = httpx.PathAt(r, 5)
	assert.False(ok)
	assert.Equal(p, "")
}

// EOF
