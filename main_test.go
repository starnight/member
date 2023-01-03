package main

import (
  "fmt"
  "strings"
  "testing"
  "time"
  "net/http/httptest"
  "net/http"
  "net/url"
  "github.com/stretchr/testify/assert"
)

func copyCookies(req *http.Request, res *httptest.ResponseRecorder) {
  req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func TestPing(t *testing.T) {
  r := setupRouter()

  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/ping", nil)
  r.ServeHTTP(res, req)

  assert.Equal(t, http.StatusOK, res.Code)
  assert.Equal(t, "pong", res.Body.String())
}

func TestRootWithoutLogin(t *testing.T) {
  r := setupRouter()

  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/", nil)
  r.ServeHTTP(res, req)

  assert.Equal(t, http.StatusOK, res.Code)
  assert.Equal(t, "Welcome", res.Body.String())
}

func TestShowdateWithoutLogin(t *testing.T) {
  r := setupRouter()

  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/showdate", nil)
  r.ServeHTTP(res, req)

  assert.Equal(t, http.StatusForbidden, res.Code)
  assert.Equal(t, "Please login first", res.Body.String())
}

func TestFailedLogin(t *testing.T) {
  r := setupRouter()

  res1 := httptest.NewRecorder()
  data := url.Values{}
  req1, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusForbidden, res1.Code)
  assert.Equal(t, "Wrong account or password", res1.Body.String())

  res2 := httptest.NewRecorder()
  data.Set("account", "foo")
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusForbidden, res2.Code)
  assert.Equal(t, "Wrong account or password", res2.Body.String())
}

func TestLoginAndShowdate(t *testing.T) {
  r := setupRouter()

  res1 := httptest.NewRecorder()
  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  req1, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)

  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("GET", "/showdate", nil)
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusOK, res2.Code)
  body := fmt.Sprintf("Welcome %s %s", data.Get("account"), time.Now().Format("2006-01-02"))
  assert.Equal(t, body, res2.Body.String()[:len(body)])
}
