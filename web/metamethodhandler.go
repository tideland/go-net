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
)

//--------------------
// METHOD HANDLER INTERFACES
//--------------------

// GetHandler has to be implemented by a handler for HEAD requests
// dispatched through the MetaMethodHandler.
type GetHandler interface {
	ServeHTTPGet(w http.ResponseWriter, r *http.Request)
}

// HeadHandler has to be implemented by a handler for HEAD requests
// dispatched through the MetaMethodHandler.
type HeadHandler interface {
	ServeHTTPHead(w http.ResponseWriter, r *http.Request)
}

// PostHandler has to be implemented by a handler for POST requests
// dispatched through the MetaMethodHandler.
type PostHandler interface {
	ServeHTTPPost(w http.ResponseWriter, r *http.Request)
}

// PutHandler has to be implemented by a handler for PUT requests
// dispatched through the MetaMethodHandler.
type PutHandler interface {
	ServeHTTPPut(w http.ResponseWriter, r *http.Request)
}

// PatchHandler has to be implemented by a handler for PATCH requests
// dispatched through the MetaMethodHandler.
type PatchHandler interface {
	ServeHTTPPatch(w http.ResponseWriter, r *http.Request)
}

// DeleteHandler has to be implemented by a handler for DELETE requests
// dispatched through the MetaMethodHandler.
type DeleteHandler interface {
	ServeHTTPDelete(w http.ResponseWriter, r *http.Request)
}

// ConnectHandler has to be implemented by a handler for CONNECT requests
// dispatched through the MetaMethodHandler.
type ConnectHandler interface {
	ServeHTTPConnect(w http.ResponseWriter, r *http.Request)
}

// OptionsHandler has to be implemented by a handler for OPTIONS requests
// dispatched through the MetaMethodHandler.
type OptionsHandler interface {
	ServeHTTPOptions(w http.ResponseWriter, r *http.Request)
}

// TraceHandler has to be implemented by a handler for TRACE requests
// dispatched through the MetaMethodHandler.
type TraceHandler interface {
	ServeHTTPTrace(w http.ResponseWriter, r *http.Request)
}

//--------------------
// META METHOD HANDLER
//--------------------

// MetaMethodHandler checks if the handler contains a handler with
// a matching interface for the HTTP method.
type MetaMethodHandler struct {
	handler http.Handler
}

// NewMetaMethodHandler creates a meta HTTP method handler.
func NewMetaMethodHandler(handler http.Handler) *MetaMethodHandler {
	if handler == nil {
		panic("need handler")
	}
	return &MetaMethodHandler{
		handler: handler,
	}
}

// ServeHTTP implements http.Handler.
func (mmh *MetaMethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h, ok := mmh.handler.(GetHandler); ok {
			h.ServeHTTPGet(w, r)
			return
		}
	case http.MethodHead:
		if h, ok := mmh.handler.(HeadHandler); ok {
			h.ServeHTTPHead(w, r)
			return
		}
	case http.MethodPost:
		if h, ok := mmh.handler.(PostHandler); ok {
			h.ServeHTTPPost(w, r)
			return
		}
	case http.MethodPut:
		if h, ok := mmh.handler.(PutHandler); ok {
			h.ServeHTTPPut(w, r)
			return
		}
	case http.MethodPatch:
		if h, ok := mmh.handler.(PatchHandler); ok {
			h.ServeHTTPPatch(w, r)
			return
		}
	case http.MethodDelete:
		if h, ok := mmh.handler.(DeleteHandler); ok {
			h.ServeHTTPDelete(w, r)
			return
		}
	case http.MethodConnect:
		if h, ok := mmh.handler.(ConnectHandler); ok {
			h.ServeHTTPConnect(w, r)
			return
		}
	case http.MethodOptions:
		if h, ok := mmh.handler.(OptionsHandler); ok {
			h.ServeHTTPOptions(w, r)
			return
		}
	case http.MethodTrace:
		if h, ok := mmh.handler.(TraceHandler); ok {
			h.ServeHTTPTrace(w, r)
			return
		}
	}
	mmh.handler.ServeHTTP(w, r)
}

// EOF
