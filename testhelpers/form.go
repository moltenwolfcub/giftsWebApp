package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func CreateFormRequest(target string, contents map[string]string) *http.Request {
	form := url.Values{}
	for k, v := range contents {
		form.Set(k, v)
	}

	body := strings.NewReader(form.Encode())
	req := httptest.NewRequest("POST", target, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}
