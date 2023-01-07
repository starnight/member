package middleware

import (
  "github.com/gin-gonic/gin"
  "github.com/utrack/gin-csrf"
)

func AddCSRFToken(c *gin.Context) {
  c.Writer.Header().Set("X-CSRF-TOKEN", csrf.GetToken(c))
  c.Next()
}
