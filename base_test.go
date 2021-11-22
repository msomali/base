package base

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"testing"
)

type User struct {
	XMLName xml.Name `xml:"user"`
	Name    string   `json:"name,omitempty" xml:"name"`
	Age     int64    `json:"age,omitempty" xml:"age"`
	Email   string   `json:"email,omitempty" xml:"email"`
	Job     string   `json:"job,omitempty" xml:"job"`
}

func Test_BuildRequest(t *testing.T) {

	rv := NewReceiver(os.Stderr, true)

	user := struct {
		XMLName xml.Name `xml:"user"`
		Name    string   `json:"name" xml:"name"`
		Age     int64    `json:"age" xml:"age"`
	}{
		Name: "John Doe",
		Age:  27,
	}

	user1 := struct {
		XMLName xml.Name `xml:"user"`
		Name    string   `json:"name" xml:"name"`
		Age     int64    `json:"age" xml:"age"`
		Email   string   `json:"email" xml:"email"`
		Job     string   `json:"job" xml:"job"`
	}{
		Name:  "John Doe",
		Age:   27,
		Email: "jdoe@anon.com",
		Job:   "imposter",
	}

	headerModifier := func(headers map[string]string) RequestModifier {
		return func(req *http.Request) {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}
	}

	payloadModifier := func(payload interface{}) RequestModifier {
		return func(req *http.Request) {
			pt := categorizeContentType(req.Header.Get("Content-Type"))
			buf, _ := MarshalPayload(pt, payload)
			req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		}
	}

	request := NewRequestBuilder("login request", http.MethodPost, "https://google.com/login").
		Endpoint("/account-id").
		BasicAuth(&BasicAuth{Username: "johndoe", Password: "jd2021"}).
		Headers(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json"}).
		Payload(user).
		QueryParams(map[string]string{
			"page":  "1",
			"limit": "10"}).
		Build()

	newHeaders := map[string]string{
		"Content-Type":  "application/xml",
		"Accept":        "*/*",
		"Authorization": "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
	}

	rq, err := NewRequestWithContext(context.TODO(), request, headerModifier(newHeaders), payloadModifier(user1))
	if err != nil {
		t.Errorf("error creating request: %s", err.Error())
	}

	if rq.Header.Get("Content-Type") != "application/xml" {
		t.Errorf("Content-Type header is not set correctly")
	}

	u := new(User)

	_, err = rv.Receive(context.TODO(), request.Name, rq, u)
	if err != nil {
		t.Errorf("error receiving request: %s", err.Error())
	}

}
