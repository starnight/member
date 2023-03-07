package main

import (
  "github.com/gin-gonic/gin"
  "github.com/starnight/member/database"

  "bytes"
  "os"
  "strings"
  "testing"
  "golang.org/x/net/html"
  "net/http/httptest"
  "net/http"
  "net/url"
  "github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
  var cfg configSet

  os.Args = []string{"example"}
  cfg = parseArgs()
  assert.Equal(t, ":8080", cfg.AddrStr)

  os.Setenv("PORT", "8787")
  cfg = parseArgs()
  assert.Equal(t, ":8787", cfg.AddrStr)

  os.Args = []string{"example", "--port", "8765"}
  cfg = parseArgs()
  assert.Equal(t, ":8765", cfg.AddrStr)
}

func TestSetupDB(t *testing.T) {
  setupDB()

  dbstr := database.GetDBStr(gin.Mode())
  db := database.ConnectDB(dbstr)

  assert.NotNil(t, db)

  utils := database.UserUtils{DB: db}

  cnt, err := utils.Count()

  assert.Equal(t, cnt, int64(0))
  assert.Nil(t, err)
}

func TestPing(t *testing.T) {
  r := setupRouter()

  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/ping", nil)
  r.ServeHTTP(res, req)

  assert.Equal(t, http.StatusOK, res.Code)
  assert.Equal(t, "pong", res.Body.String())
}

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

func TestAdd1stUser(t *testing.T) {
  r := setupRouter()

  /* Request root path when none existed user */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusFound, res1.Code)
  assert.Equal(t, "/add1stuser", res1.Header().Get("Location"))

  /* Request add1stuser with GET method to have session and CSRF token */
  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("GET", "/add1stuser", nil)
  r.ServeHTTP(res2, req2)

  expected_adduser := "<h1>Add User</h1>"
  expected_csrf := "name=\"_csrf\""

  assert.Equal(t, http.StatusOK, res2.Code)
  assert.True(t, strings.Contains(res2.Body.String(), expected_adduser))
  assert.True(t, strings.Contains(res2.Body.String(), expected_csrf))

  csrf_token := getCSRFToken(res2)
  assert.True(t, len(csrf_token) > 0)

  /* Requests with session, the CSRF token, and correct POST form fields */
  res3 := httptest.NewRecorder()
  data1 := []byte(`{"account":"foo","passwd":"bar","email":"foo@bar.idv"}`)
  req3, _ := http.NewRequest("POST", "/add1stuser", bytes.NewBuffer(data1))
  req3.Header.Set("Content-Type", "application/json")
  req3.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.Equal(t, "", res2.Body.String())

  /* Request add1stuser with GET method again */
  res4 := httptest.NewRecorder()
  req4, _ := http.NewRequest("GET", "/add1stuser", nil)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusFound, res4.Code)
  assert.Equal(t, "/", res4.Header().Get("Location"))

  /* Requests with session, the CSRF token, and correct POST form fields again */
  res5 := httptest.NewRecorder()
  data2 := []byte(`{"account":"fooagain","passwd":"bar","email":"fooagain@bar.idv"}`)
  req5, _ := http.NewRequest("POST", "/add1stuser", bytes.NewBuffer(data2))
  req5.Header.Set("Content-Type", "application/json")
  req5.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req5, res2)
  r.ServeHTTP(res5, req5)

  assert.Equal(t, http.StatusForbidden, res5.Code)
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

func TestLogin(t *testing.T) {
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
  req3, _ := http.NewRequest("GET", "/", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  expected_index := "<h1>Welcome foo</h1>\n"
  assert.Equal(t, http.StatusOK, res3.Code)
  assert.True(t, strings.Contains(res3.Body.String(), expected_index))
}

func TestAddUserWrong(t *testing.T) {
  r := setupRouter()

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  csrf_token := getCSRFToken(res1)
  assert.True(t, len(csrf_token) > 0)

  res2 := httptest.NewRecorder()
  data1 := url.Values{}
  data1.Set("account", "foo")
  data1.Set("passwd", "bar")
  data1.Set("_csrf", csrf_token)
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data1.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusFound, res2.Code)
  assert.Equal(t, "/", res2.Header().Get("Location"))
  assert.Equal(t, "", res2.Body.String())

  /* Request adduser with GET method to have session and CSRF token */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/adduser", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  expected_adduser := "<h1>Add User</h1>"
  expected_csrf := "name=\"_csrf\""

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.True(t, strings.Contains(res3.Body.String(), expected_adduser))
  assert.True(t, strings.Contains(res3.Body.String(), expected_csrf))

  csrf_token = getCSRFToken(res3)
  assert.True(t, len(csrf_token) > 0)

  /* Requests with session and the CSRF token, but bad POST form fields */
  res4 := httptest.NewRecorder()
  data2 := url.Values{}
  data2.Set("_csrf", csrf_token)
  req4, _ := http.NewRequest("POST", "/adduser", strings.NewReader(data2.Encode()))
  req4.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusBadRequest, res4.Code)
  assert.Equal(t, "Wrong account, password or email address", res4.Body.String())

  res5 := httptest.NewRecorder()
  data2.Set("account", "foo2")
  req5, _ := http.NewRequest("POST", "/adduser", strings.NewReader(data2.Encode()))
  req5.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req5, res2)
  r.ServeHTTP(res5, req5)

  assert.Equal(t, http.StatusBadRequest, res5.Code)
  assert.Equal(t, "Wrong account, password or email address", res5.Body.String())

  res6 := httptest.NewRecorder()
  data2.Set("account", "foo2")
  data2.Set("passwd", "bar2")
  req6, _ := http.NewRequest("POST", "/adduser", strings.NewReader(data2.Encode()))
  req6.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req6, res2)
  r.ServeHTTP(res6, req6)

  assert.Equal(t, http.StatusBadRequest, res6.Code)
  assert.Equal(t, "Wrong account, password or email address", res6.Body.String())
}

func TestAddUserSuccess(t *testing.T) {
  r := setupRouter()
  utils := database.UserUtils{DB: database.ConnectDB("")}

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  csrf_token := getCSRFToken(res1)
  assert.True(t, len(csrf_token) > 0)

  res2 := httptest.NewRecorder()
  data1 := url.Values{}
  data1.Set("account", "foo")
  data1.Set("passwd", "bar")
  data1.Set("_csrf", csrf_token)
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data1.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusFound, res2.Code)
  assert.Equal(t, "/", res2.Header().Get("Location"))
  assert.Equal(t, "", res2.Body.String())

  /* Request adduser with GET method to have session and CSRF token */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/adduser", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  expected_adduser := "<h1>Add User</h1>"
  expected_csrf := "name=\"_csrf\""

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.True(t, strings.Contains(res3.Body.String(), expected_adduser))
  assert.True(t, strings.Contains(res3.Body.String(), expected_csrf))

  csrf_token = getCSRFToken(res3)
  assert.True(t, len(csrf_token) > 0)

  /* Requests with session, the CSRF token, and correct POST form fields to add a guest */
  res4 := httptest.NewRecorder()
  data2 := []byte(`{"account":"foo2","passwd":"bar2","email":"foo2@bar2.idv","groups":["Guest"]}`)
  req4, _ := http.NewRequest("POST", "/adduser", bytes.NewBuffer(data2))
  req4.Header.Set("Content-Type", "application/json")
  req4.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusOK, res4.Code)
  assert.Equal(t, "", res4.Body.String())

  user2, err := utils.Get("foo2", hashSHA512("bar2"))
  assert.Nil(t, err)
  assert.Equal(t, "foo2", user2.Account)
  assert.Equal(t, hashSHA512("bar2"), user2.Passwd)
  assert.Equal(t, "foo2@bar2.idv", user2.Email)
  assert.Equal(t, 1, len(user2.Groups))
  assert.Equal(t, "Guest", user2.Groups[0].Name)

  /* Requests with session, the CSRF token, and correct POST form fields to add another administrator */
  res5 := httptest.NewRecorder()
  data3 := []byte(`{"account":"foo3","passwd":"bar3","email":"foo3@bar3.idv","groups":["Guest","Administrator"]}`)
  req5, _ := http.NewRequest("POST", "/adduser", bytes.NewBuffer(data3))
  req5.Header.Set("Content-Type", "application/json")
  req5.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req5, res2)
  r.ServeHTTP(res5, req5)

  assert.Equal(t, http.StatusOK, res5.Code)
  assert.Equal(t, "", res5.Body.String())

  user3, err := utils.Get("foo3", hashSHA512("bar3"))
  assert.Nil(t, err)
  assert.Equal(t, "foo3", user3.Account)
  assert.Equal(t, hashSHA512("bar3"), user3.Passwd)
  assert.Equal(t, "foo3@bar3.idv", user3.Email)
  assert.Equal(t, "Guest", user3.Groups[0].Name)
  assert.Equal(t, "Administrator", user3.Groups[1].Name)
}

func TestGuestAddUser(t *testing.T) {
  r := setupRouter()

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  csrf_token := getCSRFToken(res1)
  assert.True(t, len(csrf_token) > 0)

  res2 := httptest.NewRecorder()
  data1 := url.Values{}
  data1.Set("account", "foo2")
  data1.Set("passwd", "bar2")
  data1.Set("_csrf", csrf_token)
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data1.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusFound, res2.Code)
  assert.Equal(t, "/", res2.Header().Get("Location"))
  assert.Equal(t, "", res2.Body.String())

  /* Guest requests adduser with GET method */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/adduser", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusForbidden, res3.Code)
  assert.Equal(t, "", res3.Body.String())

  /* Guest requests updateuser with GET method to borrow CSRF token */
  res4 := httptest.NewRecorder()
  req4, _ := http.NewRequest("GET", "/updateuser/2", nil)
  copyCookies(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusOK, res4.Code)
  csrf_token = getCSRFToken(res4)
  assert.True(t, len(csrf_token) > 0)

  /* Guest requests adduser with POST method */
  res5 := httptest.NewRecorder()
  data2 := []byte(`{"account":"foo4","passwd":"bar4","email":"foo4@bar4.idv","groups":["Administrator"]}`)
  req5, _ := http.NewRequest("POST", "/adduser", bytes.NewBuffer(data2))
  req5.Header.Set("Content-Type", "application/json")
  req5.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req5, res4)
  r.ServeHTTP(res5, req5)

  assert.Equal(t, http.StatusBadRequest, res5.Code)
}

func TestGuestUpdateUser(t *testing.T) {
  r := setupRouter()
  utils := database.UserUtils{DB: database.ConnectDB("")}

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  csrf_token := getCSRFToken(res1)
  assert.True(t, len(csrf_token) > 0)

  res2 := httptest.NewRecorder()
  data1 := url.Values{}
  data1.Set("account", "foo2")
  data1.Set("passwd", "bar2")
  data1.Set("_csrf", csrf_token)
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data1.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusFound, res2.Code)
  assert.Equal(t, "/", res2.Header().Get("Location"))
  assert.Equal(t, "", res2.Body.String())

  /* Request updateuser with GET method to have session and CSRF token */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/updateuser/2", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  expected_updateuser := "<h1>Update User</h1>"
  expected_roles := "let available_groups = [\"Guest\"];"
  expected_csrf := "name=\"_csrf\""

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.True(t, strings.Contains(res3.Body.String(), expected_updateuser))
  assert.True(t, strings.Contains(res3.Body.String(), expected_roles))
  assert.True(t, strings.Contains(res3.Body.String(), expected_csrf))

  csrf_token = getCSRFToken(res3)
  assert.True(t, len(csrf_token) > 0)

  /* Requests with session, the CSRF token, and correct POST form fields to update self */
  res4 := httptest.NewRecorder()
  data2 := []byte(`{"email":"foo2@bar.idv","groups":["Administrator"]}`)
  req4, _ := http.NewRequest("POST", "/updateuser/2", bytes.NewBuffer(data2))
  req4.Header.Set("Content-Type", "application/json")
  req4.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusOK, res4.Code)
  assert.Equal(t, "", res4.Body.String())

  user2, err := utils.Get("foo2", hashSHA512("bar2"))
  assert.Nil(t, err)
  assert.Equal(t, "foo2", user2.Account)
  assert.Equal(t, hashSHA512("bar2"), user2.Passwd)
  assert.Equal(t, "foo2@bar.idv", user2.Email)
  assert.Equal(t, 1, len(user2.Groups))
  assert.Equal(t, "Guest", user2.Groups[0].Name)

  /* Request updateuser with GET method to update an administrator */
  res5 := httptest.NewRecorder()
  req5, _ := http.NewRequest("GET", "/updateuser/3", nil)
  copyCookies(req5, res2)
  r.ServeHTTP(res5, req5)

  assert.Equal(t, http.StatusForbidden, res5.Code)
  assert.Equal(t, "", res5.Body.String())

  /* Requests with session, the CSRF token, and correct POST form fields to update an administrator */
  res6 := httptest.NewRecorder()
  req6, _ := http.NewRequest("POST", "/updateuser/3", bytes.NewBuffer(data2))
  req6.Header.Set("Content-Type", "application/json")
  req6.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req6, res2)
  r.ServeHTTP(res6, req6)

  assert.Equal(t, http.StatusForbidden, res6.Code)
  assert.Equal(t, "", res6.Body.String())

  /* Request updateuser with GET method to update an unexisted user */
  res7 := httptest.NewRecorder()
  req7, _ := http.NewRequest("GET", "/updateuser/10", nil)
  copyCookies(req7, res2)
  r.ServeHTTP(res7, req7)

  assert.Equal(t, http.StatusForbidden, res7.Code)
  assert.Equal(t, "", res7.Body.String())

  /* Requests with session, the CSRF token, and correct POST form fields to update an unexisted user */
  res8 := httptest.NewRecorder()
  req8, _ := http.NewRequest("POST", "/updateuser/10", bytes.NewBuffer(data2))
  req8.Header.Set("Content-Type", "application/json")
  req8.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req8, res2)
  r.ServeHTTP(res8, req8)

  assert.Equal(t, http.StatusForbidden, res8.Code)

  /* Request updateuser with GET method to update a wrong user */
  res9 := httptest.NewRecorder()
  req9, _ := http.NewRequest("GET", "/updateuser/ABC10", nil)
  copyCookies(req9, res2)
  r.ServeHTTP(res9, req9)

  assert.Equal(t, http.StatusNotFound, res9.Code)
  assert.Equal(t, "", res9.Body.String())

  /* Requests with session, the CSRF token, and correct POST form fields to update a wrong user */
  res10 := httptest.NewRecorder()
  req10, _ := http.NewRequest("POST", "/updateuser/ABC10", bytes.NewBuffer(data2))
  req10.Header.Set("Content-Type", "application/json")
  req10.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req10, res2)
  r.ServeHTTP(res10, req10)

  assert.Equal(t, http.StatusNotFound, res10.Code)
  assert.Equal(t, "", res10.Body.String())
}

func TestAdministratorUpdateUser(t *testing.T) {
  r := setupRouter()
  utils := database.UserUtils{DB: database.ConnectDB("")}

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  csrf_token := getCSRFToken(res1)
  assert.True(t, len(csrf_token) > 0)

  res2 := httptest.NewRecorder()
  data1 := url.Values{}
  data1.Set("account", "foo")
  data1.Set("passwd", "bar")
  data1.Set("_csrf", csrf_token)
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data1.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusFound, res2.Code)
  assert.Equal(t, "/", res2.Header().Get("Location"))
  assert.Equal(t, "", res2.Body.String())

  /* Request updateuser with GET method to have session and CSRF token */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/updateuser/2", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  expected_updateuser := "<h1>Update User</h1>"
  expected_roles := "let available_groups = [\"Guest\",\"Administrator\""
  expected_csrf := "name=\"_csrf\""

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.True(t, strings.Contains(res3.Body.String(), expected_updateuser))
  assert.True(t, strings.Contains(res3.Body.String(), expected_roles))
  assert.True(t, strings.Contains(res3.Body.String(), expected_csrf))

  csrf_token = getCSRFToken(res3)
  assert.True(t, len(csrf_token) > 0)

  /* Requests with session, the CSRF token, and correct POST form fields to update a Guest */
  res4 := httptest.NewRecorder()
  data2 := []byte(`{"email":"foo2@bar.idvv","groups":["Administrator"]}`)
  req4, _ := http.NewRequest("POST", "/updateuser/2", bytes.NewBuffer(data2))
  req4.Header.Set("Content-Type", "application/json")
  req4.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req4, res2)
  r.ServeHTTP(res4, req4)

  assert.Equal(t, http.StatusOK, res4.Code)
  assert.Equal(t, "", res4.Body.String())

  user2, err := utils.Get("foo2", hashSHA512("bar2"))
  assert.Nil(t, err)
  assert.Equal(t, "foo2", user2.Account)
  assert.Equal(t, hashSHA512("bar2"), user2.Passwd)
  assert.Equal(t, "foo2@bar.idvv", user2.Email)
  assert.Equal(t, "Administrator", user2.Groups[0].Name)

  /* Request updateuser with GET method to update an Administrator*/
  res5 := httptest.NewRecorder()
  req5, _ := http.NewRequest("GET", "/updateuser/2", nil)
  copyCookies(req5, res2)
  r.ServeHTTP(res5, req5)

  assert.Equal(t, http.StatusOK, res5.Code)
  assert.True(t, strings.Contains(res5.Body.String(), expected_updateuser))
  assert.True(t, strings.Contains(res5.Body.String(), expected_roles))
  assert.True(t, strings.Contains(res5.Body.String(), expected_csrf))

  csrf_token = getCSRFToken(res5)
  assert.True(t, len(csrf_token) > 0)

  /* Requests with session, the CSRF token, and correct POST form fields to update an Administrator */
  res6 := httptest.NewRecorder()
  data3 := []byte(`{"email":"foo2@bar.idv","groups":["Guest"]}`)
  req6, _ := http.NewRequest("POST", "/updateuser/2", bytes.NewBuffer(data3))
  req6.Header.Set("Content-Type", "application/json")
  req6.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req6, res2)
  r.ServeHTTP(res6, req6)

  assert.Equal(t, http.StatusOK, res6.Code)
  assert.Equal(t, "", res6.Body.String())

  user3, err := utils.Get("foo2", hashSHA512("bar2"))
  assert.Nil(t, err)
  assert.Equal(t, "foo2", user3.Account)
  assert.Equal(t, hashSHA512("bar2"), user3.Passwd)
  assert.Equal(t, "foo2@bar.idv", user3.Email)
  assert.Equal(t, "Guest", user3.Groups[0].Name)

  /* Request updateuser with GET method to update an unexisted user */
  res7 := httptest.NewRecorder()
  req7, _ := http.NewRequest("GET", "/updateuser/10", nil)
  copyCookies(req7, res2)
  r.ServeHTTP(res7, req7)

  assert.Equal(t, http.StatusNotFound, res7.Code)
  assert.Equal(t, "", res7.Body.String())

  /* Requests with session, the CSRF token, and correct POST form fields to update an unexisted user */
  res8 := httptest.NewRecorder()
  req8, _ := http.NewRequest("POST", "/updateuser/10", bytes.NewBuffer(data3))
  req8.Header.Set("Content-Type", "application/json")
  req8.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req8, res2)
  r.ServeHTTP(res8, req8)

  assert.Equal(t, http.StatusNotFound, res8.Code)
  assert.Equal(t, "", res8.Body.String())

  /* Request updateuser with GET method to update a wrong user */
  res9 := httptest.NewRecorder()
  req9, _ := http.NewRequest("GET", "/updateuser/ABC10", nil)
  copyCookies(req9, res2)
  r.ServeHTTP(res9, req9)

  assert.Equal(t, http.StatusNotFound, res9.Code)
  assert.Equal(t, "", res9.Body.String())

  /* Requests with session, the CSRF token, and correct POST form fields to update a wrong user */
  res10 := httptest.NewRecorder()
  req10, _ := http.NewRequest("POST", "/updateuser/ABC10", bytes.NewBuffer(data3))
  req10.Header.Set("Content-Type", "application/json")
  req10.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req10, res2)
  r.ServeHTTP(res10, req10)

  assert.Equal(t, http.StatusNotFound, res10.Code)
  assert.Equal(t, "", res10.Body.String())
}

func TestLogout(t *testing.T) {
  r := setupRouter()

  /* Have the session and the CSRF token for following POST request */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)
  csrf_token := getCSRFToken(res1)
  assert.True(t, len(csrf_token) > 0)

  res2 := httptest.NewRecorder()
  data1 := url.Values{}
  data1.Set("account", "foo")
  data1.Set("passwd", "bar")
  data1.Set("_csrf", csrf_token)
  req2, _ := http.NewRequest("POST", "/login", strings.NewReader(data1.Encode()))
  req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusFound, res2.Code)
  assert.Equal(t, "/", res2.Header().Get("Location"))
  assert.Equal(t, "", res2.Body.String())

  /* Logout */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/logout", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusFound, res3.Code)
  assert.Equal(t, "/", res3.Header().Get("Location"))

  /* Request index again */
  res4 := httptest.NewRecorder()
  req4, _ := http.NewRequest("GET", "/", nil)
  copyCookies(req4, res3)
  r.ServeHTTP(res4, req4)

  expected_index := "<h1>Welcome</h1>\n"

  assert.Equal(t, http.StatusOK, res4.Code)
  assert.True(t, strings.Contains(res4.Body.String(), expected_index))
}
