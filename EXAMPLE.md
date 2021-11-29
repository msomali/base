```go
package examples

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/base"
	"github.com/techcraftlabs/base/io"
	stdio "io"
	"net/http"
	"time"
)

type (
	User struct {
		RequestName  string `json:"name,omitempty"`
		Age   int    `json:"age,omitempty"`
		Email string `json:"email,omitempty"`
	}

	Resp struct {
		Message   string `json:"message,omitempty"`
		Error     string `json:"error,omitempty"`
		Timestamp int64  `json:"-"`
	}
	
	Client struct {
		logger stdio.Writer
		debugMode bool
		rv base.Receiver
		rp base.Replier
	}
	
	ClientOption func(client *Client)
)

const (
	defaultWriter = io.Stderr
	defaultDebugMode = false
)

func WithDebugMode(mode bool)ClientOption{
	return func(client *Client) {
		client.debugMode = mode
	}
}


func WithLogger(writer stdio.Writer)ClientOption{
	return func(client *Client) {
		client.logger = writer
	}
}
func NewClient(opts...ClientOption)*Client{
	c := new(Client)
	
	c = &Client{
		logger:    defaultWriter,
		debugMode: defaultDebugMode,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.rp = base.NewReplier(c.logger,c.debugMode)
	c.rv = base.NewReceiver(c.logger,c.debugMode)
	return c
}

func (c *Client)UserHandler(writer http.ResponseWriter, r *http.Request)  {
	ctx,cancel := context.WithTimeout(context.Background(),time.Minute)

	defer cancel()
	user := new(User)
	receipt, err := c.rv.Receive(ctx,"user", r,user)
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

	c.rp.Reply(writer,response)
}

```

```go
package examples

import (
	"context"
	"github.com/techcraftlabs/base"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler(t *testing.T) {
	u := User{
		RequestName:  "KingKaka",
		Age:   40,
		Email: "kaka@mkubwa.com",
	}

	var reqOpts []base.RequestOption
	bAuthOpt := base.WithBasicAuth("adminusername", "adminpassword")
	tokenAuth := base.WithMoreHeaders(map[string]string{
		"X-API-Key": "uongounonebegehgdahfsssfrtsrstfs",
	})
	reqOpts = append(reqOpts, bAuthOpt,tokenAuth)
	freq := base.NewRequest(http.MethodPost, "https://facebook.com", u, reqOpts...)

	req, err := base.NewRequestWithContext(context.TODO(), freq)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	c := NewClient(WithDebugMode(true))
	handler := http.HandlerFunc(c.UserHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"message":"username: adminusername, password: adminpassword apikey: uongounonebegehgdahfsssfrtsrstfs","error":"no error"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

```