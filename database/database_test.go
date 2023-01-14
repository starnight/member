package database

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestGetDBStr(t *testing.T) {
  str1 := GetDBStr("test")

  assert.Equal(t, str1, "test.db")

  str2 := GetDBStr("product")

  assert.Equal(t, str2, "product.db")
}

func TestConnectDB(t *testing.T) {
  db1 := ConnectDB(GetDBStr("test"))

  assert.NotNil(t, db1)

  db2 := ConnectDB(GetDBStr("test"))

  assert.Equal(t, db2, db1)
}

func TestCreateTables(t *testing.T) {
  db := ConnectDB(GetDBStr("test"))

  assert.NotNil(t, db)

  CreateTables(db)
}
