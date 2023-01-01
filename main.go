package main

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"

  "github.com/starnight/member/middleware"
  "github.com/starnight/member/routes"
)

func setupRouter() *gin.Engine {
  r := gin.Default()

  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

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
