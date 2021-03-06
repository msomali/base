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
	"errors"
	"fmt"
	stdio "io"
	"net/http"
	"strings"
)

const errStatusCodeMargin = 400

var DoErr = errors.New("result code is above or equal to 400")

// Do perform http request and return *Response. It takes *Request and body as input. It will inspect
// the header of the http.Response then determine the Content-Type of the response body. It will then
// unmarshal the content of the response body to the specified type. Error returned by this function
// is operation error. In case the response status code is equal or above to 400 and the operations like
// unmarshalling or reading header have all gone correctly the error will be nil but Response.Error will not.
func (c *Client) Do(ctx context.Context, request *Request, body interface{}, modifiers ...RequestModifier) (*Response, error) {

	var (
		rn               = request.Name
		errDecodingBody  = errors.New("error while decoding response body")
		errUnknownHeader = errors.New("unknown content-type header")
	)

	var (
		_, cancel    = context.WithTimeout(ctx, defaultTimeout)
		req          *http.Request
		res          *http.Response
		reqBodyBytes []byte
		resBodyBytes []byte
	)
	defer cancel()
	defer func(debug bool) {
		if debug {
			req.Body = stdio.NopCloser(bytes.NewBuffer(reqBodyBytes))
			if res == nil {
				c.logOut(strings.ToUpper(rn), req, nil)

				return
			}
			res.Body = stdio.NopCloser(bytes.NewBuffer(resBodyBytes))
			c.logOut(strings.ToUpper(rn), req, res)
		}
	}(c.DebugMode)
	req, err := NewRequestWithContext(ctx, request, modifiers...)

	if err != nil {
		return nil, err
	}

	if req.Body != nil {
		reqBodyBytes, _ = stdio.ReadAll(req.Body)
	}

	req.Body = stdio.NopCloser(bytes.NewBuffer(reqBodyBytes))
	res, doErr := c.Http.Do(req)

	if doErr != nil {
		return nil, doErr
	}

	if res.Body != nil {
		resBodyBytes, _ = stdio.ReadAll(res.Body)
	}

	response := new(Response)
	statusCode := res.StatusCode
	response.StatusCode = statusCode
	response.HTTP = res

	//change res.Header to map[string]string
	header := make(map[string]string)
	for k, v := range res.Header {
		header[k] = v[0]
	}
	response.HeaderMap = header

	contentType := res.Header.Get("Content-Type")
	headers := make(map[string]string)
	for k, v := range res.Header {
		headers[strings.ToLower(k)] = v[0]
	}

	response.HeaderMap = headers
	cType := categorizeContentType(contentType)

	isJSON := cType == JsonPayload
	isXML := cType == XmlPayload || cType == TextXmlPayload
	isOK := statusCode < errStatusCodeMargin

	if body != nil {
		if isJSON {
			dErr := json.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(body)
			isDecodeErr := dErr != nil && !errors.Is(dErr, stdio.EOF)

			if isDecodeErr {
				return nil, fmt.Errorf("%w: %v", dErr, errDecodingBody)
			}

			response.Body = body

			if !isOK {
				response.Error = DoErr
				return response, nil
			}

			return response, nil

		} else if isXML {

			dErr := xml.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(body)
			isDecodeErr := dErr != nil && !errors.Is(dErr, stdio.EOF)
			if isDecodeErr {
				return nil, fmt.Errorf("%w: %v", dErr, errDecodingBody)
			}

			response.Body = body
			if !isOK {
				response.Error = DoErr
				return response, nil
			}
			return response, nil

		} else {
			//response.Error = errUnknownHeader
			return nil, errUnknownHeader
		}
	}

	if !isOK {
		response.Error = DoErr
	}
	return response, nil

}
