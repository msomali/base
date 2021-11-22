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
