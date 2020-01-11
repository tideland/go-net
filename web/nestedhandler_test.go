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
	"tideland.dev/go/audit/environments"
	"tideland.dev/go/net/httpx"
	"tideland.dev/go/net/web"
)

//--------------------
// TESTS
//--------------------

// TestNestedHandlerNoHandler tests the mapping of requests to a
// nested handler w/o sub-handlers.
func TestNestedHandlerNoHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	nh := web.NewNestedHandler()

	wa.Handle("/nh", nh)

	wreq := wa.CreateRequest(http.MethodGet, "/nh")
	wresp := wreq.Do()

	wresp.AssertStatusCodeEquals(http.StatusNotFound)
	wresp.AssertBodyMatches("")
}

// TestNestedHandler tests the mapping of requests to a number of
// nested individual handlers.
func TestNestedHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	nh := web.NewNestedHandler()

	nh.AppendHandlerFunc("foo", func(w http.ResponseWriter, r *http.Request) {
		reply := ""
		if f, ok := httpx.PathAt(r, 0); ok {
			reply = f
		}
		if f, ok := httpx.PathAt(r, 1); ok {
			reply += " && " + f
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
		_, err := w.Write([]byte(reply))
		assert.NoError(err)
	})
	nh.AppendHandlerFunc("bar", func(w http.ResponseWriter, r *http.Request) {
		reply := ""
		if f, ok := httpx.PathAt(r, 2); ok {
			reply = f
		}
		if f, ok := httpx.PathAt(r, 3); ok {
			reply += " && " + f
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
		_, err := w.Write([]byte(reply))
		assert.NoError(err)
	})

	wa.Handle("/foo/", nh)

	tests := []struct {
		path       string
		statusCode int
		body       string
	}{
		{
			path:       "/",
			statusCode: http.StatusNotFound,
			body:       "",
		}, {
			path:       "/foo/",
			statusCode: http.StatusOK,
			body:       "foo",
		}, {
			path:       "/foo/4711",
			statusCode: http.StatusOK,
			body:       "foo && 4711",
		}, {
			path:       "/foo/4711/bar",
			statusCode: http.StatusOK,
			body:       "bar",
		}, {
			path:       "/foo/4711/bar/1",
			statusCode: http.StatusOK,
			body:       "bar && 1",
		}, {
			path:       "/foo/4711/bar/1/nothingelse",
			statusCode: http.StatusNotFound,
			body:       "",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.path)
		wreq := wa.CreateRequest(http.MethodGet, test.path)
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

// EOF
