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
	"context"
	"fmt"
	"net/http"
	"strings"
)

type (

	// Request encapsulate details of a request to be sent to server
	// Endpoint is dynamic that is appended to URL
	// e.g if the url is www.server.com/users/user-id, user-id is the endpoint
	Request struct {
		Name        string
		Method      string
		URL         string
		Endpoint    string
		BasicAuth   *BasicAuth
		Payload     interface{}
		Headers     map[string]string
		QueryParams map[string]string
	}

	RequestBuilder struct {
		name        string
		method      string
		url         string
		endpoint    string
		basicAuth   *BasicAuth
		payload     interface{}
		headers     map[string]string
		queryParams map[string]string
	}

	requestBuilder interface {
		Payload(interface{}) *RequestBuilder
		Headers(map[string]string) *RequestBuilder
		BasicAuth(auth *BasicAuth) *RequestBuilder
		QueryParams(params map[string]string) *RequestBuilder
		Endpoint(endpoint string) *RequestBuilder
		Build() *Request
	}

	RequestOption func(request *RequestBuilder)

	RequestInformer interface {
		fmt.Stringer
		Endpoint() string
		Method() string
		// Name returns the name of the request
		Name() string

		// MNO returns the name of the mno
		MNO() string
		Group() string
	}

	// RequestModifier is a func that can be injected into NewRequestWithContext
	// and transform the request before it is sent to server
	// be carefully using this func, as it can change the request completely
	RequestModifier func(request *http.Request) error
)

func (r *RequestBuilder) Payload(i interface{}) *RequestBuilder {
	r.payload = i
	return r
}

func (r *RequestBuilder) Headers(m map[string]string) *RequestBuilder {
	r.headers = m
	return r
}

func (r *RequestBuilder) BasicAuth(auth *BasicAuth) *RequestBuilder {
	r.basicAuth = auth
	return r
}

func (r *RequestBuilder) QueryParams(params map[string]string) *RequestBuilder {
	r.queryParams = params
	return r
}

func (r *RequestBuilder) Endpoint(endpoint string) *RequestBuilder {
	r.endpoint = endpoint
	return r
}

func (r *RequestBuilder) Build() *Request {
	return &Request{
		Name:        r.name,
		Method:      r.method,
		URL:         r.url,
		Endpoint:    r.endpoint,
		BasicAuth:   r.basicAuth,
		Payload:     r.payload,
		Headers:     r.headers,
		QueryParams: r.queryParams,
	}
}

func NewRequestBuilder(name, method, basePath string) *RequestBuilder {
	var defaultRequestHeaders = map[string]string{
		"Content-Type": cTypeJson,
	}

	if method == "" {
		method = http.MethodGet
	}
	rb := &RequestBuilder{
		name:    name,
		method:  method,
		url:     basePath,
		headers: defaultRequestHeaders,
	}

	return rb
}

var _ requestBuilder = (*RequestBuilder)(nil)

func MakeInternalRequest(basePath, endpoint string, requestType RequestInformer, payload interface{}, opts ...RequestOption) *Request {
	url := appendEndpoint(basePath, endpoint)
	method := requestType.Method()
	return NewRequest(requestType.String(), method, url, payload, opts...)
}

func NewRequest(name, method, url string, payload interface{}, opts ...RequestOption) *Request {

	rb := NewRequestBuilder(name, method, url)
	rb.payload = payload
	for _, opt := range opts {
		opt(rb)
	}

	return rb.Build()
}

func WithQueryParams(params map[string]string) RequestOption {
	return func(request *RequestBuilder) {
		request.queryParams = params
	}
}

func WithEndpoint(endpoint string) RequestOption {
	return func(request *RequestBuilder) {
		request.endpoint = endpoint
	}
}

// WithRequestHeaders replaces all the available HeaderMap with new ones
// WithMoreHeaders appends HeaderMap does not replace them
func WithRequestHeaders(headers map[string]string) RequestOption {
	return func(request *RequestBuilder) {
		request.headers = headers
	}
}

// WithMoreHeaders appends HeaderMap does not replace them like WithRequestHeaders
func WithMoreHeaders(headers map[string]string) RequestOption {
	return func(request *RequestBuilder) {
		for key, value := range headers {
			request.headers[key] = value
		}
	}
}

//// See 2 (end of page 4) https://www.ietf.org/rfc/rfc2617.txt
//// "To receive authorization, the client sends the userid and password,
//// separated by a single colon (":") character, within a base64
//// encoded string in the credentials."
//// It is not meant to be urlencoded.
//func basicAuth(username, password string) string {
//	auth := username + ":" + password
//	return base64.StdEncoding.EncodeToString([]byte(auth))
//}

// WithBasicAuth add password and username to request HeaderMap
func WithBasicAuth(username, password string) RequestOption {
	return func(request *RequestBuilder) {
		request.basicAuth = &BasicAuth{
			Username: username,
			Password: password,
		}
	}
}

func (request *Request) AddHeader(key, value string) {
	request.Headers[key] = value
}

func appendEndpoint(url string, endpoint string) string {
	url, endpoint = strings.TrimSpace(url), strings.TrimSpace(endpoint)
	urlHasSuffix, endpointHasPrefix := strings.HasSuffix(url, "/"), strings.HasPrefix(endpoint, "/")

	bothTrue := urlHasSuffix == true && endpointHasPrefix == true
	bothFalse := urlHasSuffix == false && endpointHasPrefix == false
	notEqual := urlHasSuffix != endpointHasPrefix

	if notEqual {
		return fmt.Sprintf("%s%s", url, endpoint)
	}

	if bothFalse {
		return fmt.Sprintf("%s/%s", url, endpoint)
	}

	if bothTrue {
		endp := strings.TrimPrefix(endpoint, "/")
		return fmt.Sprintf("%s%s", url, endp)
	}

	return ""
}

// NewRequestWithContext takes a *Request and transform into *http.Request with a context.Context
// It takes number of RequestModifier to modify the created *http.Request before it is used
// the modifiers will be applied in the order they are passed
// The modifier logic is completely up to the implementor of the modifier, so care should be taken
// to not modify the request in a way that it will break the sending logic.
func NewRequestWithContext(ctx context.Context, request *Request, modifiers ...RequestModifier) (req *http.Request, err error) {

	cType := request.Headers["Content-Type"]
	pType := categorizeContentType(cType)
	requestURL := request.URL
	requestEndpoint := request.Endpoint
	if requestEndpoint != "" {
		request.URL = appendEndpoint(requestURL, requestEndpoint)
	}
	if request.Payload == nil {
		req, err = http.NewRequestWithContext(ctx, request.Method, request.URL, nil)
		if err != nil {
			return nil, err
		}
	} else {
		buffer, err := MarshalPayload(pType, request.Payload)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequestWithContext(ctx, request.Method, request.URL, buffer)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	for name, value := range request.QueryParams {
		values := req.URL.Query()
		values.Add(name, value)
		req.URL.RawQuery = values.Encode()
	}

	ba := request.BasicAuth

	if ba != nil {
		req.SetBasicAuth(ba.Username, ba.Password)
	}

	for _, modifier := range modifiers {
		err := modifier(req)
		if err != nil {
			return nil, fmt.Errorf("error applying modifier: %w", err)
		}
	}

	return req, nil
}
