package database

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo"
  passwd := "bar"
  email := "foo@bar.idv"

  user, err := utils.Add(account, passwd, email)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
}

func TestGetUserSuccess(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo"
  passwd := "bar"
  email := "foo@bar.idv"

  user, err := utils.Get(account, passwd)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
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
