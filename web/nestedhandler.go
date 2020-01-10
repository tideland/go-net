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
	"strings"
)

//--------------------
// NESTED HANDLER
//--------------------

// NestedHandler allows to nest handler following the
// pattern /handlerA/{entityID-A}/handlerB/{entityID-B}.
type NestedHandler struct {
	handlerIDs  []string
	handlers    []http.Handler
	handlersLen int
}

// NewNestedHandler creates an empty nested handler.
func NewNestedHandler() *NestedHandler {
	return &NestedHandler{}
}

// AppendHandler adds one handler to the stack of handlers.
func (nh *NestedHandler) AppendHandler(id string, h http.Handler) {
	nh.handlerIDs = append(nh.handlerIDs, id)
	nh.handlers = append(nh.handlers, h)
	nh.handlersLen++
}

// AppendHandlerFunc adds one handler function to the stack of handlers.
func (nh *NestedHandler) AppendHandlerFunc(id string, hf func(http.ResponseWriter, *http.Request)) {
	nh.AppendHandler(id, http.HandlerFunc(hf))
}

// ServeHTTP implements http.Handler.
func (nh *NestedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := nh.handler(r.URL.Path)
	handler.ServeHTTP(w, r)
}

// handler retrieves the correct handler from the stack.
func (nh *NestedHandler) handler(path string) http.Handler {
	path = strings.Trim(path, "/")
	fields := strings.Split(path, "/")
	fieldsLen := len(fields)
	index := (fieldsLen - 1) / 2
	if (fieldsLen == 1 && fields[0] == "") || index >= nh.handlersLen {
		return http.NotFoundHandler()
	}
	return nh.handlers[index]
}

// EOF
