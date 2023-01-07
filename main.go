package main

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/utrack/gin-csrf"

  "github.com/starnight/member/middleware"
  "github.com/starnight/member/routes"
)

func setupRouter() *gin.Engine {
  r := gin.Default()

  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))
  r.Use(csrf.Middleware(csrf.Options{
    Secret: "secret123",
    ErrorFunc: func(c *gin.Context) {
      c.String(400, "CSRF token mismatch")
      c.Abort()
    },
  }))
  r.Use(middleware.AddCSRFToken)

  public := r.Group("/")
  routes.PublicRoutes(public)

  private := r.Group("/")
  private.Use(middleware.AuthRequired)
  routes.PrivateRoutes(private)

  return r
}

func main() {
  r := setupRouter()
  r.Run(":8080")
}
