package vanguard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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
}

func (f *serverFixture) SendRequest(method, path string, body interface{}) bool {
	assert := assert.New(f.t)

	f.res = httptest.NewRecorder()

	var r io.Reader
	if body != nil {
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

func TestServerProvision(t *testing.T) {
	assert := assert.New(t)

	tmpfile, err := ioutil.TempFile(os.TempDir(), "vanguard-test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Chmod(0400)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	yamlFile := fmt.Sprintf(`---
tasks:
  - name: modify tmpfile
    file:
      path: "%s"
      mode: 0644
`, tmpfile.Name())

	f := ServerFixture(t)
	f.SendRequest("POST", "/v1/provision", yamlFile)

	assert.Equal(http.StatusOK, f.res.Code)
	f.SendRequest("GET", fmt.Sprintf("/v1/runs/%s/wait", f.res.Body.String()), nil)

	st, err := os.Stat(tmpfile.Name())
	if assert.NoError(err) {
		assert.Equal(os.FileMode(0644), st.Mode())
	}
}
