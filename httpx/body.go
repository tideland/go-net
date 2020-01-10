// Tideland Go Network - HTTP Extensions
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx // import "tideland.dev/go/net/httpx"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"tideland.dev/go/trace/failure"
)

//--------------------
// BODY TOOLS
//--------------------

// ReadBody retrieves the whole body out of a HTTP request or response
// and returns it as byte slice.
func ReadBody(r io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, failure.Annotate(err, "cannot read body")
	}
	r.Close()
	return data, nil
}

// WriteBody writes the whole body into a HTTP request or response.
func WriteBody(w io.Writer, data []byte) error {
	if _, err := w.Write(data); err != nil {
		return failure.Annotate(err, "cannot write body")
	}
	return nil
}

// UnmarshalBody parses the body data of a request or response based on the
// content type header stores the result in the value pointed by v. Conten types
// JSON and XML expect the according types as arguments, all text types
// expect a string, and all others too, but the data is encoded in BASE64.
func UnmarshalBody(r io.ReadCloser, h http.Header, v interface{}) error {
	// First read.
	data, err := ReadBody(r)
	if err != nil {
		return err
	}
	// Then unmarshal based on content type.
	switch {
	case ContainsContentType(h, contentTypesText):
		switch tv := v.(type) {
		case *string:
			*tv = string(data)
		case *interface{}:
			*tv = string(data)
		default:
			return failure.New("invalid value argument for text; need string or empty interface")
		}
		return nil
	case ContainsContentType(h, ContentTypeJSON):
		if err = json.Unmarshal(data, &v); err != nil {
			return failure.Annotate(err, "cannot unmarshal JSON")
		}
		return nil
	case ContainsContentType(h, ContentTypeXML):
		if err = xml.Unmarshal(data, &v); err != nil {
			return failure.Annotate(err, "cannot unmarshal XML")
		}
		return nil
	default:
		sd := base64.StdEncoding.EncodeToString(data)
		switch tv := v.(type) {
		case *string:
			*tv = sd
		case *interface{}:
			*tv = sd
		default:
			return failure.New("invalid value argument for base64; need string or empty interface")
		}
		return nil
	}
}

// MarshalBody allows to directly marshal a value into a writer depending on
// the content type.
func MarshalBody(w io.Writer, h http.Header, v interface{}) error {
	// First marshal based on content type.
	var data []byte
	var err error
	switch {
	case ContainsContentType(h, contentTypesText):
		data = []byte(fmt.Sprintf("%v", v))
	case ContainsContentType(h, ContentTypeJSON):
		data, err = json.Marshal(v)
		if err != nil {
			return failure.Annotate(err, "cannot marshal to JSON")
		}
	case ContainsContentType(h, ContentTypeXML):
		data, err = xml.Marshal(v)
		if err != nil {
			return failure.Annotate(err, "cannot marshal to XML")
		}
	default:
		vbs, ok := v.([]byte)
		if !ok {
			return failure.New("invalid value argument")
		}
		data = vbs
	}
	// Then write the body.
	return WriteBody(w, data)
}

// EOF
