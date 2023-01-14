package database

import (
  "gorm.io/gorm"
)

type User struct {
  gorm.Model
  Account string `gorm:"unique;not null"`
  Passwd string `gorm:"not null"`
}

type UserUtils struct {
  DB *gorm.DB
}

func (utils *UserUtils) Add(account string, passwd string) (User, error) {
  user := User{Account: account, Passwd: passwd}
  res := utils.DB.Create(&user)
  return user, res.Error
}

func (utils *UserUtils) Get(account string, passwd string) (User, error) {
  var user User
  res := utils.DB.Where("Account = ? AND passwd = ?", account, passwd).First(&user)
  return user, res.Error
}

func (utils *UserUtils) Count() (int64, error) {
  var cnt int64 = 0
  res := utils.DB.Model(&User{}).Count(&cnt)
  return cnt, res.Error
}

func (utils *UserUtils) CreateUserTables() {
  utils.DB.AutoMigrate(&User{})
}
