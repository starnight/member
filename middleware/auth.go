package middleware

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
)

func AuthenticationRequired(c *gin.Context) {
  session := sessions.Default(c)

  account := session.Get("account")
  if (account == nil) {
    c.String(http.StatusForbidden, "Please login first")
    c.Abort()
    return
  }

  c.Next()
}
