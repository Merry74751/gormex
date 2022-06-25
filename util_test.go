package gorm_expand

import "testing"

func TestEndWith(t *testing.T) {
	s := "Hello world"
	s2 := "world"
	t.Log(endWith(s, s2))

	s3 := "wrol"
	t.Log(endWith(s, s3))
}

func TestToUnderLineCase(t *testing.T) {
	s := "HelloWorld"
	t.Log(toUnderLineCase(s))
}
