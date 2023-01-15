package main

import (
  "fmt"
  "time"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
)

func index_guest(c *gin.Context) {
  c.File("./resource/index_guest.html")
}

func Index(c *gin.Context) {
  session := sessions.Default(c)

  account := session.Get("account")
  if (account == nil) {
    index_guest(c)
    return
  }

  c.HTML(http.StatusOK, "index.tmpl", gin.H{
    "account": account,
  })
}

func Ping (c *gin.Context) {
  c.String(http.StatusOK, "pong")
}

func Login (c *gin.Context) {
  account := c.PostForm("account")
  passwd := c.PostForm("passwd")

  if (account == "" || passwd == "") {
    c.String(http.StatusForbidden, "Wrong account or password")
    c.Abort()
    return
  }

  session := sessions.Default(c)
  session.Set("account", account)
  session.Save()

  c.Status(http.StatusOK)
}

func Showdate(c *gin.Context) {
  session := sessions.Default(c)
  account := session.Get("account")
  currentTime := time.Now()

  res := fmt.Sprintf("Welcome %s %s", account, currentTime.String())
  c.String(http.StatusOK, res)
}
