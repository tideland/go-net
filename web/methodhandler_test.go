// Tideland Go Network - Web - Unit Tests
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package web_test // import "tideland.dev/go/net/web_test"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/net/web"
)

//--------------------
// TESTS
//--------------------

// TestInvalidMethodHandler tests the panics for invalid values
// passed to a MethodHandler.
func TestInvalidMethodHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mh := web.NewMethodHandler()

	assert.Panics(func() {
		mh.HandleFunc("", makeMethodEcho(assert))
	}, "invalid HTTP method")

	assert.Panics(func() {
		mh.HandleFunc("DO-SOMETHING", makeMethodEcho(assert))
	}, "invalid HTTP method")

	assert.Panics(func() {
		mh.HandleFunc(http.MethodGet, nil)
	}, "need handler function")

	mh.HandleFunc(http.MethodGet, makeMethodEcho(assert))

	assert.Panics(func() {
		mh.HandleFunc(http.MethodGet, makeMethodEcho(assert))
	}, "multiple registrations for GET")
}

// TestMethodHandler tests the multiplexing of methods to different handler.
func TestMethodHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	mh := web.NewMethodHandler()

	mh.HandleFunc(http.MethodGet, makeMethodEcho(assert))
	mh.HandleFunc(http.MethodPatch, makeMethodEcho(assert))
	mh.HandleFunc(http.MethodOptions, makeMethodEcho(assert))

	wa.Handle("/mh/", mh)

	tests := []struct {
		method     string
		statusCode int
		body       string
	}{
		{
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			body:       "METHOD: GET!",
		}, {
			method:     http.MethodHead,
			statusCode: http.StatusMethodNotAllowed,
			body:       "",
		}, {
			method:     http.MethodPost,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodPut,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodPatch,
			statusCode: http.StatusOK,
			body:       "METHOD: PATCH!",
		}, {
			method:     http.MethodDelete,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodConnect,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			body:       "METHOD: OPTIONS!",
		}, {
			method:     http.MethodTrace,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.method)
		wreq := wa.CreateRequest(test.method, "/mh/")
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

// EOF
