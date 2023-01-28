package main

import (
  "fmt"
  "strconv"
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

func UpdateUserHTML(c *gin.Context) {
  tg_id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
  if (err != nil) {
    c.Status(http.StatusNotFound)
    c.Abort()
    return
  }
  tg_id := uint(tg_id64)

  utils := database.UserUtils{DB: database.ConnectDB("")}
  session := sessions.Default(c)
  cur_id := session.Get("id").(uint)
  cur_user, _ := utils.GetById(cur_id)

  /* Only:
   * - Userself can update self
   * - Administrator can update others
   */
  if (tg_id != cur_id && cur_user.Role != database.Administrator) {
    c.Status(http.StatusForbidden)
    c.Abort()
    return
  }

  tg_user, err2 := utils.GetById(tg_id)
  if (err2 != nil) {
    c.Status(http.StatusNotFound)
    c.Abort()
    return
  }

  roles := make(map[string]uint32)
  roles["Guest"] = database.Guest
  if (cur_user.Role == database.Administrator) {
    roles["Administrator"] = database.Administrator
  }

  c.HTML(http.StatusOK, "updateuser.tmpl", gin.H{
    "tg_id": tg_id,
    "tg_account": tg_user.Account,
    "tg_email": tg_user.Email,
    "tg_role": tg_user.Role,
    "roles": roles,
    "_csrf": csrf.GetToken(c),
  })
}

func UpdateUser(c *gin.Context) {
  tg_id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
  if (err != nil) {
    c.Status(http.StatusNotFound)
    c.Abort()
    return
  }
  tg_id := uint(tg_id64)

  utils := database.UserUtils{DB: database.ConnectDB("")}
  session := sessions.Default(c)
  cur_id := session.Get("id").(uint)
  cur_user, _ := utils.GetById(cur_id)

  /* Only:
   * - Userself can update self
   * - Administrator can update others
   */
  if (tg_id != cur_id && cur_user.Role != database.Administrator) {
    c.Status(http.StatusForbidden)
    c.Abort()
    return
  }

  tg_user, err2 := utils.GetById(tg_id)
  if (err2 != nil) {
    c.Status(http.StatusNotFound)
    c.Abort()
    return
  }

  tg_email := c.PostForm("email")
  tg_role := c.PostForm("role")

  needupdate := false

  if (tg_email != "") {
    tg_user.Email = tg_email
    needupdate = true
  }

  if (tg_role != "" && cur_user.Role == database.Administrator) {
    switch tg_role {
      case "Administrator":
        tg_user.Role = database.Administrator
	needupdate = true
      case "Guest":
        tg_user.Role = database.Guest
	needupdate = true
    }
  }

  if (needupdate) {
    utils.Update(&tg_user)
  }
}

func Showdate(c *gin.Context) {
  session := sessions.Default(c)
  account := session.Get("account")
  currentTime := time.Now()

  res := fmt.Sprintf("Welcome %s %s", account, currentTime.String())
  c.String(http.StatusOK, res)
}
