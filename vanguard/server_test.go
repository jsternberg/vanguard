package vanguard

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type serverFixture struct {
	t   *testing.T
	s   *Server
	res *httptest.ResponseRecorder
}

func ServerFixture(t *testing.T) *serverFixture {
	f := &serverFixture{}
	f.SetUp(t)
	return f
}

func (f *serverFixture) SetUp(t *testing.T) {
	f.t = t
	f.s = NewServer()
	f.res = httptest.NewRecorder()
}

func (f *serverFixture) SendRequest(method, path string, body interface{}) bool {
	assert := assert.New(f.t)

	var r io.Reader
	if body == nil {
		switch obj := body.(type) {
		case string:
			r = bytes.NewBufferString(obj)
		default:
			out, err := json.Marshal(obj)
			if !assert.NoError(err) {
				return false
			}
			r = bytes.NewBuffer(out)
		}
	}

	req, err := http.NewRequest(method, path, r)
	if assert.NoError(err) {
		f.s.handler.ServeHTTP(f.res, req)
		return true
	}
	return false
}

func (f *serverFixture) Json() interface{} {
	var out interface{}
	decoder := json.NewDecoder(f.res.Body)
	if err := decoder.Decode(&out); err != nil {
		f.t.Fatal(err)
	}
	return out
}

func TestServerPing(t *testing.T) {
	assert := assert.New(t)

	f := ServerFixture(t)
	f.SendRequest("GET", "/v1/ping", nil)

	if assert.Equal(200, f.res.Code) {
		serverInfo := f.Json().(map[string]interface{})
		assert.Equal(Version, serverInfo["version"])
	}
}
