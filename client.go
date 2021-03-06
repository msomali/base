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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/techcraftlabs/base/io"
	stdio "io"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

const (
	defaultTimeout = 60 * time.Second
)

type (
	Client struct {
		mu        sync.Mutex
		Http      *http.Client
		Logger    stdio.Writer // for logging purposes
		DebugMode bool
		certPool  *x509.CertPool
	}

	ClientOption func(client *Client)
)

//SetLogger set logger for client
func (c *Client) SetLogger(writer stdio.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if writer != nil {
		c.Logger = writer
	}
}

//SetDebugMode set debug mode for client
func (c *Client) SetDebugMode(debugMode bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DebugMode = debugMode
}

func NewClient(opts ...ClientOption) *Client {
	defClient := &http.Client{
		Timeout: defaultTimeout,
	}
	client := &Client{
		mu:        sync.Mutex{},
		Http:      defClient,
		Logger:    io.StdErr,
		DebugMode: true,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) logPayload(t PayloadType, prefix string, payload interface{}) {
	buf, _ := MarshalPayload(t, payload)
	_, _ = c.Logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, buf.String())))
}

func (c *Client) log(name string, request *http.Request) {

	if request != nil {
		reqDump, _ := httputil.DumpRequest(request, true)
		_, err := fmt.Fprintf(c.Logger, "\n\n%s REQUEST: %s\n\n", name, reqDump)
		if err != nil {
			fmt.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
			return
		}
		return
	}
	return
}

// logOut is like log except this is for outgoing client requests:
// http.Request that is supposed to be sent to tigo
func (c *Client) logOut(name string, request *http.Request, response *http.Response) {

	if request != nil {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, err := fmt.Fprintf(c.Logger, "\n\n%s REQUEST (OUTGOING)\n%s\n\n", name, reqDump)
		if err != nil {
			fmt.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
		}
	}

	if response != nil {
		respDump, _ := httputil.DumpResponse(response, true)
		_, err := fmt.Fprintf(c.Logger, "\n\n%s RESPONSE\n%s\n\n", name, respDump)
		if err != nil {
			fmt.Printf("error while logging %s response: %v\n",
				strings.ToLower(name), err)
		}
	}

	return
}

// WithDebugMode set debug mode to true or false
func WithDebugMode(debugMode bool) ClientOption {
	return func(client *Client) {
		client.DebugMode = debugMode

	}
}

// WithLogger set a Logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out stdio.Writer) ClientOption {
	return func(client *Client) {
		if out == nil {
			return
		}
		client.Logger = out
	}
}

// WithHTTPClient when called unset the present http.Client and replace it
// with c. In case user tries to pass a nil value referencing the pkg
// i.e. WithHTTPClient(nil), it will be ignored and the pkg will not be replaced
// Note: the new pkg Transport will be modified. It will be wrapped by another
// middleware that enables pkg to
func WithHTTPClient(httpClient *http.Client) ClientOption {

	// TODO check if its really necessary to set the default Timeout to 1 minute

	return func(client *Client) {
		if httpClient == nil {
			return
		}

		client.Http = httpClient
	}
}

func WithCACert(caCert []byte) ClientOption {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return func(client *Client) {
		if caCert == nil {
			return
		}

		c := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
			CheckRedirect: client.Http.CheckRedirect,
			Jar:           client.Http.Jar,
			Timeout:       client.Http.Timeout,
		}

		client.Http = c
	}

}
