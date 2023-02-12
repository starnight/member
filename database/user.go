package database

import (
  "gorm.io/gorm"
)

type User struct {
  gorm.Model
  Account string `gorm:"unique;not null"`
  Passwd string `gorm:"not null"`
  Email string `gorm:"unique;not null"`
  Groups []Group `gorm:"many2many:user_groups;"`
}

type UserUtils struct {
  DB *gorm.DB
}

func (utils *UserUtils) Add(account string, passwd string, email string, groups []Group) (User, error) {
  user := User{
    Account: account,
    Passwd: passwd,
    Email: email,
    Groups: groups,
  }
  res := utils.DB.Create(&user)
  return user, res.Error
}

func (utils *UserUtils) Get(account string, passwd string) (User, error) {
  var user User
  res := utils.DB.Where("Account = ? AND passwd = ?", account, passwd).Preload("Groups").First(&user)
  return user, res.Error
}

func (utils *UserUtils) GetById(id uint) (User, error) {
  var user User
  res := utils.DB.Where("ID = ?", id).Preload("Groups").First(&user)
  return user, res.Error
}

func (utils *UserUtils) Update(user *User) {
  utils.DB.Save(user)
  utils.DB.Model(user).Association("Groups").Replace(user.Groups)
}

func (utils *UserUtils) Count() (int64, error) {
  var cnt int64 = 0
  res := utils.DB.Model(&User{}).Count(&cnt)
  return cnt, res.Error
}

func (utils *UserUtils) CreateUserTables() {
  utils.DB.AutoMigrate(&User{})
}

func InitUserTables(db *gorm.DB) {
  utils := UserUtils{DB: db}
  utils.CreateUserTables()
}
