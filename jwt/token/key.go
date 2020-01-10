// Tideland Go Network - JSON Web Token - Crypto
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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"

	"tideland.dev/go/trace/failure"
)

//--------------------
// KEY
//--------------------

// Key is the used key to sign a token. The real implementation
// controls signing and verification.
type Key interface{}

// ReadECPrivateKey reads a PEM formated ECDSA private key
// from the passed reader.
func ReadECPrivateKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, failure.New("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, failure.New("cannot decode the PEM")
	}
	var parsed *ecdsa.PrivateKey
	if parsed, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		return nil, failure.Annotate(err, "cannot parse the ECDSA")
	}
	return parsed, nil
}

// ReadECPublicKey reads a PEM encoded ECDSA public key
// from the passed reader.
func ReadECPublicKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, failure.New("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, failure.New("cannot decode the PEM")
	}
	var parsed interface{}
	parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, failure.Annotate(err, "cannot parse the ECDSA")
		}
		parsed = certificate.PublicKey
	}
	publicKey, ok := parsed.(*ecdsa.PublicKey)
	if !ok {
		return nil, failure.New("passed key is no ECDSA key")
	}
	return publicKey, nil
}

// ReadRSAPrivateKey reads a PEM encoded PKCS1 or PKCS8 private key
// from the passed reader.
func ReadRSAPrivateKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, failure.New("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, failure.New("cannot decode the PEM")
	}
	var parsed interface{}
	parsed, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, failure.Annotate(err, "cannot parse the RSA")
		}
	}
	privateKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, failure.New("passed key is no RSA key")
	}
	return privateKey, nil
}

// ReadRSAPublicKey reads a PEM encoded PKCS1 or PKCS8 public key
// from the passed reader.
func ReadRSAPublicKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, failure.New("cannot read the PEM")
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, failure.New("cannot decode the PEM")
	}
	var parsed interface{}
	parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, failure.Annotate(err, "cannot parse the RSA")
		}
		parsed = certificate.PublicKey
	}
	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, failure.New("passed key is no RSA key")
	}
	return publicKey, nil
}

// EOF
