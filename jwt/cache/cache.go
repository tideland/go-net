// Tideland Go Network - JSON Web Token - Cache
//
// Copyright (C) 2016-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cache // import "tideland.dev/go/net/jwt/cache"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"net/http"
	"strings"
	"time"

	"tideland.dev/go/net/jwt/token"
	"tideland.dev/go/trace/failure"
	"tideland.dev/go/trace/logger"
)

//--------------------
// CACHE ENTRY
//--------------------

// cacheEntry manages a token and its access time.
type cacheEntry struct {
	jwt      *token.JWT
	accessed time.Time
}

//--------------------
// CACHE
//--------------------

// defaultTimeout is the default timeout for synchronous actions.
const defaultTimeout = 5 * time.Second

// Cache provides a caching for tokens so that these
// don't have to be decoded or verified multiple times.
type Cache struct {
	ctx        context.Context
	entries    map[string]*cacheEntry
	ttl        time.Duration
	leeway     time.Duration
	interval   time.Duration
	maxEntries int
	actionc    chan func()
}

// New creates a new JWT caching. The ttl value controls
// the time a cached token may be unused before cleanup. The
// leeway is used for the time validation of the token itself.
// The duration of the interval controls how often the background
// cleanup is running. Final configuration parameter is the maximum
// number of entries inside the cache. If these grow too fast the
// ttl will be temporarily reduced for cleanup.
func New(ctx context.Context, ttl, leeway, interval time.Duration, maxEntries int) *Cache {
	c := &Cache{
		ctx:        ctx,
		entries:    map[string]*cacheEntry{},
		ttl:        ttl,
		leeway:     leeway,
		interval:   interval,
		maxEntries: maxEntries,
		actionc:    make(chan func(), 1),
	}
	go c.backend()
	return c
}

// Get tries to retrieve a token from the cache.
func (c *Cache) Get(st string) (*token.JWT, error) {
	var jwt *token.JWT
	aerr := c.doSync(func() {
		if c.entries == nil {
			return
		}
		entry, ok := c.entries[st]
		if !ok {
			return
		}
		if !entry.jwt.IsValid(c.leeway) {
			// Remove invalid token.
			delete(c.entries, st)
		}
		entry.accessed = time.Now()
		jwt = entry.jwt
	}, defaultTimeout)
	if aerr != nil {
		return nil, aerr
	}
	return jwt, nil
}

// RequestDecode tries to retrieve a token from the cache by
// the requests authorization header. Otherwise it decodes it and
// puts it.
func (c *Cache) RequestDecode(req *http.Request) (*token.JWT, error) {
	var jwt *token.JWT
	var err error
	aerr := c.doSync(func() {
		var st string
		if st, err = c.requestToken(req); err != nil {
			return
		}
		if jwt, err = c.Get(st); err != nil {
			return
		}
		if jwt, err = token.Decode(st); err != nil {
			return
		}
		_, err = c.Put(jwt)
	}, defaultTimeout)
	if aerr != nil {
		return nil, aerr
	}
	return jwt, err
}

// RequestVerify tries to retrieve a token from the cache by
// the requests authorization header. Otherwise it verifies it and
// puts it.
func (c *Cache) RequestVerify(req *http.Request, key token.Key) (*token.JWT, error) {
	var jwt *token.JWT
	var err error
	aerr := c.doSync(func() {
		var st string
		if st, err = c.requestToken(req); err != nil {
			return
		}
		if jwt, err = c.Get(st); err != nil {
			return
		}
		if jwt, err = token.Verify(st, key); err != nil {
			return
		}
		_, err = c.Put(jwt)
	}, defaultTimeout)
	if aerr != nil {
		return nil, aerr
	}
	return jwt, err
}

// Put adds a token to the cache and return the total number of entries.
func (c *Cache) Put(jwt *token.JWT) (int, error) {
	var l int
	err := c.doSync(func() {
		if c.entries == nil {
			l = 0
			return
		}
		if jwt.IsValid(c.leeway) {
			c.entries[jwt.String()] = &cacheEntry{jwt, time.Now()}
			lenEntries := len(c.entries)
			if lenEntries > c.maxEntries {
				ttl := int64(c.ttl) / int64(lenEntries) * int64(c.maxEntries)
				c.cleanup(time.Duration(ttl))
			}
		}
		l = len(c.entries)
	}, defaultTimeout)
	return l, err
}

// Cleanup manually tells the cache to cleanup.
func (c *Cache) Cleanup() error {
	return c.doSync(func() {
		if c.entries == nil {
			return
		}
		c.cleanup(c.ttl)
	}, defaultTimeout)
}

// requestToken retrieves an authentication token out of a request.
func (c *Cache) requestToken(req *http.Request) (string, error) {
	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return "", failure.New("request contains no authorization header")
	}
	fields := strings.Fields(authorization)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return "", failure.New("invalid authorization header: %q", authorization)
	}
	return fields[1], nil
}

// cleanup checks for invalid or unused tokens.
func (c *Cache) cleanup(ttl time.Duration) {
	valids := map[string]*cacheEntry{}
	now := time.Now()
	for token, entry := range c.entries {
		if entry.jwt.IsValid(c.leeway) {
			if entry.accessed.Add(ttl).After(now) {
				// Everything fine.
				valids[token] = entry
			}
		}
	}
	c.entries = valids
}

// doSync performs a function in the backend synchronously.
func (c *Cache) doSync(action func(), timeout time.Duration) error {
	donec := make(chan struct{})
	c.actionc <- func() {
		action()
		close(donec)
	}
	select {
	case <-donec:
		return nil
	case <-time.After(timeout):
		return failure.New("cache action timeout")
	}
}

// backend is the goroutine of the cache.
func (c *Cache) backend() {
	ticker := time.NewTicker(c.interval)
	for {
		select {
		case <-c.ctx.Done():
			c.entries = map[string]*cacheEntry{}
			ticker.Stop()
			return
		case action := <-c.actionc:
			action()
		case <-ticker.C:
			go func() {
				if err := c.Cleanup(); err != nil {
					logger.Errorf("JWT cache: %v", err)
				}
			}()
		}
	}
}

// EOF
