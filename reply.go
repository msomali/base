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
	"sync"
)

var (
	_ Replier = (*replier)(nil)
)

type (
	replier struct {
		mu        sync.Mutex
		Logger    io.Writer
		DebugMode bool
	}
	Replier interface {
		Reply(writer http.ResponseWriter, r *Response, opts ...OptionFunc)
	}
)

func (rp *replier) update(params *Params) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	if params != nil {
		rp.DebugMode = params.DebugMode
		rp.Logger = params.Logger
	}
}

func (rp *replier) Reply(writer http.ResponseWriter, response *Response, opts ...OptionFunc) {
	params := &Params{
		DebugMode: rp.DebugMode,
		Logger:    rp.Logger,
	}
	for _, opt := range opts {
		opt(params)
	}

	rp.update(params)

	responseFmt, err := responseFormat(response)
	if err != nil {
		return
	}
	defer func(debug bool) {
		if debug {
			_, _ = rp.Logger.Write([]byte(responseFmt))
		}
	}(rp.DebugMode)

	reply(writer, response)
}

func NewReplier(writer io.Writer, debug bool) Replier {
	return &replier{
		mu:        sync.Mutex{},
		Logger:    writer,
		DebugMode: debug,
	}
}

func reply(writer http.ResponseWriter, r *Response) {
	if r.Body == nil {

		for key, value := range r.HeaderMap {
			writer.Header().Add(key, value)
		}
		writer.WriteHeader(r.StatusCode)

		return
	}

	cType := r.HeaderMap["Content-Type"]
	pType := categorizeContentType(cType)
	payload := r.Body
	switch pType {
	case XmlPayload:
		payload1, err := xml.MarshalIndent(payload, "", "  ")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		for key, value := range r.HeaderMap {
			writer.Header().Set(key, value)
		}

		writer.WriteHeader(r.StatusCode)
		writer.Header().Set("Content-Type", cTypeAppXml)
		_, err = writer.Write(payload1)
		return

	case JsonPayload:
		payload1, err := json.Marshal(payload)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		for key, value := range r.HeaderMap {
			writer.Header().Set(key, value)
		}
		writer.Header().Set("Content-Type", cTypeJson)
		writer.WriteHeader(r.StatusCode)
		_, err = writer.Write(payload1)
		return
	}
}
