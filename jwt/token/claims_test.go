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
	"encoding/json"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/net/jwt/token"
)

//--------------------
// TESTS
//--------------------

// TestClaimsMarshalling tests the marshalling of Claims to JSON and back.
func TestClaimsMarshalling(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims marshalling")
	// First with uninitialised or empty token.
	var c token.Claims
	jsonValue, err := json.Marshal(c)
	assert.Equal(string(jsonValue), "{}")
	assert.Nil(err)
	c = token.NewClaims()
	jsonValue, err = json.Marshal(c)
	assert.Equal(string(jsonValue), "{}")
	assert.Nil(err)
	// Now fill it.
	c.Set("foo", "yadda")
	c.Set("bar", 12345)
	assert.Length(c, 2)
	jsonValue, err = json.Marshal(c)
	assert.NotNil(jsonValue)
	assert.Nil(err)
	var uc token.Claims
	err = json.Unmarshal(jsonValue, &uc)
	assert.Nil(err)
	assert.Length(uc, 2)
	foo, ok := uc.Get("foo")
	assert.Equal(foo, "yadda")
	assert.True(ok)
	bar, ok := uc.GetInt("bar")
	assert.Equal(bar, 12345)
	assert.True(ok)
}

// TestClaimsBasic tests the low level operations on token.
func TestClaimsBasic(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims basic functions handling")
	// First with uninitialised token.
	var c token.Claims
	ok := c.Contains("foo")
	assert.False(ok)
	nothing, ok := c.Get("foo")
	assert.Nil(nothing)
	assert.False(ok)
	old := c.Set("foo", "bar")
	assert.Nil(old)
	old = c.Delete("foo")
	assert.Nil(old)
	// Now initialise it.
	c = token.NewClaims()
	ok = c.Contains("foo")
	assert.False(ok)
	nothing, ok = c.Get("foo")
	assert.Nil(nothing)
	assert.False(ok)
	old = c.Set("foo", "bar")
	assert.Nil(old)
	ok = c.Contains("foo")
	assert.True(ok)
	foo, ok := c.Get("foo")
	assert.Equal(foo, "bar")
	assert.True(ok)
	old = c.Set("foo", "yadda")
	assert.Equal(old, "bar")
	// Finally delete it.
	old = c.Delete("foo")
	assert.Equal(old, "yadda")
	old = c.Delete("foo")
	assert.Nil(old)
	ok = c.Contains("foo")
	assert.False(ok)
}

// TestClaimsString tests the string operations on token.
func TestClaimsString(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims string handling")
	c := token.NewClaims()
	nothing := c.Set("foo", "bar")
	assert.Nil(nothing)
	var foo string
	foo, ok := c.GetString("foo")
	assert.Equal(foo, "bar")
	assert.True(ok)
	c.Set("foo", 4711)
	foo, ok = c.GetString("foo")
	assert.Equal(foo, "4711")
	assert.True(ok)
}

// TestClaimsBool tests the bool operations on token.
func TestClaimsBool(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims bool handling")
	c := token.NewClaims()
	c.Set("foo", true)
	c.Set("bar", false)
	c.Set("baz", "T")
	c.Set("bingo", "0")
	c.Set("yadda", "nope")
	foo, ok := c.GetBool("foo")
	assert.True(foo)
	assert.True(ok)
	bar, ok := c.GetBool("bar")
	assert.False(bar)
	assert.True(ok)
	baz, ok := c.GetBool("baz")
	assert.True(baz)
	assert.True(ok)
	bingo, ok := c.GetBool("bingo")
	assert.False(bingo)
	assert.True(ok)
	yadda, ok := c.GetBool("yadda")
	assert.False(yadda)
	assert.False(ok)
}

// TestClaimsInt tests the int operations on token.
func TestClaimsInt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims int handling")
	c := token.NewClaims()
	c.Set("foo", 4711)
	c.Set("bar", "4712")
	c.Set("baz", 4713.0)
	c.Set("yadda", "nope")
	foo, ok := c.GetInt("foo")
	assert.Equal(foo, 4711)
	assert.True(ok)
	bar, ok := c.GetInt("bar")
	assert.Equal(bar, 4712)
	assert.True(ok)
	baz, ok := c.GetInt("baz")
	assert.Equal(baz, 4713)
	assert.True(ok)
	yadda, ok := c.GetInt("yadda")
	assert.Equal(yadda, 0)
	assert.False(ok)
}

// TestClaimsFloat64 tests the float64 operations on token.
func TestClaimsFloat64(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims float64 handling")
	c := token.NewClaims()
	c.Set("foo", 4711)
	c.Set("bar", "4712")
	c.Set("baz", 4713.0)
	c.Set("yadda", "nope")
	foo, ok := c.GetFloat64("foo")
	assert.Equal(foo, 4711.0)
	assert.True(ok)
	bar, ok := c.GetFloat64("bar")
	assert.Equal(bar, 4712.0)
	assert.True(ok)
	baz, ok := c.GetFloat64("baz")
	assert.Equal(baz, 4713.0)
	assert.True(ok)
	yadda, ok := c.GetFloat64("yadda")
	assert.Equal(yadda, 0.0)
	assert.False(ok)
}

// TestClaimsTime tests the time operations on token.
func TestClaimsTime(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims time handling")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := token.NewClaims()
	c.SetTime("foo", goLaunch)
	c.Set("bar", goLaunch.Unix())
	c.Set("baz", goLaunch.Format(time.RFC3339))
	c.Set("yadda", "nope")
	foo, ok := c.GetTime("foo")
	assert.Equal(foo.Unix(), goLaunch.Unix())
	assert.True(ok)
	bar, ok := c.GetTime("bar")
	assert.Equal(bar.Unix(), goLaunch.Unix())
	assert.True(ok)
	baz, ok := c.GetTime("baz")
	assert.Equal(baz.Unix(), goLaunch.Unix())
	assert.True(ok)
	yadda, ok := c.GetTime("yadda")
	assert.Equal(yadda, time.Time{})
	assert.False(ok)
}

// nestedValue is used as a structured value of a claim.
type nestedValue struct {
	Name  string
	Value int
}

// TestClaimsMarshalledValue tests the marshalling and
// unmarshalling of structures as values.
func TestClaimsMarshalledValue(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims deep value unmarshalling")
	baz := []*nestedValue{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}
	c := token.NewClaims()
	c.Set("foo", "bar")
	c.Set("baz", baz)
	// Now marshal and unmarshal the token.
	jsonValue, err := json.Marshal(c)
	assert.NotNil(jsonValue)
	assert.Nil(err)
	var uc token.Claims
	err = json.Unmarshal(jsonValue, &uc)
	assert.Nil(err)
	assert.Length(uc, 2)
	foo, ok := uc.Get("foo")
	assert.Equal(foo, "bar")
	assert.True(ok)
	var ubaz []*nestedValue
	ok, err = uc.GetMarshalled("baz", &ubaz)
	assert.True(ok)
	assert.Nil(err)
	assert.Length(ubaz, 3)
	assert.Equal(ubaz[0].Name, "one")
	assert.Equal(ubaz[2].Value, 3)
}

// TestClaimsAudience checks the setting, getting, and
// deleting of the audience claim.
func TestClaimsAudience(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claim \"aud\"")
	audience := []string{"foo", "bar", "baz"}
	c := token.NewClaims()
	aud, ok := c.Audience()
	assert.False(ok)
	assert.Nil(aud)
	none := c.SetAudience(audience...)
	assert.Length(none, 0)
	aud, ok = c.Audience()
	assert.Equal(aud, audience)
	assert.True(ok)
	old := c.DeleteAudience()
	assert.Equal(old, aud)
	_, ok = c.Audience()
	assert.False(ok)
}

// TestClaimsExpiration checks the setting, getting, and
// deleting of the expiration claim.
func TestClaimsExpiration(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claim \"exp\"")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := token.NewClaims()
	exp, ok := c.Expiration()
	assert.False(ok)
	none := c.SetExpiration(goLaunch)
	assert.Equal(none, time.Time{})
	exp, ok = c.Expiration()
	assert.Equal(exp.Unix(), goLaunch.Unix())
	assert.True(ok)
	old := c.DeleteExpiration()
	assert.Equal(old.Unix(), exp.Unix())
	exp, ok = c.Expiration()
	assert.False(ok)
}

// TestClaimsIdentifier checks the setting, getting, and
// deleting of the identifier claim.
func TestClaimsIdentifier(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claim \"jti\"")
	identifier := "foo"
	c := token.NewClaims()
	jti, ok := c.Identifier()
	assert.False(ok)
	assert.Empty(jti)
	none := c.SetIdentifier(identifier)
	assert.Equal(none, "")
	jti, ok = c.Identifier()
	assert.Equal(jti, identifier)
	assert.True(ok)
	old := c.DeleteIdentifier()
	assert.Equal(old, jti)
	_, ok = c.Identifier()
	assert.False(ok)
}

// TestClaimsIssuedAt checks the setting, getting, and
// deleting of the issued at claim.
func TestClaimsIssuedAt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claim \"iat\"")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := token.NewClaims()
	iat, ok := c.IssuedAt()
	assert.False(ok)
	none := c.SetIssuedAt(goLaunch)
	assert.Equal(none, time.Time{})
	iat, ok = c.IssuedAt()
	assert.Equal(iat.Unix(), goLaunch.Unix())
	assert.True(ok)
	old := c.DeleteIssuedAt()
	assert.Equal(old.Unix(), iat.Unix())
	iat, ok = c.IssuedAt()
	assert.False(ok)
}

// TestClaimsIssuer checks the setting, getting, and
// deleting of the issuer claim.
func TestClaimsIssuer(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claim \"iss\"")
	issuer := "foo"
	c := token.NewClaims()
	iss, ok := c.Issuer()
	assert.False(ok)
	assert.Empty(iss)
	none := c.SetIssuer(issuer)
	assert.Equal(none, "")
	iss, ok = c.Issuer()
	assert.Equal(iss, issuer)
	assert.True(ok)
	old := c.DeleteIssuer()
	assert.Equal(old, iss)
	_, ok = c.Issuer()
	assert.False(ok)
}

// TestClaimsNotBefore checks the setting, getting, and
// deleting of the not before claim.
func TestClaimsNotBefore(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claim \"nbf\"")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c := token.NewClaims()
	nbf, ok := c.NotBefore()
	assert.False(ok)
	none := c.SetNotBefore(goLaunch)
	assert.Equal(none, time.Time{})
	nbf, ok = c.NotBefore()
	assert.Equal(nbf.Unix(), goLaunch.Unix())
	assert.True(ok)
	old := c.DeleteNotBefore()
	assert.Equal(old.Unix(), nbf.Unix())
	_, ok = c.NotBefore()
	assert.False(ok)
}

// TestClaimsSubject checks the setting, getting, and
// deleting of the subject claim.
func TestClaimsSubject(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claim \"sub\"")
	subject := "foo"
	c := token.NewClaims()
	sub, ok := c.Subject()
	assert.False(ok)
	assert.Empty(sub)
	none := c.SetSubject(subject)
	assert.Equal(none, "")
	sub, ok = c.Subject()
	assert.Equal(sub, subject)
	assert.True(ok)
	old := c.DeleteSubject()
	assert.Equal(old, sub)
	_, ok = c.Subject()
	assert.False(ok)
}

// TestClaimsValidity checks the validation of the not before
// and the expiring time.
func TestClaimsValidity(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	assert.Logf("testing claims validity")
	// Fresh token.
	now := time.Now()
	leeway := time.Minute
	c := token.NewClaims()
	valid := c.IsAlreadyValid(leeway)
	assert.True(valid)
	valid = c.IsStillValid(leeway)
	assert.True(valid)
	valid = c.IsValid(leeway)
	assert.True(valid)
	// Set times.
	nbf := now.Add(-time.Hour)
	exp := now.Add(time.Hour)
	c.SetNotBefore(nbf)
	valid = c.IsAlreadyValid(leeway)
	assert.True(valid)
	c.SetExpiration(exp)
	valid = c.IsStillValid(leeway)
	assert.True(valid)
	valid = c.IsValid(leeway)
	assert.True(valid)
	// Invalid token.
	nbf = now.Add(time.Hour)
	exp = now.Add(-time.Hour)
	c.SetNotBefore(nbf)
	c.DeleteExpiration()
	valid = c.IsAlreadyValid(leeway)
	assert.False(valid)
	valid = c.IsValid(leeway)
	assert.False(valid)
	c.DeleteNotBefore()
	c.SetExpiration(exp)
	valid = c.IsStillValid(leeway)
	assert.False(valid)
	valid = c.IsValid(leeway)
	assert.False(valid)
	c.SetNotBefore(nbf)
	valid = c.IsValid(leeway)
	assert.False(valid)
}

// EOF
