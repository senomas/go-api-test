package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuiteEnv struct {
	suite.Suite
	server *httptest.Server
}

func SetupDBMock(suite *TestSuiteEnv) (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
	var err error
	var sqlDB *sql.DB
	var mock sqlmock.Sqlmock
	sqlDB, mock, err = sqlmock.New()
	if err != nil {
		suite.Errorf(err, "Failed to open mock sql db")
	}
	var db *gorm.DB
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		suite.Errorf(err, "Failed to setup mockup gorm")
	}
	return db, sqlDB, mock
}

func (suite *TestSuiteEnv) Marshal(v any) string {
	var str string
	if bb, err := json.MarshalIndent(v, "", "\t"); err != nil {
		suite.Fail("marshal", err, v)
	} else {
		str = string(bb)
	}
	return str
}

func QuoteMeta(r string) string {
	return "^" + regexp.QuoteMeta(r) + "$"
}

func (suite *TestSuiteEnv) HttpGet(path string, responseData any) (string, any) {
	a := suite.Assert()

	var resp *http.Response
	if r, err := http.Get(suite.server.URL + path); err != nil {
		suite.Fail("Http Error", err)
	} else {
		resp = r
	}
	suite.Equal(200, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		suite.Fail("Expected Content-Type header to be set")
	}
	a.Equal("application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		suite.Fail("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			a.Fail("Unmarshal body", err)
		}
	}
	a.Equal(suite.Marshal(responseData), suite.Marshal(res), body)
	return body, res
}

func (suite *TestSuiteEnv) HttpPost(path string, requestData any, responseData any) (string, any) {
	a := suite.Assert()

	requestDataBytes, _ := json.Marshal(requestData)

	var resp *http.Response
	if r, err := http.Post(suite.server.URL+path, "application/json; charset=utf-8", bytes.NewBuffer(requestDataBytes)); err != nil {
		suite.Fail("Http Error", err)
	} else {
		resp = r
	}
	suite.Equal(200, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		suite.Fail("Expected Content-Type header to be set")
	}
	a.Equal("application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		suite.Fail("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			a.Fail("Unmarshal body", err)
		}
	}
	a.Equal(suite.Marshal(responseData), suite.Marshal(res), body)
	return body, res
}

func (suite *TestSuiteEnv) HttpPut(path string, requestData any, responseData any) (string, any) {
	a := suite.Assert()

	requestDataBytes, _ := json.Marshal(requestData)

	var resp *http.Response
	if req, err := http.NewRequest(http.MethodPut, suite.server.URL+path, bytes.NewBuffer(requestDataBytes)); err != nil {
		suite.Fail("Http Error", err)
	} else {
		client := &http.Client{}
		if r, err := client.Do(req); err != nil {
			suite.Fail("Http Error", err)
		} else {
			resp = r
		}
	}
	suite.Equal(200, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		suite.Fail("Expected Content-Type header to be set")
	}
	a.Equal("application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		suite.Fail("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			a.Fail("Unmarshal body", err)
		}
	}
	a.Equal(suite.Marshal(responseData), suite.Marshal(res), body)
	return body, res
}

func (suite *TestSuiteEnv) HttpPatch(path string, requestData any, responseData any) (string, any) {
	a := suite.Assert()

	requestDataBytes, _ := json.Marshal(requestData)

	var resp *http.Response
	if req, err := http.NewRequest(http.MethodPatch, suite.server.URL+path, bytes.NewBuffer(requestDataBytes)); err != nil {
		suite.Fail("Http Error", err)
	} else {
		client := &http.Client{}
		if r, err := client.Do(req); err != nil {
			suite.Fail("Http Error", err)
		} else {
			resp = r
		}
	}
	suite.Equal(200, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		suite.Fail("Expected Content-Type header to be set")
	}
	a.Equal("application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		suite.Fail("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			a.Fail("Unmarshal body", err)
		}
	}
	a.Equal(suite.Marshal(responseData), suite.Marshal(res), body)
	return body, res
}

func (suite *TestSuiteEnv) HttpDelete(path string, responseData any) (string, any) {
	a := suite.Assert()

	var resp *http.Response
	if req, err := http.NewRequest(http.MethodDelete, suite.server.URL+path, nil); err != nil {
		suite.Fail("Http Error", err)
	} else {
		client := &http.Client{}
		if r, err := client.Do(req); err != nil {
			suite.Fail("Http Error", err)
		} else {
			resp = r
		}
	}
	suite.Equal(200, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]

	if !ok {
		suite.Fail("Expected Content-Type header to be set")
	}
	a.Equal("application/json; charset=utf-8", val[0])
	var res map[string]json.RawMessage
	var body string

	if rb, err := ioutil.ReadAll(resp.Body); err != nil {
		suite.Fail("Read body", err)
	} else {
		body = string(rb)
		if err := json.Unmarshal(rb, &res); err != nil {
			a.Fail("Unmarshal body", err)
		}
	}
	a.Equal(suite.Marshal(responseData), suite.Marshal(res), body)
	return body, res
}
