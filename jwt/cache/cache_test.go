// Tideland Go Network - JSON Web Token - Cache - Unit Tests
//
// Copyright (C) 2016-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cache_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/net/jwt/cache"
	"tideland.dev/go/net/jwt/token"
)

//--------------------
// TESTS
//--------------------

// TestCachePutGet tests the putting and getting of tokens
// to the cache.
func TestCachePutGet(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing cache put and get")
	ctx := context.Background()
	cache := cache.New(ctx, time.Minute, time.Minute, time.Minute, 10)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := token.Encode(claims, key, token.HS512)
	assert.Nil(err)
	_, err = cache.Put(jwtIn)
	assert.NoError(err)
	token := jwtIn.String()
	jwtOut, err := cache.Get(token)
	assert.NoError(err)
	assert.Equal(jwtIn, jwtOut)
	jwtOut, err = cache.Get("is.not.there")
	assert.NoError(err)
	assert.Nil(jwtOut)
}

// TestCacheAccessCleanup tests the access based cleanup
// of the JWT cache.
func TestCacheAccessCleanup(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing cache access based cleanup")
	ctx := context.Background()
	cache := cache.New(ctx, time.Second, time.Second, time.Second, 10)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := token.Encode(claims, key, token.HS512)
	assert.NoError(err)
	_, err = cache.Put(jwtIn)
	assert.NoError(err)
	token := jwtIn.String()
	jwtOut, err := cache.Get(token)
	assert.NoError(err)
	assert.Equal(jwtIn, jwtOut)
	// Now wait a bit an try again.
	time.Sleep(5 * time.Second)
	jwtOut, err = cache.Get(token)
	assert.NoError(err)
	assert.Nil(jwtOut)
}

// TestCacheValidityCleanup tests the validity based cleanup
// of the JWT cache.
func TestCacheValidityCleanup(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing cache validity based cleanup")
	ctx := context.Background()
	cache := cache.New(ctx, time.Minute, time.Second, time.Second, 10)
	key := []byte("secret")
	now := time.Now()
	nbf := now.Add(-2 * time.Second)
	exp := now.Add(2 * time.Second)
	claims := initClaims()
	claims.SetNotBefore(nbf)
	claims.SetExpiration(exp)
	jwtIn, err := token.Encode(claims, key, token.HS512)
	assert.Nil(err)
	_, err = cache.Put(jwtIn)
	assert.NoError(err)
	token := jwtIn.String()
	jwtOut, err := cache.Get(token)
	assert.NoError(err)
	assert.Equal(jwtOut, jwtIn)
	// Now access until it is invalid and not
	// available anymore.
	var i int
	for i = 0; i < 5; i++ {
		time.Sleep(time.Second)
		jwtOut, err = cache.Get(token)
		assert.NoError(err)
		if jwtOut == nil {
			break
		}
		assert.Equal(jwtOut, jwtIn)
	}
	assert.True(i > 1 && i < 4)
}

// TestCacheLoad tests the cache load based cleanup.
func TestCacheLoad(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing cache load based cleanup")
	cacheTime := 100 * time.Millisecond
	ctx := context.Background()
	cache := cache.New(ctx, 2*cacheTime, cacheTime, cacheTime, 4)
	claims := initClaims()
	// Now fill the cache and check that it doesn't
	// grow too high.
	var i int
	for i = 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		key := []byte(fmt.Sprintf("secret-%d", i))
		jwtIn, err := token.Encode(claims, key, token.HS512)
		assert.Nil(err)
		size, err := cache.Put(jwtIn)
		assert.NoError(err)
		assert.True(size < 6)
	}
}

// TestCacheContext tests the cache stopping by context.
func TestCacheContext(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing cache stopping by context")
	ctx, cancel := context.WithCancel(context.Background())
	cache := cache.New(ctx, time.Minute, time.Minute, time.Minute, 10)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := token.Encode(claims, key, token.HS512)
	assert.NoError(err)
	_, err = cache.Put(jwtIn)
	assert.NoError(err)
	// Now cancel and test to get token.
	cancel()
	time.Sleep(10 * time.Millisecond)
	token := jwtIn.String()
	jwtOut, err := cache.Get(token)
	assert.ErrorMatch(err, ".* cache action timeout")
	assert.Nil(jwtOut)
}

//--------------------
// HELPERS
//--------------------

// initClaims creates test claims.
func initClaims() token.Claims {
	c := token.NewClaims()
	c.SetSubject("1234567890")
	c.Set("name", "John Doe")
	c.Set("admin", true)
	return c
}

// EOF
