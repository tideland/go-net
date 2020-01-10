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
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"time"

	"tideland.dev/go/net/httpx"
	"tideland.dev/go/net/jwt/cache"
	"tideland.dev/go/net/jwt/token"
)

//--------------------
// JWT HANDLER
//--------------------

// JWTHandlerConfig allows to control how the JWT handler works.
// All values are optional. In this case tokens are only decoded
// without using a cache, validated for the current time plus/minus
// a minute leeway, and there's no user defined gatekeeper function
// running afterwards.
type JWTHandlerConfig struct {
	Cache      *cache.Cache
	Key        token.Key
	Leeway     time.Duration
	Gatekeeper func(w http.ResponseWriter, r *http.Request, claims token.Claims) error
}

// JWTHandler checks for a valid token and then runs
// a gatekeeper function.
type JWTHandler struct {
	handler    http.Handler
	cache      *cache.Cache
	key        token.Key
	leeway     time.Duration
	gatekeeper func(w http.ResponseWriter, r *http.Request, claims token.Claims) error
}

// NewJWTHandler creates a handler checking for a valid JSON
// Web Token in each request.
func NewJWTHandler(handler http.Handler, config *JWTHandlerConfig) *JWTHandler {
	jw := &JWTHandler{
		handler: handler,
		leeway:  time.Minute,
	}
	if config != nil {
		if config.Cache != nil {
			jw.cache = config.Cache
		}
		if config.Key != nil {
			jw.key = config.Key
		}
		if config.Leeway != 0 {
			jw.leeway = config.Leeway
		}
		if config.Gatekeeper != nil {
			jw.gatekeeper = config.Gatekeeper
		}
	}
	return jw
}

// ServeHTTP implements the http.Handler interface. It checks for an existing
// and valid token before calling the wrapped handler.
func (jw *JWTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if jw.isAuthorized(w, r) {
		jw.handler.ServeHTTP(w, r)
	}
}

// isAuthorized checks the request for a valid token and if configured
// asks the gatekeepr if the request may pass.
func (jw *JWTHandler) isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	var jwt *token.JWT
	var err error
	switch {
	case jw.cache != nil && jw.key != nil:
		jwt, err = jw.cache.RequestVerify(r, jw.key)
	case jw.cache != nil && jw.key == nil:
		jwt, err = jw.cache.RequestDecode(r)
	case jw.cache == nil && jw.key != nil:
		jwt, err = token.RequestVerify(r, jw.key)
	default:
		jwt, err = token.RequestDecode(r)
	}
	// Now do the checks.
	if err != nil {
		jw.deny(w, r, err.Error(), http.StatusUnauthorized)
		return false
	}
	if jwt == nil {
		jw.deny(w, r, "no JSON Web Token", http.StatusUnauthorized)
		return false
	}
	if !jwt.IsValid(jw.leeway) {
		jw.deny(w, r, "the JSON Web Token claims 'nbf' and/or 'exp' are not valid", http.StatusForbidden)
		return false
	}
	if jw.gatekeeper != nil {
		err := jw.gatekeeper(w, r, jwt.Claims())
		if err != nil {
			jw.deny(w, r, "access rejected by gatekeeper: "+err.Error(), http.StatusUnauthorized)
			return false
		}
	}
	// All fine.
	return true
}

// deny sends a negative feedback to the caller.
func (jw *JWTHandler) deny(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	feedback := map[string]string{
		"statusCode": strconv.Itoa(statusCode),
		"message":    msg,
	}
	switch {
	case httpx.ContainsContentType(r.Header, httpx.ContentTypeJSON):
		b, _ := json.Marshal(feedback)
		w.WriteHeader(statusCode)
		w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeJSON)
		w.Write(b)
	case httpx.ContainsContentType(r.Header, httpx.ContentTypeXML):
		b, _ := xml.Marshal(feedback)
		w.WriteHeader(statusCode)
		w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeXML)
		w.Write(b)
	default:
		w.WriteHeader(statusCode)
		w.Header().Set(httpx.HeaderContentType, httpx.ContentTypePlain)
		w.Write([]byte(msg))
	}
}

// EOF
