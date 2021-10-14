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

import "strings"

const (
	JsonPayload PayloadType = iota
	XmlPayload
	FormPayload
	TextXml
	UnsupportedPayload
)

const (
	cTypeJson        = "application/json"
	cTypeTextXml     = "text/xml"
	cTypeAppXml      = "application/xml"
	cTypeForm        = "application/x-www-form-urlencoded"
	cTypeUnsupported = "unsupported"
)

type (
	PayloadType int
)

func (p PayloadType) String() string {
	types := []string{
		cTypeJson,
		cTypeAppXml,
		cTypeForm,
		cTypeTextXml,
		cTypeUnsupported,
	}

	return types[p]
}

func categorizeContentType(headerStr string) PayloadType {
	j := JsonPayload.String()
	x := XmlPayload.String()
	xml2 := TextXml.String()
	form := FormPayload.String()
	if strings.Contains(headerStr, j) {
		return JsonPayload
	} else if strings.Contains(headerStr, xml2) || strings.Contains(headerStr, x) {
		return XmlPayload
	} else if strings.Contains(headerStr, form) {
		return FormPayload
	} else {
		//todo: figure out proper way to return this
		return UnsupportedPayload
	}
}
