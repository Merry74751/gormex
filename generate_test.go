package gorm_expand

import "testing"

func TestGenerateSelect(t *testing.T) {
	s := generateSelect[User]()
	t.Log(s)
}

func TestGenerateCreate(t *testing.T) {
	u := User{Username: "zhangsan", Age: 18}
	s, _, _ := generateCreate(u)
	t.Log(s)

}

type User struct {
	Username string
	Age      int
	Sex      string
}
