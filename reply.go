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
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

var(
	_ Replier = (*replier)(nil)
)

type(
	replier struct {
		Logger io.Writer
		DebugMode bool
	}
	Replier interface {
		Reply(writer http.ResponseWriter, r *Response)
	}
)

func (rp *replier) Reply(writer http.ResponseWriter, response *Response) {
	responseFmt, err := responseFormat(response)
	if err != nil {
		return
	}
	defer func(debug bool) {
		if debug{
			_, _ = rp.Logger.Write([]byte(responseFmt))
		}
	}(rp.DebugMode)

	reply(writer,response)
}

func NewReplier(writer io.Writer, debug bool) Replier {
	return &replier{
		Logger:    writer,
		DebugMode: debug,
	}
}

func reply(writer http.ResponseWriter, r *Response) {
	if r.Payload == nil {

		for key, value := range r.Headers {
			writer.Header().Add(key, value)
		}
		writer.WriteHeader(r.StatusCode)

		return
	}

	pType := categorizeContentType(r.Headers["Content-Type"])

	switch pType {
	case XmlPayload:
		payload, err := xml.MarshalIndent(r.Payload, "", "  ")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		for key, value := range r.Headers {
			writer.Header().Set(key, value)
		}

		writer.WriteHeader(r.StatusCode)
		writer.Header().Set("Content-Type", cTypeAppXml)
		_, err = writer.Write(payload)
		return

	case JsonPayload:
		payload, err := json.Marshal(r.Payload)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		for key, value := range r.Headers {
			writer.Header().Set(key, value)
		}
		writer.Header().Set("Content-Type", cTypeJson)
		writer.WriteHeader(r.StatusCode)
		_, err = writer.Write(payload)
		return
	}
}
