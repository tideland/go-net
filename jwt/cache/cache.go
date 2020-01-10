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
	"sync"
	"time"

	"tideland.dev/go/net/jwt/token"
	"tideland.dev/go/together/loop"
	"tideland.dev/go/together/notifier"
	"tideland.dev/go/trace/failure"
)

//--------------------
// CACHE
//--------------------

// cacheEntry manages a token and its access time.
type cacheEntry struct {
	jwt      *token.JWT
	accessed time.Time
}

// Cache provides a caching for tokens so that these
// don't have to be decoded or verified multiple times.
type Cache struct {
	mu         sync.Mutex
	entries    map[string]*cacheEntry
	ttl        time.Duration
	leeway     time.Duration
	interval   time.Duration
	maxEntries int
	cleanupc   chan time.Duration
	loop       *loop.Loop
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
		entries:    map[string]*cacheEntry{},
		ttl:        ttl,
		leeway:     leeway,
		interval:   interval,
		maxEntries: maxEntries,
		cleanupc:   make(chan time.Duration, 5),
	}
	options := []loop.Option{loop.WithFinalizer(c.finalize)}
	if ctx != nil {
		options = append(options, loop.WithContext(ctx))
	}
	l, err := loop.Go(c.worker, options...)
	if err != nil {
		panic("new JWT cache: " + err.Error())
	}
	c.loop = l
	return c
}

// Get tries to retrieve a token from the cache.
func (c *Cache) Get(token string) (*token.JWT, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		return nil, false
	}
	entry, ok := c.entries[token]
	if !ok {
		return nil, false
	}
	if entry.jwt.IsValid(c.leeway) {
		entry.accessed = time.Now()
		return entry.jwt, true
	}
	// Remove invalid token.
	delete(c.entries, token)
	return nil, false
}

// RequestDecode tries to retrieve a token from the cache by
// the requests authorization header. Otherwise it decodes it and
// puts it.
func (c *Cache) RequestDecode(req *http.Request) (*token.JWT, error) {
	t, err := c.requestToken(req)
	if err != nil {
		return nil, err
	}
	jwt, ok := c.Get(t)
	if ok {
		return jwt, nil
	}
	jwt, err = token.Decode(t)
	if err != nil {
		return nil, err
	}
	c.Put(jwt)
	return jwt, nil
}

// RequestVerify tries to retrieve a token from the cache by
// the requests authorization header. Otherwise it verifies it and
// puts it.
func (c *Cache) RequestVerify(req *http.Request, key token.Key) (*token.JWT, error) {
	t, err := c.requestToken(req)
	if err != nil {
		return nil, err
	}
	jwt, ok := c.Get(t)
	if ok {
		return jwt, nil
	}
	jwt, err = token.Verify(t, key)
	if err != nil {
		return nil, err
	}
	c.Put(jwt)
	return jwt, nil
}

// Put adds a token to the cache and return the total number of entries.
func (c *Cache) Put(jwt *token.JWT) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		return 0
	}
	if jwt.IsValid(c.leeway) {
		c.entries[jwt.String()] = &cacheEntry{jwt, time.Now()}
		lenEntries := len(c.entries)
		if lenEntries > c.maxEntries {
			ttl := int64(c.ttl) / int64(lenEntries) * int64(c.maxEntries)
			c.cleanupc <- time.Duration(ttl)
		}
	}
	return len(c.entries)
}

// Cleanup manually tells the cache to cleanup.
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		return
	}
	c.cleanupc <- c.ttl
}

// Stop tells the cache to end working.
func (c *Cache) Stop() error {
	return c.loop.Stop(nil)
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

// worker runs a cleaning session every five minutes.
func (c *Cache) worker(nc *notifier.Closer) error {
	defer func() {
		// Cleanup entries map after stop or error.
		c.mu.Lock()
		defer c.mu.Unlock()
		c.entries = nil
	}()
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-nc.Done():
			return nil
		case ttl := <-c.cleanupc:
			c.cleanup(ttl)
		case <-ticker.C:
			c.cleanup(c.ttl)
		}
	}
}

// finalizer cleans the data when the loop ends.
func (c *Cache) finalize(err error) error {
	c.entries = map[string]*cacheEntry{}
	return err
}

// cleanup checks for invalid or unused tokens.
func (c *Cache) cleanup(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
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

// EOF
