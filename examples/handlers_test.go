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
		Name:  "KingKaka",
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
	handler := http.HandlerFunc(UserHandler)

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
