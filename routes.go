package main

import (
  "github.com/gin-gonic/gin"
)

func PublicRoutes (g *gin.RouterGroup) {
  g.GET("/", Index)
  g.GET("/ping", Ping)
  g.GET("/login", LoginHTML)
  g.POST("/login", Login)
  g.GET("/adduser", AddUserHTML)
  g.POST("/adduser", AddUser)
}

func PrivateRoutes (g *gin.RouterGroup) {
  g.GET("/showdate", Showdate)
}
