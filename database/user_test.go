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
  role := Administrator

  user, err := utils.Add(account, passwd, email, role)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
  assert.Equal(t, user.Role, role)

}

func TestGetUserSuccess(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo"
  passwd := "bar"
  email := "foo@bar.idv"
  role := Administrator

  user, err := utils.Get(account, passwd)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
  assert.Equal(t, user.Role, role)
}

func TestGetUserNone(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo1"
  passwd := "bar"

  _, err := utils.Get(account, passwd)

  assert.NotNil(t, err)
}

func TestGetUserById(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  account := "foo"
  passwd := "bar"
  email := "foo@bar.idv"
  role := Administrator

  user, err := utils.Get(account, passwd)

  assert.Nil(t, err)

  user, err = utils.GetById(user.ID)
  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
  assert.Equal(t, user.Role, role)
}

func TestUpdateUser(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}
  expected_id := uint(1)
  expected_account := "foo"
  expected_passwd := "bar"
  expected_email := "foo@bar.idv2"
  expected_role := Administrator

  user, err := utils.GetById(expected_id)
  assert.Nil(t, err)

  user.Email = expected_email

  utils.Update(&user)

  assert.Equal(t, user.ID, expected_id)
  assert.Equal(t, user.Account, expected_account)
  assert.Equal(t, user.Passwd, expected_passwd)
  assert.Equal(t, user.Email, expected_email)
  assert.Equal(t, user.Role, expected_role)
}

func TestCountUser(t *testing.T) {
  utils := UserUtils{DB: ConnectDB(GetDBStr("test"))}

  var expected_cnt int64 = 1

  cnt, err := utils.Count()

  assert.Nil(t, err)
  assert.Equal(t, cnt, expected_cnt)
}
