// Tideland Go Network - JSON Web Token
//
// Copyright (C) 2016-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token // import "tideland.dev/go/net/jwt/token"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"tideland.dev/go/trace/failure"
)

//--------------------
// JSON Web Token
//--------------------

// jwtHeader contains the JWT header fields.
type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// JWT manages the parts of a JSON Web Token and the access to those.
type JWT struct {
	claims    Claims
	key       Key
	algorithm Algorithm
	token     string
}

// Encode creates a JSON Web Token for the given claims
// based on key and algorithm.
func Encode(claims Claims, key Key, algorithm Algorithm) (*JWT, error) {
	jwt := &JWT{
		claims:    claims,
		key:       key,
		algorithm: algorithm,
	}
	headerPart, err := marshallAndEncode(jwtHeader{string(algorithm), "JWT"})
	if err != nil {
		return nil, failure.Annotate(err, "cannot encode the header")
	}
	claimsPart, err := marshallAndEncode(claims)
	if err != nil {
		return nil, failure.Annotate(err, "cannot encode the claims")
	}
	dataParts := headerPart + "." + claimsPart
	signaturePart, err := signAndEncode([]byte(dataParts), key, algorithm)
	if err != nil {
		return nil, failure.Annotate(err, " cannot encode the signature")
	}
	jwt.token = dataParts + "." + signaturePart
	return jwt, nil
}

// Decode creates a token out of a string without verification.
func Decode(token string) (*JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, failure.New("cannot decode the parts")
	}
	var header jwtHeader
	err := decodeAndUnmarshall(parts[0], &header)
	if err != nil {
		return nil, failure.Annotate(err, "cannot decode the header")
	}
	var claims Claims
	err = decodeAndUnmarshall(parts[1], &claims)
	if err != nil {
		return nil, failure.Annotate(err, "cannot decode the claims")
	}
	return &JWT{
		claims:    claims,
		algorithm: Algorithm(header.Algorithm),
		token:     token,
	}, nil
}

// Verify creates a token out of a string and varifies it against
// the passed key.
func Verify(token string, key Key) (*JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, failure.New("cannot verify the parts")
	}
	var header jwtHeader
	err := decodeAndUnmarshall(parts[0], &header)
	if err != nil {
		return nil, failure.Annotate(err, "cannot verify the header")
	}
	err = decodeAndVerify(parts, key, Algorithm(header.Algorithm))
	if err != nil {
		return nil, failure.Annotate(err, "cannot verify the signature")
	}
	var claims Claims
	err = decodeAndUnmarshall(parts[1], &claims)
	if err != nil {
		return nil, failure.Annotate(err, "cannot verify the claims")
	}
	return &JWT{
		claims:    claims,
		key:       key,
		algorithm: Algorithm(header.Algorithm),
		token:     token,
	}, nil
}

// Claims returns the claims payload of the token.
func (jwt *JWT) Claims() Claims {
	return jwt.claims
}

// Key returns the key of the token only when it is a result of encoding or verification.
func (jwt *JWT) Key() (Key, error) {
	if jwt.key == nil {
		return nil, failure.New("no key available, only after encoding or verifying")
	}
	return jwt.key, nil
}

// Algorithm returns the algorithm of the token after encoding, decoding, or verification.
func (jwt *JWT) Algorithm() Algorithm {
	return jwt.algorithm
}

// IsValid is a convenience method checking the registered claims if the token is valid.
func (jwt *JWT) IsValid(leeway time.Duration) bool {
	return jwt.claims.IsValid(leeway)
}

// String implements the fmt.Stringer interface.
func (jwt *JWT) String() string {
	return jwt.token
}

//--------------------
// PRIVATE HELPERS
//--------------------

// marshallAndEncode marshals the passed value to JSON and
// creates a BASE64 string out of it.
func marshallAndEncode(value interface{}) (string, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", failure.Annotate(err, "error marshalling to JSON")
	}
	encoded := base64.RawURLEncoding.EncodeToString(jsonValue)
	return encoded, nil
}

// decodeAndUnmarshall decodes a BASE64 encoded JSON string and
// unmarshals it into the passed value.
func decodeAndUnmarshall(part string, value interface{}) error {
	decoded, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return failure.Annotate(err, "part of the token contains invalid data")
	}
	err = json.Unmarshal(decoded, value)
	if err != nil {
		return failure.Annotate(err, "error unmarshalling from JSON")
	}
	return nil
}

// signAndEncode creates the signature for the data part (header and
// payload) of the token using the passed key and algorithm. The result
// is then encoded to BASE64.
func signAndEncode(data []byte, key Key, algorithm Algorithm) (string, error) {
	sig, err := algorithm.Sign(data, key)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(sig)
	return encoded, nil
}

// decodeAndVerify decodes a BASE64 encoded signature and verifies
// the correct signing of the data part (header and payload) using the
// passed key and algorithm.
func decodeAndVerify(parts []string, key Key, algorithm Algorithm) error {
	data := []byte(parts[0] + "." + parts[1])
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return failure.Annotate(err, "part of the token contains invalid data")
	}
	return algorithm.Verify(data, sig, key)
}

// EOF
