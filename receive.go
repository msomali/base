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

package base

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	stdio "io"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

var (
	_ Receiver = (*receiver)(nil)
)

type (
	receiver struct {
		mu sync.Mutex
		Logger    stdio.Writer
		DebugMode bool
	}

	Receiver interface {
		Receive(ctx context.Context, rn string, r *http.Request, v interface{},opts...OptionFunc) (*Receipt, error)
	}

	BasicAuth struct {
		Username string
		Password string
	}

	Receipt struct {
		Request       *http.Request
		BearerToken   string
		BasicAuth     BasicAuth
		ApiKey        string
		RemoteAddress string
		ForwardedFor  string
	}
)

func NewReceiver(writer stdio.Writer, debug bool) Receiver {
	return &receiver{
		mu: sync.Mutex{},
		Logger:    writer,
		DebugMode: debug,
	}
}

func (rc *receiver)update(params *Params) {
    rc.mu.Lock()
    defer rc.mu.Unlock()
	if params != nil {
		rc.Logger = params.Logger
		rc.DebugMode = params.DebugMode
	}
}

func (rc *receiver) Receive(ctx context.Context, rn string, r *http.Request, v interface{},opts...OptionFunc) (*Receipt, error) {
	params := &Params{
        DebugMode: rc.DebugMode,
        Logger:    rc.Logger,
    }

	for _, opt := range opts {
		opt(params)
	}

	// update receiver in case the options changed
	rc.update(params)

	var (
		bodyBytes []byte
		err       error
	)

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
	if ok {
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
	if r.Body != nil {
		bodyBytes, err = stdio.ReadAll(r.Body)
	}

	if err != nil {
		return nil, err
	}

	// restore request body
	r.Body = stdio.NopCloser(bytes.NewBuffer(bodyBytes))

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
		r.Body = stdio.NopCloser(bytes.NewBuffer(bodyBytes))
		return receipt, err

	case XmlPayload:
		r.Body = stdio.NopCloser(bytes.NewBuffer(bodyBytes))
		return receipt, xml.Unmarshal(bodyBytes, v)
	}

	return receipt, err
}

// logRequest is called to print the details of http.Request received
func (rc *receiver) logRequest(name string, request *http.Request) {

	rn := strings.ToUpper(fmt.Sprintf("%s request (RECEIVED)", name))
	if request != nil && rc.DebugMode {
		reqDump, _ := httputil.DumpRequest(request, true)
		_, err := fmt.Fprintf(rc.Logger, "\n\n%s : %s\n\n", rn, reqDump)
		if err != nil {
			fmt.Printf("Error while logging %s request: %v\n",
				strings.ToLower(name), err)
			return
		}
		return
	}
	return
}


