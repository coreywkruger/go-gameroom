package main

import (
  "errors"
  "log"
)

var memberList = make(map[int]*Member)

// Member -
type Member struct {
  ID int
  WS *Connection
}

// RegisterMember -
func RegisterMember(id int, c *Connection) (*Member, error) {
  log.Println("registering ", id)
  memberList[id] = &Member{
    ID: id,
    WS: c,
  }
  log.Println(len(memberList))
  return memberList[id], nil
}

// GetMember -
func GetMember(id int) (*Member, error) {
  m := memberList[id]
  if m != nil {
    log.Println("Found ", m.ID)
    return m, nil
  }
  return nil, errors.New("Member not found")
}

// GetAllMembers -
func GetAllMembers() map[int]*Member {
  return memberList
}
