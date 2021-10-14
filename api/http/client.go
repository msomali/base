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

package http

import (
	"context"
	"fmt"
	"github.com/techcraftt/base"
	"github.com/techcraftt/base/api"
	"github.com/techcraftt/base/internal"
	"net/http"
)

var _ base.Service = (*Client)(nil)

type (
	Client struct {
		BaseURL string
		Port    uint64
		base    *internal.BaseClient
	}
)

func (c *Client) Divide(a int64, b int64) (int64, error) {
	req := api.Request{
		A: a,
		B: b,
	}
	request := internal.NewRequest(http.MethodGet, c.requestURL("div"), req)

	response := new(api.Response)
	do, err := c.base.Do(context.TODO(), "divide", request, response)

	if err != nil {
		return 0, err
	}

	if do.Error != nil {
		return 0, fmt.Errorf("%s:%s", response.Err, response.Message)
	}

	return response.Answer, nil
}

func NewClient(base string, port uint64, debug bool) *Client {
	return &Client{
		BaseURL: base,
		Port:    port,
		base:    internal.NewBaseClient(internal.WithDebugMode(debug)),
	}
}

func (c *Client) requestURL(endpoint string) string {
	return fmt.Sprintf("http://%s:%d/%s", c.BaseURL, c.Port, endpoint)
}
