// Tideland Go Network - Web
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package web // import "tideland.dev/go/net/web"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"

	"tideland.dev/go/net/httpx"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// MethodAll allows to define a handle for all non-explicit
	// configured HTTP methods.
	MethodAll = "*"
)

//--------------------
// METHOD HANDLER
//--------------------

// MethodHandler distributes request depending on the HTTP method
// to subhandlers.
type MethodHandler struct {
	handlers map[string]http.Handler
}

// NewMethodHandler creates an empty HTTP method handler.
func NewMethodHandler() *MethodHandler {
	return &MethodHandler{
		handlers: make(map[string]http.Handler),
	}
}

// Handle adds the handler based on the method.
func (mh *MethodHandler) Handle(method string, handler http.Handler) {
	if !httpx.IsValidMethod(method) {
		panic("invalid HTTP method")
	}
	if handler == nil {
		panic("need handler")
	}
	if _, exist := mh.handlers[method]; exist {
		panic("multiple registrations for " + method)
	}
	mh.handlers[method] = handler
}

// HandleFunc adds the handler function based on the method.
func (mh *MethodHandler) HandleFunc(method string, hf func(http.ResponseWriter, *http.Request)) {
	if hf == nil {
		panic("need handler function")
	}
	mh.Handle(method, http.HandlerFunc(hf))
}

// ServeHTTP implements http.Handler.
func (mh *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := mh.handlers[r.Method]
	if !ok {
		handler, ok = mh.handlers[MethodAll]
		if !ok {
			http.Error(w, "no matching method handler found", http.StatusMethodNotAllowed)
			return
		}
	}
	handler.ServeHTTP(w, r)
}

// EOF
