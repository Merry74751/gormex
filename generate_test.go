package gorm_expand

import "testing"

func TestGenerateCreate(t *testing.T) {
	u := User{Username: "zhangsan", Age: 18}
	s, _ := generateInsert(u)
	t.Log(s)

}

type User struct {
	Username string
	Age      int
	Sex      string
}
