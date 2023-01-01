package main

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"

  "github.com/starnight/member/controllers"
)

func setupRouter() *gin.Engine {
  r := gin.Default()

  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

  r.GET("/", controllers.Index)
  r.GET("/ping", controllers.Ping)
  r.POST("/login", controllers.Login)
  r.GET("/showdate", controllers.Showdate)

  return r
}

func main() {
  r := setupRouter()
  r.Run(":8080")
}
