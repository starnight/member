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

var user_utils = database.UserUtils{DB: database.ConnectDB("")}
var group_utils = database.GroupUtils{DB: database.ConnectDB("")}

func hashSHA512(text string) string {
  txtByte := []byte(text)
  sha512 := sha512.New()

  sha512.Write(txtByte)

  hashedtxt := sha512.Sum(nil)

  return hex.EncodeToString(hashedtxt)
}

func IsAuthorized(id uint, group_name string) bool {
  group_names := []string{group_name}

  return user_utils.IsInGroups(id, group_names)
}

func index_guest(c *gin.Context) {
  c.File("./resource/index_guest.html")
}

func Index(c *gin.Context) {
  cnt, _ := user_utils.Count()
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
  var group database.Group

  if (c.Request.URL.Path == "/add1stuser") {
    /* The first user should be an Administrator for following management */
    group, _ = group_utils.Get("Administrator")
  } else {
    switch c.PostForm("role") {
      case "Administrator":
        group, _ = group_utils.Get("Administrator")
      default:
        group, _ = group_utils.Get("Guest")
    }
  }

  groups := []database.Group{group}

  if (account == "" || passwd == "" || email =="") {
    c.String(http.StatusBadRequest, "Wrong account, password or email address")
    c.Abort()
    return
  }

  user_utils.Add(account, hashSHA512(passwd), email, groups)

  c.Status(http.StatusOK)
}

func Add1stUserHTML(c *gin.Context) {
  cnt, _ := user_utils.Count()
  if (cnt != 0) {
    c.Redirect(http.StatusFound, "/")
    return
  }

  AddUserHTML(c)
}

func Add1stUser(c *gin.Context) {
  cnt, _ := user_utils.Count()
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

  user, err := user_utils.Get(account, hashSHA512(passwd))

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

  session := sessions.Default(c)
  cur_id := session.Get("id").(uint)

  /* Only:
   * - Userself can update self
   * - Administrator can update others
   */
  if (tg_id != cur_id && !IsAuthorized(cur_id, "Administrator")) {
    c.Status(http.StatusForbidden)
    c.Abort()
    return
  }

  tg_user, err2 := user_utils.GetById(tg_id)
  if (err2 != nil) {
    c.Status(http.StatusNotFound)
    c.Abort()
    return
  }

  roles := make(map[string]uint)
  gst_grp, _ := group_utils.Get("Guest")
  roles["Guest"] = gst_grp.ID
  if (IsAuthorized(cur_id, "Administrator")) {
    adm_grp, _ := group_utils.Get("Administrator")
    roles["Administrator"] = adm_grp.ID
  }

  c.HTML(http.StatusOK, "updateuser.tmpl", gin.H{
    "tg_id": tg_id,
    "tg_account": tg_user.Account,
    "tg_email": tg_user.Email,
    "tg_role": tg_user.Groups[0].ID,
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

  session := sessions.Default(c)
  cur_id := session.Get("id").(uint)

  /* Only:
   * - Userself can update self
   * - Administrator can update others
   */
  if (tg_id != cur_id && !IsAuthorized(cur_id, "Administrator")) {
    c.Status(http.StatusForbidden)
    c.Abort()
    return
  }

  tg_user, err2 := user_utils.GetById(tg_id)
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

  if (tg_role != "" && IsAuthorized(cur_id, "Administrator")) {
    switch tg_role {
      case "Administrator":
        grp, _ := group_utils.Get("Administrator")
        tg_user.Groups = []database.Group{grp}
	needupdate = true
      case "Guest":
        grp, _ := group_utils.Get("Guest")
        tg_user.Groups = []database.Group{grp}
	needupdate = true
    }
  }

  if (needupdate) {
    user_utils.Update(&tg_user)
  }
}

func Showdate(c *gin.Context) {
  session := sessions.Default(c)
  account := session.Get("account")
  currentTime := time.Now()

  res := fmt.Sprintf("Welcome %s %s", account, currentTime.String())
  c.String(http.StatusOK, res)
}

func Logout(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()
  session.Options(sessions.Options{MaxAge: -1})
  session.Save()
  c.Redirect(http.StatusFound, "/")
}
