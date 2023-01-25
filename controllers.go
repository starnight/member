package main

import (
  "fmt"
  "time"
  "crypto/sha512"
  "encoding/hex"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/utrack/gin-csrf"
  "github.com/starnight/member/database"
)

func hashSHA512(text string) string {
  txtByte := []byte(text)
  sha512 := sha512.New()

  sha512.Write(txtByte)

  hashedtxt := sha512.Sum(nil)

  return hex.EncodeToString(hashedtxt)
}

func IsAuthorized(id uint, role uint32) bool {
  utils := database.UserUtils{DB: database.ConnectDB("")}
  cur_user, _ := utils.GetById(id)
  res := false

  if (cur_user.Role == role) {
    res = true
  }

  return res
}

func index_guest(c *gin.Context) {
  c.File("./resource/index_guest.html")
}

func Index(c *gin.Context) {
  utils := database.UserUtils{DB: database.ConnectDB("")}
  cnt, _ := utils.Count()
  if (cnt == 0) {
    c.Redirect(http.StatusFound, "/add1stuser")
    return
  }

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

func AddUserHTML(c *gin.Context) {
  add1stuser := false

  if (c.Request.URL.Path == "/add1stuser") {
    add1stuser = true
  } else {
    /* Only Administrator can add users.  Others are forbidden */
    session := sessions.Default(c)
    id := session.Get("id").(uint)
    if (!IsAuthorized(id, database.Administrator)) {
      c.Status(http.StatusForbidden)
      c.Abort()
      return
    }
  }

  c.HTML(http.StatusOK, "adduser.tmpl", gin.H{
    "IsAdd1stUser": add1stuser,
    "_csrf": csrf.GetToken(c),
  })
}

func AddUser(c *gin.Context) {
  account := c.PostForm("account")
  passwd := c.PostForm("passwd")
  email := c.PostForm("email")
  role := database.Guest

  if (c.Request.URL.Path == "/add1stuser") {
    /* The first user should be an Administrator for following management */
    role = database.Administrator
  } else {
    /* Only Administrator can add users.  Others are forbidden */
    session := sessions.Default(c)
    id := session.Get("id").(uint)
    if (!IsAuthorized(id, database.Administrator)) {
      c.Status(http.StatusForbidden)
      c.Abort()
      return
    }

    switch c.PostForm("role") {
      case "Administrator":
        role = database.Administrator
    }
  }

  if (account == "" || passwd == "" || email =="") {
    c.String(http.StatusBadRequest, "Wrong account, password or email address")
    c.Abort()
    return
  }

  utils := database.UserUtils{DB: database.ConnectDB("")}
  utils.Add(account, hashSHA512(passwd), email, role)

  c.Status(http.StatusOK)
}

func Add1stUserHTML(c *gin.Context) {
  utils := database.UserUtils{DB: database.ConnectDB("")}
  cnt, _ := utils.Count()
  if (cnt != 0) {
    c.Redirect(http.StatusFound, "/")
    return
  }

  AddUserHTML(c)
}

func Add1stUser(c *gin.Context) {
  utils := database.UserUtils{DB: database.ConnectDB("")}
  cnt, _ := utils.Count()
  if (cnt != 0) {
    c.Status(http.StatusForbidden)
    return
  }

  AddUser(c)
}

func LoginHTML(c *gin.Context) {
  c.HTML(http.StatusOK, "login.tmpl", gin.H{
    "_csrf": csrf.GetToken(c),
  })
}

func Login (c *gin.Context) {
  account := c.PostForm("account")
  passwd := c.PostForm("passwd")

  utils := database.UserUtils{DB: database.ConnectDB("")}
  user, err := utils.Get(account, hashSHA512(passwd))

  if (err != nil) {
    c.String(http.StatusForbidden, "Wrong account or password")
    c.Abort()
    return
  }

  session := sessions.Default(c)
  session.Set("id", user.ID)
  session.Set("account", user.Account)
  session.Save()

  c.Redirect(http.StatusFound, "/")
}

func Showdate(c *gin.Context) {
  session := sessions.Default(c)
  account := session.Get("account")
  currentTime := time.Now()

  res := fmt.Sprintf("Welcome %s %s", account, currentTime.String())
  c.String(http.StatusOK, res)
}
