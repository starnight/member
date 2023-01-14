package database

import (
  "fmt"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

var _db *gorm.DB = nil

func GetDBStr(mode string) string {
  if (mode == "test") {
    return "test.db"
  }

  return "product.db"
}

func ConnectDB(name string) *gorm.DB {
  if (_db != nil) {
    return _db
  }

  var err error

  _db, err = gorm.Open(sqlite.Open(name), &gorm.Config{})
  if err != nil {
    err_str := fmt.Sprintf("failed to connect database %s\n", name)
    panic(err_str)
  }

  return _db
}

func CreateTables(db *gorm.DB) {
  utils := UserUtils{DB: db}
  utils.CreateUserTables()
}
