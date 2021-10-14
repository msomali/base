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

import "fmt"

type (
	// Response contains details to be sent as a response to tigo
	// when tigo make callback request, name check request or payment
	// request.
	Response struct {
		StatusCode int               `json:"status_code"`
		Payload    interface{}       `json:"payload"`
		Headers    map[string]string `json:"headers"`
		Error      error             `json:"error"`
	}

	ResponseOption func(response *Response)
)

// NewResponse create a response to be sent back to Tigo. HTTP Status code, payload and its
// type need to be specified. Other fields like  Response.Error and Response.Headers can be
// changed using WithMoreResponseHeaders (add headers), WithResponseHeaders (replace all the
// existing ) and WithResponseError to add error its default value is nil, default value of
// Response.Headers is
// defaultResponseHeader = map[string]string{
//		"Content-Type": ContentTypeXml,
// }
func NewResponse(status int, payload interface{}, opts ...ResponseOption) *Response {

	var (
		defaultResponseHeader = map[string]string{
			"Content-Type": cTypeJson,
		}
	)

	response := &Response{
		StatusCode: status,
		Payload:    payload,
		Headers:    defaultResponseHeader,
		Error:      nil,
	}

	for _, opt := range opts {
		opt(response)
	}

	return response
}

func WithResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *Response) {
		response.Headers = headers
	}
}

func WithMoreResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *Response) {
		for key, value := range headers {
			response.Headers[key] = value
		}
	}
}

func WithDefaultXMLHeader() ResponseOption {
	return func(response *Response) {
		response.Headers["Content-Type"] = cTypeAppXml
	}
}

func WithResponseError(err error) ResponseOption {

	return func(response *Response) {
		response.Error = err
	}
}

func ResponseFmt(response *Response) (string,error) {

	var(
		errMsg string
	)
	if response == nil{
		return "",fmt.Errorf("response is nil")
	}
	hs := response.Headers
	statusCode := response.StatusCode

	if response.Error != nil{
		errMsg = response.Error.Error()
	}
	if response.Error == nil{
		errMsg = "nil"
	}

	contentType := response.Headers["Content-Type"]
	payloadType := categorizeContentType(contentType)
	buffer, err := MarshalPayload(payloadType,response.Payload)
	if err != nil{
		return "",err
	}
	payload := buffer.String()

	fmtString := fmt.Sprintf("\nRESPONSE DUMP:\nstatus code: %d\nheaders: %v\nerror: %s\npayload: %s\n",statusCode,hs,errMsg,payload)
	return fmtString, nil
}
