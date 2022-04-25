package test_lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Api struct {
	Server *httptest.Server
	T      *testing.T
}

func Marshal(t *testing.T, v any) string {
	var str string
	if bb, err := json.MarshalIndent(v, "", "\t"); err != nil {
		t.Fatal("marshal", err, v)
	} else {
		str = string(bb)
	}
	return str
}

func (api *Api) Marshal(v any) string {
	var str string
	if bb, err := json.MarshalIndent(v, "", "\t"); err != nil {
		api.T.Fatal("marshal", err, v)
	} else {
		str = string(bb)
	}
	return str
}

func QuoteMeta(r string) string {
	return "^" + regexp.QuoteMeta(r) + "$"
}

func (api *Api) HttpGet(path string, statusCode int, responseData any) (string, any) {
	var resp *http.Response
	if r, err := http.Get(api.Server.URL + path); err != nil {
		api.T.Fatal("Http Error", err)
	} else {
		resp = r
	}
	assert.Equal(api.T, statusCode, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		api.T.Fatal("Expected Content-Type header to be set")
	}
	assert.Equal(api.T, "application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		api.T.Fatal("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			assert.Fail(api.T, "Unmarshal body", err)
		}
	}
	assert.Equal(api.T, api.Marshal(responseData), api.Marshal(res), body)
	return body, res
}

func (api *Api) HttpPost(path string, requestData any, statusCode int, responseData any) (string, any) {
	requestDataBytes, _ := json.Marshal(requestData)

	var resp *http.Response
	if r, err := http.Post(api.Server.URL+path, "application/json; charset=utf-8", bytes.NewBuffer(requestDataBytes)); err != nil {
		api.T.Fatal("Http Error", err)
	} else {
		resp = r
	}
	assert.Equal(api.T, statusCode, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		api.T.Fatal("Expected Content-Type header to be set")
	}
	assert.Equal(api.T, "application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		api.T.Fatal("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			assert.Fail(api.T, "Unmarshal body", err)
		}
	}
	assert.Equal(api.T, api.Marshal(responseData), api.Marshal(res), body)
	return body, res
}

func (api *Api) HttpPut(path string, requestData any, statusCode int, responseData any) (string, any) {
	requestDataBytes, _ := json.Marshal(requestData)

	var resp *http.Response
	if req, err := http.NewRequest(http.MethodPut, api.Server.URL+path, bytes.NewBuffer(requestDataBytes)); err != nil {
		api.T.Fatal("Http Error", err)
	} else {
		client := &http.Client{}
		if r, err := client.Do(req); err != nil {
			api.T.Fatal("Http Error", err)
		} else {
			resp = r
		}
	}
	assert.Equal(api.T, statusCode, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		api.T.Fatal("Expected Content-Type header to be set")
	}
	assert.Equal(api.T, "application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		api.T.Fatal("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			assert.Fail(api.T, "Unmarshal body", err)
		}
	}
	assert.Equal(api.T, api.Marshal(responseData), api.Marshal(res), body)
	return body, res
}

func (api *Api) HttpPatch(path string, requestData any, statusCode int, responseData any) (string, any) {
	requestDataBytes, _ := json.Marshal(requestData)

	var resp *http.Response
	if req, err := http.NewRequest(http.MethodPatch, api.Server.URL+path, bytes.NewBuffer(requestDataBytes)); err != nil {
		api.T.Fatal("Http Error", err)
	} else {
		client := &http.Client{}
		if r, err := client.Do(req); err != nil {
			api.T.Fatal("Http Error", err)
		} else {
			resp = r
		}
	}
	assert.Equal(api.T, statusCode, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		api.T.Fatal("Expected Content-Type header to be set")
	}
	assert.Equal(api.T, "application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		api.T.Fatal("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			assert.Fail(api.T, "Unmarshal body", err)
		}
	}
	assert.Equal(api.T, api.Marshal(responseData), api.Marshal(res), body)
	return body, res
}

func (api *Api) HttpDelete(path string, statusCode int, responseData any) (string, any) {
	var resp *http.Response
	if req, err := http.NewRequest(http.MethodDelete, api.Server.URL+path, nil); err != nil {
		api.T.Fatal(api.T, "Http Error", err)
	} else {
		client := &http.Client{}
		if r, err := client.Do(req); err != nil {
			api.T.Fatal(api.T, "Http Error", err)
		} else {
			resp = r
		}
	}
	assert.Equal(api.T, statusCode, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		api.T.Fatal("Expected Content-Type header to be set")
	}
	assert.Equal(api.T, "application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		api.T.Fatal("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			assert.Fail(api.T, "Unmarshal body", err)
		}
	}
	assert.Equal(api.T, api.Marshal(responseData), api.Marshal(res), body)
	return body, res
}
