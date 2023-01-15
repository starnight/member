package main

import (
  "github.com/gin-gonic/gin"
)

func PublicRoutes (g *gin.RouterGroup) {
  g.GET("/", Index)
  g.GET("/ping", Ping)
  g.GET("/login", LoginHTML)
  g.POST("/login", Login)
}

func PrivateRoutes (g *gin.RouterGroup) {
  g.GET("/showdate", Showdate)
}
