package main

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/utrack/gin-csrf"

  "github.com/starnight/member/middleware"
  "github.com/starnight/member/database"
)

func setupDB() {
  dbstr := database.GetDBStr(gin.Mode())
  db := database.ConnectDB(dbstr)
  utils := database.UserUtils{DB: db}

  _, err := utils.Count()

  if (err != nil) {
    database.InitTables(db)
  }
}

func setupRouter() *gin.Engine {
  r := gin.Default()

  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))
  r.Use(csrf.Middleware(csrf.Options{
    Secret: "secret123",
    ErrorFunc: func(c *gin.Context) {
      c.String(400, "CSRF token mismatch")
      c.Abort()
    },
  }))

  r.LoadHTMLGlob("template/*.tmpl")

  public := r.Group("/")
  PublicRoutes(public)

  private := r.Group("/")
  private.Use(middleware.AuthRequired)
  PrivateRoutes(private)

  return r
}

func main() {
  setupDB()

  r := setupRouter()
  r.Run(":8080")
}
