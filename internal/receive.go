/*
 * MIT License
 *
 * Copyright (c) 2021 TECHCRAFT TECHNOLOGIES CO LTD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/techcraftt/base/internal/io"
	stdio "io"
	"net/http"
	"net/http/httputil"
	"strings"
)

var (
	_ Receiver = (*receiver)(nil)
)

type (
	receiver struct {
		Logger    stdio.Writer
		DebugMode bool
	}

	Receiver interface {
		Receive(ctx context.Context, rn string, r *http.Request, v interface{}) (*Receipt, error)
	}

	BasicAuth struct {
		Username string
		Password string
	}

	Receipt struct {
		Request *http.Request
		BearerToken string
		BasicAuth BasicAuth
		ApiKey string
		RemoteAddress string
		ForwardedFor string
	}

	ReceiveParams struct {
		DebugMode bool
		Logger    stdio.Writer
	}

	ReceiveOption func(params *ReceiveParams)
)

func NewReceiver(writer stdio.Writer, debug bool) Receiver {
	return &receiver{
		Logger:    writer,
		DebugMode: debug,
	}
}

func (rc *receiver) Receive(ctx context.Context, rn string, r *http.Request, v interface{}) (*Receipt, error) {
	receipt := new(Receipt)


	// capture bearer token
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) >= 2 {
		reqToken = strings.TrimSpace(splitToken[1])
		receipt.BearerToken = reqToken
	}

	// capture basic auth
	username, password, ok := r.BasicAuth()
	if ok{
		ba := BasicAuth{
			Username: username,
			Password: password,
		}
		receipt.BasicAuth = ba
	}

	receipt.RemoteAddress = r.RemoteAddr
	receipt.ForwardedFor = r.Header.Get("X-Forwarded-For")
	receipt.ApiKey = r.Header.Get("X-Api-key")

	rClone := r.Clone(ctx)
	receipt.Request = rClone
	contentType := r.Header.Get("Content-Type")
	payloadType := categorizeContentType(contentType)
	body, err := stdio.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// restore request body
	r.Body = stdio.NopCloser(bytes.NewBuffer(body))

	defer func(debug bool) {
		if debug {
			rc.logRequest(rn, r)
		}
	}(rc.DebugMode)

	if v == nil {
		return receipt, nil
	}

	switch payloadType {
	case JsonPayload:
		err := json.NewDecoder(r.Body).Decode(v)
		defer func(Body stdio.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)
		r.Body = stdio.NopCloser(bytes.NewBuffer(body))
		return receipt, err

	case XmlPayload:
		r.Body = stdio.NopCloser(bytes.NewBuffer(body))
		return receipt, xml.Unmarshal(body, v)
	}

	return receipt, err
}

func ReceiveDebugMode(mode bool) ReceiveOption {
	return func(params *ReceiveParams) {
		params.DebugMode = mode
	}
}

func ReceiveLogger(writer stdio.Writer) ReceiveOption {
	return func(params *ReceiveParams) {
		params.Logger = writer
	}
}

// ReceivePayloadWithParams takes *http.Request from clients like during then unmarshal the provided
// request into given interface v
// The expected Content-Type should also be declared. If its cTypeJson or
// application/xml
func ReceivePayloadWithParams(r *http.Request, v interface{}, opts ...ReceiveOption) error {

	rp := &ReceiveParams{
		DebugMode: true,
		Logger:    io.Stderr,
	}

	for _, opt := range opts {
		opt(rp)
	}
	contentType := r.Header.Get("Content-Type")
	payloadType := categorizeContentType(contentType)
	body, err := stdio.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if v == nil {
		return fmt.Errorf("v can not be nil")
	}
	// restore request body
	r.Body = stdio.NopCloser(bytes.NewBuffer(body))

	switch payloadType {
	case JsonPayload:
		err := json.NewDecoder(r.Body).Decode(v)
		defer func(Body stdio.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)
		r.Body = stdio.NopCloser(bytes.NewBuffer(body))
		return err

	case XmlPayload:
		r.Body = stdio.NopCloser(bytes.NewBuffer(body))
		return xml.Unmarshal(body, v)
	}

	return err
}

// ReceivePayload takes *http.Request from clients like during then unmarshal the provided
// request into given interface v
// The expected Content-Type should also be declared. If its cTypeJson or
// application/xml
func ReceivePayload(r *http.Request, v interface{}) error {

	contentType := r.Header.Get("Content-Type")
	payloadType := categorizeContentType(contentType)
	body, err := stdio.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if v == nil {
		return fmt.Errorf("v can not be nil")
	}
	// restore request body
	r.Body = stdio.NopCloser(bytes.NewBuffer(body))

	switch payloadType {
	case JsonPayload:
		err := json.NewDecoder(r.Body).Decode(v)
		defer func(Body stdio.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)
		r.Body = stdio.NopCloser(bytes.NewBuffer(body))
		return err

	case XmlPayload:
		r.Body = stdio.NopCloser(bytes.NewBuffer(body))
		return xml.Unmarshal(body, v)
	}

	return err
}

// logRequest is called to print the details of http.Request received
func (rc *receiver) logRequest(name string, request *http.Request) {

	if request != nil && rc.DebugMode {
		reqDump, _ := httputil.DumpRequest(request, true)
		_, err := fmt.Fprintf(rc.Logger, "%s REQUEST: %s\n", name, reqDump)
		if err != nil {
			fmt.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
			return
		}
		return
	}
	return
}
