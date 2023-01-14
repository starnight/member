package database

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo"
  passwd := "bar"

  user, err := utils.Add(account, passwd)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
}

func TestGetUserSuccess(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo"
  passwd := "bar"

  user, err := utils.Get(account, passwd)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
}

func TestGetUserNone(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo1"
  passwd := "bar"

  _, err := utils.Get(account, passwd)

  assert.NotNil(t, err)
}

func TestCountUser(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}

  var expected_cnt int64 = 1

  cnt, err := utils.Count()

  assert.Nil(t, err)
  assert.Equal(t, cnt, expected_cnt)
}
