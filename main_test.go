package main

import (
  "fmt"
  "strings"
  "testing"
  "time"
  "golang.org/x/net/html"
  "net/http/httptest"
  "net/http"
  "net/url"
  "github.com/stretchr/testify/assert"
)

func copyCookies(req *http.Request, res *httptest.ResponseRecorder) {
  req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func getElementById(id string, n *html.Node) (element *html.Node, ok bool) {
  for _, a := range n.Attr {
    if (a.Key == "id" && a.Val == id) {
      return n, true
    }
  }

  element = nil
  ok = false

  for c := n.FirstChild; c != nil; c = c.NextSibling {
    element, ok = getElementById(id, c)
    if ok {
      break
    }
  }

  return element, ok
}

func getCSRFToken(res *httptest.ResponseRecorder) (token string) {
  token = ""

  root, _ := html.Parse(res.Body)
  csrf_node, ok := getElementById("_csrf", root)

  if (!ok) {
    return token
  }

  for _, a := range csrf_node.Attr {
    if (a.Key == "value") {
      token = a.Val
      break
    }
  }

  return token
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

  expected_body := "<h1>Welcome</h1>"

  assert.Equal(t, http.StatusOK, res.Code)
  assert.Equal(t, expected_body, res.Body.String()[:len(expected_body)])
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
  req2, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res2, req2)

  expected_login := "<h1>Login</h1>"
  expected_csrf := "name=\"_csrf\""

  assert.Equal(t, http.StatusOK, res2.Code)
  assert.True(t, strings.Contains(res2.Body.String(), expected_login))
  assert.True(t, strings.Contains(res2.Body.String(), expected_csrf))

  csrf_token := getCSRFToken(res2)
  assert.True(t, len(csrf_token) > 0)

  res3 := httptest.NewRecorder()
  data.Set("_csrf", csrf_token)
  req3, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

  /* According to gin-csrf's description [1], the middleware has to be used with
   * gin-contrib/sessions. The HTTP requests must include the session information
   * in its HTTP header, otherwise the middleware cannot check the CSRF token. It
   * needs some information stored in the session.
   *
   * [1]: https://github.com/utrack/gin-csrf
   */
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusForbidden, res3.Code)
  assert.Equal(t, "Wrong account or password", res3.Body.String())

  res4 := httptest.NewRecorder()
  data.Set("account", "foo")
  req4, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req4.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusForbidden, res4.Code)
  assert.Equal(t, "Wrong account or password", res4.Body.String())
}

func TestLoginAndShowdate(t *testing.T) {
  r := setupRouter()

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  csrf_token := getCSRFToken(res1)
  assert.True(t, len(csrf_token) > 0)

  res2 := httptest.NewRecorder()
  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  data.Set("_csrf", csrf_token)
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusFound, res2.Code)
  assert.Equal(t, "/", res2.Header().Get("Location"))
  assert.Equal(t, "", res2.Body.String())

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
