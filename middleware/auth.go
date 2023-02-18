package middleware

import (
  "strings"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/starnight/member/database"
)

func AuthenticationRequired(c *gin.Context) {
  session := sessions.Default(c)

  account := session.Get("account")
  if (account == nil) {
    c.String(http.StatusForbidden, "Please login first")
    c.Abort()
    return
  }

  c.Next()
}

func IsInGroups(id uint, group_names []string) bool {
  user_utils := database.UserUtils{DB: database.ConnectDB("")}
  user, _ := user_utils.GetById(id)

  for _, user_group := range user.Groups {
    for _, group_name := range group_names {
      if (user_group.Name == group_name) {
        return true
      }
    }
  }

  return false
}

func AuthorizationRequired(c *gin.Context) {
  session := sessions.Default(c)
  id := session.Get("id").(uint)

  uri_grps := map[string][]string{
    "/adduser": []string{"Administrator"},
  }

  for uri, grps := range uri_grps {
    if (strings.HasPrefix(c.Request.URL.Path, uri)) {
      if (!IsInGroups(id, grps)) {
        c.Status(http.StatusForbidden)
        c.Abort()
        return
      }
      break
    }
  }
}
