package main

import (
  "fmt"

  "github.com/alexflint/go-arg"

  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/utrack/gin-csrf"

  "github.com/starnight/member/middleware"
  "github.com/starnight/member/database"
)

type configSet struct {
  AddrStr string
}

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
  return fmt.Sprintf(":%d", port)
}

func parseArgs() configSet {
  cfg := configSet{}
  var args struct {
    Port uint `arg:"env:PORT" default:"8080" help:"listening port"`
  }

  arg.MustParse(&args)

  cfg.AddrStr = dealAddr(args.Port)

  return cfg
}

type IGinEngine interface {
  /* Follow gin's Engine format by each version */
  /* Run() https://github.com/gin-gonic/gin/blob/v1.9.0/gin.go#L376 */
  Run(addr ...string) (err error)
}

func doRun(r IGinEngine, cfg configSet) {
  r.Run(cfg.AddrStr)
}

func main() {
  cfg := parseArgs()

  setupDB()

  r := setupRouter()
  doRun(r, cfg)
}
