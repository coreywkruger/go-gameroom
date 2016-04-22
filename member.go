package main

import (
  "errors"
  "log"
)

// list of all members
var memberList = make(map[int]*Member)

// Member - member of connection group
type Member struct {
  ID         int
  Connection *Connection
}

// RegisterMember - creates new member
func RegisterMember(id int, c *Connection) (*Member, error) {
  log.Println("registering ", id)
  memberList[id] = &Member{
    ID:         id,
    Connection: c,
  }
  log.Println("# of members: ", len(memberList))
  return memberList[id], nil
}

// GetMember - gets member by id
func GetMember(id int) (*Member, error) {
  m := memberList[id]
  if m != nil {
    log.Println("Getting Member: ", m.ID)
    return m, nil
  }
  return nil, errors.New("Member not found")
}

// GetAllMembers - gets all members
func GetAllMembers() map[int]*Member {
  return memberList
}
