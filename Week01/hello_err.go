package main

import (
"database/sql"
"log"

"github.com/pkg/errors"
)

var (
DataNotExist = errors.New("model: data not exists")
)

type Student struct {
ID   uint
Name string
}

func (s *Student) MockDbNoDataError() error {
return sql.ErrNoRows
}

func (s *Student) MockFindError() error {
if err := s.MockDbNoDataError(); err != nil {
if errors.Is(err, sql.ErrNoRows) {
err = DataNotExist
}

return errors.Wrap(err, "student model")
//return err
}

return nil
}

func BizUserDetail(uid uint) (*Student, error) {
s := &Student{ID: uid}
if err := s.MockFindError(); err != nil {
return nil, errors.WithMessagef(err, "biz query user: %d detail", uid)
}

return s, nil
}

func main() {
stu, err := BizUserDetail(1)
if err != nil {
if errors.Is(err, DataNotExist) {
log.Printf("use not exists %+v\n", err)
return
}

log.Printf("query user detail failed: %+v\n", err)
return
}

log.Printf("user info: %+v\n", stu)
}
