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
	"tideland.dev/go/net/web"
)

//--------------------
// TESTS
//--------------------

// TestInvalidMMetaethodHandler tests the panic if the given handler
// for the MetaMethodHandler is invalid.
func TestInvalidMetaMethodHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	assert.Panics(func() {
		web.NewMetaMethodHandler(nil)
	}, "need handler")
}

// TestMetaMethodHandler tests the wrapping of a handler as a
// MetaMethodHandler.
func TestMethodWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	wa.Handle("/mmh/", web.NewMetaMethodHandler(mmHandler{}))

	tests := []struct {
		method     string
		statusCode int
		body       string
	}{
		{
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodHead,
			statusCode: http.StatusBadRequest,
			body:       "",
		}, {
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       "METHOD: PUT!",
		}, {
			method:     http.MethodPatch,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodDelete,
			statusCode: http.StatusNoContent,
			body:       "",
		}, {
			method:     http.MethodConnect,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodOptions,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodTrace,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.method)
		wreq := wa.CreateRequest(test.method, "/mmh/")
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

//--------------------
// HELPING HANDLER
//--------------------

// mmHandler provides some of the methods for the MetaMethodHandler.
type mmHandler struct{}

func (h mmHandler) ServeHTTPPut(w http.ResponseWriter, r *http.Request) {
	reply := "METHOD: " + r.Method + "!"
	w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
	if _, err := w.Write([]byte(reply)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h mmHandler) ServeHTTPDelete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNoContent)
}

func (h mmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

// EOF
