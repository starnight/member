package database

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

var user_utils = UserUtils{DB: ConnectDB(GetDBStr("test"))}
var group_utils = GroupUtils{DB: ConnectDB(GetDBStr("test"))}

func TestAddUser(t *testing.T) {
  account := "foo"
  passwd := "bar"
  email := "foo@bar.idv"
  group, _ := group_utils.Get("Administrator")
  groups := []Group{group}

  user, err := user_utils.Add(account, passwd, email, groups)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
  assert.Equal(t, user.Groups[0], group)
}

func TestGetUserSuccess(t *testing.T) {
  account := "foo"
  passwd := "bar"
  email := "foo@bar.idv"
  group, _ := group_utils.Get("Administrator")

  user, err := user_utils.Get(account, passwd)

  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
  assert.Equal(t, user.Groups[0], group)
}

func TestGetUserNone(t *testing.T) {
  account := "foo1"
  passwd := "bar"

  _, err := user_utils.Get(account, passwd)

  assert.NotNil(t, err)
}

func TestGetUserById(t *testing.T) {
  account := "foo"
  passwd := "bar"
  email := "foo@bar.idv"
  group, _ := group_utils.Get("Administrator")

  user, err := user_utils.Get(account, passwd)

  assert.Nil(t, err)

  user, err = user_utils.GetById(user.ID)
  assert.Nil(t, err)
  assert.Equal(t, user.Account, account)
  assert.Equal(t, user.Passwd, passwd)
  assert.Equal(t, user.Email, email)
  assert.Equal(t, user.Groups[0], group)
}

func TestUpdateUser(t *testing.T) {
  expected_id := uint(1)
  expected_account := "foo"
  expected_passwd := "bar"
  expected_email := "foo@bar.idv2"
  group0, _ := group_utils.Get("Guest")
  group1, _ := group_utils.Get("Administrator")

  user, err := user_utils.GetById(expected_id)
  assert.Nil(t, err)

  user.Email = expected_email
  user.Groups = []Group{group0, group1}

  user_utils.Update(&user)

  assert.Equal(t, user.ID, expected_id)
  assert.Equal(t, user.Account, expected_account)
  assert.Equal(t, user.Passwd, expected_passwd)
  assert.Equal(t, user.Email, expected_email)
  assert.Equal(t, user.Groups[0], group0)
  assert.Equal(t, user.Groups[1], group1)
}

func TestCountUser(t *testing.T) {
  var expected_cnt int64 = 1

  cnt, err := user_utils.Count()

  assert.Nil(t, err)
  assert.Equal(t, cnt, expected_cnt)
}

func TestIsInGroups(t *testing.T) {
  id := uint(1)
  in_group_names := []string{"Guest", "Administrator"}
  not_group_names := []string{"HaHa", "Point"}

  res1 := user_utils.IsInGroups(id, in_group_names)
  res2 := user_utils.IsInGroups(id, not_group_names)

  assert.Equal(t, true, res1)
  assert.Equal(t, false, res2)
}
