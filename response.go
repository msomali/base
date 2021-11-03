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
	"fmt"
	"strings"
)

type (
	// Response contains details to be sent as a response to tigo
	// when tigo make callback request, name check request or payment
	// request.
	Response struct {
		statusCode int
		payload    interface{}
		headers    map[string]string
		error      error
	}

	ResponseBuilder struct {
		statusCode int
		payload    interface{}
		headers    map[string]string
		error      error
	}

	responseBuilder interface {
		StatusCode(int)*ResponseBuilder
		Payload(interface{})*ResponseBuilder
		Headers(map[string]string)*ResponseBuilder
		Error(error)*ResponseBuilder
		Build()*Response
	}

	ResponseOption func(response *ResponseBuilder)
)

func (r *ResponseBuilder) StatusCode(i int) *ResponseBuilder {
	r.statusCode = i
	return r
}

func (r *ResponseBuilder) Payload(i interface{}) *ResponseBuilder {
	r.payload = i
	return r
}

func (r *ResponseBuilder) Headers(m map[string]string) *ResponseBuilder {
	r.headers = m
	return r
}

func (r *ResponseBuilder) Error(err error) *ResponseBuilder {
	r.error = err
	return r
}

func (r *ResponseBuilder) Build() *Response {
	return &Response{
		statusCode: r.statusCode,
		payload:    r.payload,
		headers:    r.headers,
		error:      r.error,
	}
}

func NewResponseBuilder()*ResponseBuilder{
	return &ResponseBuilder{
		statusCode: 200,
		payload:    nil,
		headers: map[string]string{
			"Content-Type":cTypeJson,
		},
		error:      nil,
	}
}


// NewResponse create a response to be sent back to Tigo. HTTP Status code, payload and its
// type need to be specified. Other fields like  Response.error and Response.headers can be
// changed using WithMoreResponseHeaders (add headers), WithResponseHeaders (replace all the
// existing ) and WithResponseError to add error its default value is nil, default value of
// Response.headers is
// defaultResponseHeader = map[string]string{
//		"Content-Type": ContentTypeXml,
// }
func NewResponse(status int, payload interface{}, opts ...ResponseOption) *Response {


	var (
		defaultResponseHeader = map[string]string{
			"Content-Type": cTypeJson,
		}
	)

	rb := &ResponseBuilder{
		statusCode: status,
		payload:    payload,
		headers:    defaultResponseHeader,
		error:      nil,
	}

	for _, opt := range opts {
		opt(rb)
	}

	return rb.Build()
}

func WithResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *ResponseBuilder) {
		response.headers = headers
	}
}

func WithMoreResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *ResponseBuilder) {
		for key, value := range headers {
			response.headers[key] = value
		}
	}
}

func WithDefaultXMLHeader() ResponseOption {
	return func(response *ResponseBuilder) {
		response.headers["Content-Type"] = cTypeAppXml
	}
}

func WithResponseError(err error) ResponseOption {

	return func(response *ResponseBuilder) {
		response.error = err
	}
}

func responseFormat(response *Response) (string,error) {

	var(
		errMsg string
	)
	if response == nil{
		return "",fmt.Errorf("response is nil")
	}
	hs := response.headers
	statusCode := response.statusCode

	builder := strings.Builder{}
	for key, val := range hs {
		builder.WriteString(fmt.Sprintf("%s: %s\n",key,val))
	}

	headersString := builder.String()
	if response.error != nil{
		errMsg = response.error.Error()
	}
	if response.error == nil{
		errMsg = "nil"
	}

	contentType := response.headers["Content-Type"]
	payloadType := categorizeContentType(contentType)
	buffer, err := MarshalPayload(payloadType,response.payload)
	if err != nil{
		return "",err
	}
	payload := buffer.String()

	fmtString := fmt.Sprintf("\nRESPONSE DUMP:\nstatus code: %d\nheaders: %sother details:\nerror: %s\npayload: %s\n",statusCode,headersString,errMsg,payload)
	return fmtString, nil
}
