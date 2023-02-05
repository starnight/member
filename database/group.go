package database

import (
  "gorm.io/gorm"
)

type Group struct {
  gorm.Model
  Name string `gorm:"unique;not null"`
}

type GroupUtils struct {
  DB *gorm.DB
}

func (utils *GroupUtils) Create(name string) (Group, error) {
  group := Group{
    Name: name,
  }
  res := utils.DB.Create(&group)
  return group, res.Error
}

func (utils *GroupUtils) Get(name string) (Group, error) {
  var group Group
  res := utils.DB.Where("Name = ?", name).First(&group)
  return group, res.Error
}

func (utils *GroupUtils) CreateGroupTables() {
  utils.DB.AutoMigrate(&Group{})
}

func (utils *GroupUtils) FeedGroupTables() {
  init_groups := []string{"Guest", "Administrator"}
  for _, name := range init_groups {
    utils.Create(name)
  }
}

func InitGroupTables(db *gorm.DB) {
  utils := GroupUtils{DB: db}
  utils.CreateGroupTables()
  utils.FeedGroupTables()
}
