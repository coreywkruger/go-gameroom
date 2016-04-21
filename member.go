package main

import (
  "errors"
  "log"
)

var memberList = make(map[string]*Member)

// Member -
type Member struct {
  ID string
  WS *Connection
}

// RegisterMember -
func RegisterMember(id string, c *Connection) (*Member, error) {
  log.Println("registering ", id)
  memberList[id] = &Member{
    ID: id,
    WS: c,
  }
  return memberList[id], nil
}

// GetMember -
func GetMember(id string) (*Member, error) {
  m := memberList[id]
  if m != nil {
    log.Println("Found ", m.ID)
    return m, nil
  }
  return nil, errors.New("Member not found")
}
