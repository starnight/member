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

/* According to gin-csrf's description [1], the middleware has to be used with
 * gin-contrib/sessions. The HTTP requests must include the session information
 * in its HTTP header, otherwise the middleware cannot check the CSRF token. It
 * needs some information stored in the session.
 *
 * [1]: https://github.com/utrack/gin-csrf
 */
func copyCSRFToken(req *http.Request, res *httptest.ResponseRecorder) {
  req.Header.Set("X-CSRF-TOKEN", res.Header().Get("X-CSRF-TOKEN"))
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

  /* Bare request */
  res1 := httptest.NewRecorder()
  data := url.Values{}
  req1, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusBadRequest, res1.Code)
  assert.Equal(t, "CSRF token mismatch", res1.Body.String())

  /* Requests with session and the CSRF token, but bad POST form fields */
  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("GET", "/", nil)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusOK, res2.Code)
  assert.True(t, len(res2.Header().Get("X-CSRF-TOKEN")) > 0)

  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  copyCookies(req3, res2)
  copyCSRFToken(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusForbidden, res3.Code)
  assert.Equal(t, "Wrong account or password", res3.Body.String())

  res4 := httptest.NewRecorder()
  data.Set("account", "foo")
  req4, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req4.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req4, res2)
  copyCSRFToken(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusForbidden, res4.Code)
  assert.Equal(t, "Wrong account or password", res4.Body.String())
}

func TestLoginAndShowdate(t *testing.T) {
  r := setupRouter()

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  assert.True(t, len(res1.Header().Get("X-CSRF-TOKEN")) > 0)

  res2 := httptest.NewRecorder()
  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  copyCSRFToken(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusOK, res1.Code)

  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/showdate", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusOK, res3.Code)
  body := fmt.Sprintf("Welcome %s %s", data.Get("account"), time.Now().Format("2006-01-02"))
  assert.Equal(t, body, res3.Body.String()[:len(body)])

  res4 := httptest.NewRecorder()
  req4, _ := http.NewRequest("GET", "/", nil)
  copyCookies(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusOK, res4.Code)
  assert.Equal(t, "<h1>Welcome foo</h1>\n", res4.Body.String())
}
