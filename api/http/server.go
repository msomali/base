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
	"github.com/gorilla/mux"
	"github.com/techcraftt/base"
	"github.com/techcraftt/base/api"
	"github.com/techcraftt/base/internal"
	stdhttp "net/http"
)

var _ base.Service = (*svc)(nil)

type (
	svc    struct{}
	Server struct {
		Port    uint64
		server  *stdhttp.Server
		Debug   bool
		service base.Service
	}
)

func NewServer(port uint64, debug bool) *Server {
	sv := &Server{
		Port:    port,
		Debug:   debug,
		service: &svc{},
	}

	h := sv.handler()

	sv.server = &stdhttp.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h,
	}

	return sv
}
func sendErrMessage(writer stdhttp.ResponseWriter, code int, err api.Response) {
	response := internal.NewResponse(code, err)
	internal.Reply(writer,response)
}

func (server *Server) divisionHandler(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
	defer func() {
		if r := recover(); r != nil {
			errMessage := api.Response{
				Err:     "division by zero",
				Message: "division by zero is not good dont do it",
			}
			sendErrMessage(writer, stdhttp.StatusBadRequest, errMessage)
			return
		}
	}()
	var req api.Request
	err := internal.ReceivePayload(request, &req)
	if err != nil {
		errMessage := api.Response{
			Err:     err.Error(),
			Message: "failed to obtain request body",
		}
		sendErrMessage(writer, stdhttp.StatusInternalServerError, errMessage)
		return
	}

	rs, err := server.service.Divide(req.A, req.B)
	if err != nil {
		errMessage := api.Response{
			Err:     err.Error(),
			Message: "failed to perform division",
		}

		r := internal.NewResponse(stdhttp.StatusInternalServerError, errMessage)
		internal.Reply(writer,r)
		return
	}
	result := api.Response{
		Answer: rs,
	}
	response := internal.NewResponse(stdhttp.StatusOK, result)
	internal.Reply(writer,response)
}

func (server *Server) handler() stdhttp.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/div", server.divisionHandler).Methods(stdhttp.MethodGet)
	return r
}

func (server *Server) ListenAndServe() error {

	return server.server.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

func (s *svc) Add(a, b int64) (int64, error) {
	return a + b, nil
}

func (s *svc) Divide(a, b int64) (int64, error) {
	answer := a / b
	return answer, nil
}
