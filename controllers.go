package main

import (
  "strconv"
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

type addUser struct {
  Account string `json:"account"`
  Passwd string `json:"passwd"`
  Email string `json:"email"`
  Grp_names []string `json:"groups"`
}

func AddUser(c *gin.Context) {
  var req_addUser addUser
  var group database.Group
  var grp_names []string

  c.BindJSON(&req_addUser)

  if (req_addUser.Account == "" || req_addUser.Passwd == "" || req_addUser.Email == "") {
    c.String(http.StatusBadRequest, "Wrong account, password or email address")
    c.Abort()
    return
  }

  if (c.Request.URL.Path == "/add1stuser") {
    grp_names = []string{"Administrator", "Guest"}
  } else {
    grp_names = req_addUser.Grp_names
  }

  groups := []database.Group{}
  for _, grp_name := range grp_names {
    group, _ = group_utils.Get(grp_name)
    groups = append(groups, group)
  }

  user_utils.Add(req_addUser.Account, hashSHA512(req_addUser.Passwd), req_addUser.Email, groups)

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

  tg_grp_names := []string{}
  for _, grp := range tg_user.Groups {
    tg_grp_names = append(tg_grp_names, grp.Name)
  }

  grp_names := []string{"Guest"}
  if (IsAuthorized(cur_id, "Administrator")) {
    grp_names = append(grp_names, "Administrator")
  }

  c.HTML(http.StatusOK, "updateuser.tmpl", gin.H{
    "tg_id": tg_id,
    "tg_account": tg_user.Account,
    "tg_email": tg_user.Email,
    "tg_groups": tg_grp_names,
    "groups": grp_names,
    "_csrf": csrf.GetToken(c),
  })
}

type updateUser struct {
  Email string `json:"email"`
  Grp_names []string `json:"groups"`
}

func UpdateUser(c *gin.Context) {
  var req_updateUser updateUser

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

  c.BindJSON(&req_updateUser)

  needupdate := false

  if (req_updateUser.Email != "") {
    tg_user.Email = req_updateUser.Email
    needupdate = true
  }

  grp_names := req_updateUser.Grp_names
  if (len(grp_names) != 0 && IsAuthorized(cur_id, "Administrator")) {
    groups := []database.Group{}
    for _, grp_name := range grp_names {
      group, err_grp := group_utils.Get(grp_name)
      if (err_grp != nil) {
        continue
      }
      groups = append(groups, group)
      tg_user.Groups = groups
      needupdate = true
    }
  }

  if (needupdate) {
    user_utils.Update(&tg_user)
  }
}

func Logout(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()
  session.Options(sessions.Options{MaxAge: -1})
  session.Save()
  c.Redirect(http.StatusFound, "/")
}
