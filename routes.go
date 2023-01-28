package main

import (
  "github.com/gin-gonic/gin"
)

func PublicRoutes (g *gin.RouterGroup) {
  g.GET("/", Index)
  g.GET("/ping", Ping)
  g.GET("/login", LoginHTML)
  g.POST("/login", Login)
  g.GET("/add1stuser", Add1stUserHTML)
  g.POST("/add1stuser", Add1stUser)
}

func PrivateRoutes (g *gin.RouterGroup) {
  g.GET("/adduser", AddUserHTML)
  g.POST("/adduser", AddUser)
  g.GET("/updateuser/:id", UpdateUserHTML)
  g.POST("/updateuser/:id", UpdateUser)
  g.GET("/showdate", Showdate)
}
