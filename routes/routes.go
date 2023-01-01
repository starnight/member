package routes

import (
  "github.com/gin-gonic/gin"

  "github.com/starnight/member/controllers"
)

func PublicRoutes (g *gin.RouterGroup) {
  g.GET("/", controllers.Index)
  g.GET("/ping", controllers.Ping)
  g.POST("/login", controllers.Login)
}

func PrivateRoutes (g *gin.RouterGroup) {
  g.GET("/showdate", controllers.Showdate)
}
