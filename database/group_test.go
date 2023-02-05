package database

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestGetGroups(t *testing.T) {
  utils := GroupUtils{DB: ConnectDB(GetDBStr("test"))}
  expect_groups := []string{"Guest", "Administrator"}

  for _, name := range expect_groups {
    group, err := utils.Get(name)
    assert.Nil(t, err)
    assert.Equal(t, name, group.Name)
  }
}
