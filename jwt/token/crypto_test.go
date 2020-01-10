// Tideland Go Network - JSON Web Token - Unit Tests
//
// Copyright (C) 2016-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token_test

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/net/jwt/token"
)

//--------------------
// TESTS
//--------------------

var (
	esTests = []token.Algorithm{token.ES256, token.ES384, token.ES512}
	hsTests = []token.Algorithm{token.HS256, token.HS384, token.HS512}
	psTests = []token.Algorithm{token.PS256, token.PS384, token.PS512}
	rsTests = []token.Algorithm{token.RS256, token.RS384, token.RS512}
	data    = []byte("the quick brown fox jumps over the lazy dog")
)

// TestESAlgorithms tests the ECDSA algorithms.
func TestESAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	for _, algo := range esTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		assert.Nil(err)
	}
}

// TestHSAlgorithms tests the HMAC algorithms.
func TestHSAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	key := []byte("secret")
	for _, algo := range hsTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, key)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, key)
		assert.Nil(err)
	}
}

// TestPSAlgorithms tests the RSAPSS algorithms.
func TestPSAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	for _, algo := range psTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		assert.Nil(err)
	}
}

// TestRSAlgorithms tests the RSA algorithms.
func TestRSAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	for _, algo := range rsTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		assert.Nil(err)
	}
}

// TestNoneAlgorithm tests the none algorithm.
func TestNoneAlgorithm(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing algorithm \"none\"")
	// Sign.
	signature, err := token.NONE.Sign(data, "")
	assert.Nil(err)
	assert.Empty(signature)
	// Verify.
	err = token.NONE.Verify(data, signature, "")
	assert.Nil(err)
}

// TestNotMatchingAlgorithm checks when algorithms of
// signing and verifying don't match.'
func TestNotMatchingAlgorithm(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	esPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	esPublicKey := esPrivateKey.Public()
	assert.Nil(err)
	hsKey := []byte("secret")
	rsPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	rsPublicKey := rsPrivateKey.Public()
	assert.Nil(err)
	noneKey := ""
	errorMatch := ".* combination of algorithm .* and key type .*"
	tests := []struct {
		description string
		algorithm   token.Algorithm
		key         token.Key
		signKeys    []token.Key
		verifyKeys  []token.Key
	}{
		{"ECDSA", token.ES512, esPrivateKey,
			[]token.Key{hsKey, rsPrivateKey, noneKey}, []token.Key{hsKey, rsPublicKey, noneKey}},
		{"HMAC", token.HS512, hsKey,
			[]token.Key{esPrivateKey, rsPrivateKey, noneKey}, []token.Key{esPublicKey, rsPublicKey, noneKey}},
		{"RSA", token.RS512, rsPrivateKey,
			[]token.Key{esPrivateKey, hsKey, noneKey}, []token.Key{esPublicKey, hsKey, noneKey}},
		{"RSAPSS", token.PS512, rsPrivateKey,
			[]token.Key{esPrivateKey, hsKey, noneKey}, []token.Key{esPublicKey, hsKey, noneKey}},
		{"none", token.NONE, noneKey,
			[]token.Key{esPrivateKey, hsKey, rsPrivateKey}, []token.Key{esPublicKey, hsKey, rsPublicKey}},
	}
	// Run the tests.
	for _, test := range tests {
		assert.Logf("testing %q algorithm key type mismatch", test.description)
		for _, key := range test.signKeys {
			_, err := test.algorithm.Sign(data, key)
			assert.ErrorMatch(err, errorMatch)
		}
		signature, err := test.algorithm.Sign(data, test.key)
		assert.Nil(err)
		for _, key := range test.verifyKeys {
			err = test.algorithm.Verify(data, signature, key)
			assert.ErrorMatch(err, errorMatch)
		}
	}
}

// TestESTools tests the tools for the reading of PEM encoded
// ECDSA keys.
func TestESTools(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing \"ECDSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	privateBytes, err := x509.MarshalECPrivateKey(privateKeyIn)
	assert.Nil(err)
	privateBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := token.ReadECPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := token.ReadECPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	signature, err := token.ES512.Sign(data, privateKeyOut)
	assert.Nil(err)
	err = token.ES512.Verify(data, signature, publicKeyOut)
	assert.Nil(err)
}

// TestRSTools tests the tools for the reading of PEM encoded
// RSA keys.
func TestRSTools(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing \"RSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKeyIn)
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := token.ReadRSAPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := token.ReadRSAPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	signature, err := token.RS512.Sign(data, privateKeyOut)
	assert.Nil(err)
	err = token.RS512.Verify(data, signature, publicKeyOut)
	assert.Nil(err)
}

// EOF
