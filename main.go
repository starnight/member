package main

import (
  "fmt"
  "os"

  "github.com/alexflint/go-arg"

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
  private.Use(middleware.AuthenticationRequired)
  private.Use(middleware.AuthorizationRequired)
  PrivateRoutes(private)

  return r
}

func dealAddr(port uint) string {
  var addrStr string

  if (port > 0) {
    addrStr = fmt.Sprintf(":%d", port)
  } else if (os.Getenv("PORT") != "") {
    addrStr = ":" + os.Getenv("PORT")
  } else {
    addrStr = ":8080"
  }

  return addrStr
}

type configSet struct {
  AddrStr string
}

func parseArgs() configSet {
  cfg := configSet{}
  var args struct {
    Port uint `help:"listening port"`
  }

  arg.MustParse(&args)

  cfg.AddrStr = dealAddr(args.Port)

  return cfg
}

func main() {
  cfg := parseArgs()

  setupDB()

  r := setupRouter()
  r.Run(cfg.AddrStr)
}
