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

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/environments"
)

//--------------------
// WEB ASSERTER AND HELPERS
//--------------------

// StartTestServer initialises and starts the asserter for the tests.
func startWebAsserter(assert *asserts.Asserts) *environments.WebAsserter {
	wa := environments.NewWebAsserter(assert)
	return wa
}

// makeMethodEcho creates a handler echoing the HTTP method.
func makeMethodEcho(assert *asserts.Asserts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := "METHOD: " + r.Method + "!"
		w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
		if _, err := w.Write([]byte(reply)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// EOF
