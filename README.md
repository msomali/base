# base
the base code for creating mobile money api clients using golang 

## build request
```go

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


```

## modifiers

```go

RequestModifier func(request *http.Request)

//example
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
```
