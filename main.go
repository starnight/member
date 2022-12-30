package main

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
)

func setupRouter() *gin.Engine {
  r := gin.Default()

  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

  r.GET("/", index)
  r.GET("/ping", ping)
  r.POST("/login", login)

  return r
}

func main() {
  r := setupRouter()
  r.Run(":8080")
}
