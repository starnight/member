package middleware

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"

  "strings"
  "testing"
  "net/http/httptest"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
)

func PublicRoutes (g *gin.RouterGroup) {
  g.POST("/login", func (c *gin.Context) {
    session := sessions.Default(c)
    session.Set("account", "foo")
    session.Set("id", uint(1))
    session.Save()
    c.Status(http.StatusOK)
  })
}

func PrivateRoutes (g *gin.RouterGroup) {
  g.GET("/private", func (c *gin.Context){
    session := sessions.Default(c)
    account := session.Get("account").(string)
    c.String(http.StatusOK, account)
  })
}

func copyCookies(req *http.Request, res *httptest.ResponseRecorder) {
  req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func TestAuthenticationRequired(t *testing.T) {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

  public := r.Group("/")
  PublicRoutes(public)

  private := r.Group("/")
  private.Use(AuthenticationRequired)
  PrivateRoutes(private)

  /* Must get forbidden, because has not login */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/private", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusForbidden, res1.Code)
  assert.Equal(t, "Please login first", res1.Body.String())

  /* Login */
  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("POST", "/login", nil)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusOK, res2.Code)

  /* Then, request again with the cookie for the session */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/private", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.Equal(t, "foo", res3.Body.String())
}

type mock_userUtils struct {
  mock.Mock
}

func (m_utils *mock_userUtils) IsInGroups(id uint, group_names []string) bool {
  args := m_utils.Called(id, group_names)
  return args.Bool(0)
}

func TestAuthorizationCheckOk(t *testing.T) {
  m_utils := mock_userUtils{}
  grp_names := []string{"Administrator"}

  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

  public := r.Group("/")
  PublicRoutes(public)

  r.GET("/adduser", func(c *gin.Context) {
    authorizationCheck(c, &m_utils)
    if (c.IsAborted()) {
      return
    }
    c.Status(http.StatusOK)
  })

  /* Login */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("POST", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)

  /* Add an user successfully */
  m_utils.On("IsInGroups", uint(1), grp_names).Return(true)
  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("GET", "/adduser", nil)
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusOK, res2.Code)
  m_utils.AssertCalled(t, "IsInGroups", uint(1), grp_names)
}

func TestAuthorizationCheckFailed(t *testing.T) {
  m_utils := mock_userUtils{}
  grp_names := []string{"Administrator"}

  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

  public := r.Group("/")
  PublicRoutes(public)

  r.GET("/adduser", func(c *gin.Context) {
    authorizationCheck(c, &m_utils)
    if (c.IsAborted()) {
      return
    }
    c.Status(http.StatusOK)
  })

  /* Login */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("POST", "/login", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusOK, res1.Code)

  /* Add an user failed */
  m_utils.On("IsInGroups", uint(1), grp_names).Return(false)
  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("GET", "/adduser", nil)
  copyCookies(req2, res1)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusForbidden, res2.Code)
  m_utils.AssertCalled(t, "IsInGroups", uint(1), grp_names)
}
