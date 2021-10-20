package examples

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/base"
	"github.com/techcraftlabs/base/io"
	"net/http"
	"time"
)

type (
	User struct {
		Name  string `json:"name,omitempty"`
		Age   int    `json:"age,omitempty"`
		Email string `json:"email,omitempty"`
	}

	Resp struct {
		Message   string `json:"message,omitempty"`
		Error     string `json:"error,omitempty"`
		Timestamp int64  `json:"-"`
	}
)

func UserHandler(writer http.ResponseWriter, r *http.Request)  {
	ctx,cancel := context.WithTimeout(context.Background(),time.Minute)
	rp := base.NewReplier(io.Stderr,false)
	rv := base.NewReceiver(io.Stderr,false)
	defer cancel()
	user := new(User)
	receipt, err := rv.Receive(ctx,"user", r,user)
	if err != nil {
		http.Error(writer, err.Error(),http.StatusInternalServerError)
		return
	}

	uname, pass := receipt.BasicAuth.Username, receipt.BasicAuth.Password
	token := receipt.ApiKey
	message := fmt.Sprintf("username: %s, password: %s apikey: %s",uname,pass,token)

	payload := Resp{
		Message:   message,
		Error:     "no error",
		Timestamp: time.Now().UnixNano(),
	}

	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
		"X-Server-Info":"local-man",
		"X-Timestamp": fmt.Sprintf("%v",time.Now().Format(time.RFC850)),
	})

	response := base.NewResponse(http.StatusOK,payload,headersOption)

	rp.Reply(writer,response)
}
