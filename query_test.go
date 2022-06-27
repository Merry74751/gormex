package gorm_expand

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"testing"
	"time"
)

type Users struct {
	Id         uint      `json:"id,omitempty"`
	Username   string    `json:"username,omitempty"`
	Password   string    `json:"password,omitempty"`
	Status     *uint     `json:"status,omitempty"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

func (u Users) TableName() string {
	return "user"
}

func (u Users) String() string {
	bytes, err := json.Marshal(&u)
	if err != nil {

	}
	return string(bytes)
}

func connection() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	s, err := db.DB()
	if err != nil {
		log.Println(err)
	}
	s.SetMaxIdleConns(10)
	s.SetMaxOpenConns(100)
	s.SetConnMaxLifetime(time.Hour)
	if err != nil {
		log.Println(err)
	}
	return db
}

func TestMapper_Insert(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	i := uint(0)

	u := Users{
		Username: "老八",
		Password: "123456",
		Status:   &i,
	}
	err := m.Insert(u)
	if err != nil {
		t.Log(err)
	}
}

func TestMapper_DeleteById(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	err := m.DeleteById(3)
	if err != nil {
		t.Log(err)
	}
}

func TestMapper_UpdateById(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	i := uint(1)

	u := Users{Id: 2, Username: "王五2", Status: &i}
	err := m.UpdateById(u)
	if err != nil {
		t.Log(err)
	}
}

func TestMapper_GetById(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	user, err := m.GetById(2)
	if err != nil {
		t.Log(err)
	}
	t.Log(user)
}

func TestMapper_Get(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	u := Users{Id: 1, Username: "王五2"}
	users, err := m.Get(u)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(users)
	}
}

func TestMapper_List(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	list, err := m.List()
	if err != nil {
		t.Log(err)
	} else {
		for i := range list {
			t.Log(list[i])
		}
	}
}

func TestMapper_ListByCondition(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	i := uint(0)
	u := Users{Status: &i}
	list, err := m.ListByCondition(u)
	if err != nil {
		t.Log(err)
	} else {
		for i := range list {
			t.Log(list[i])
		}
	}
}

func TestMapper_Page(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	page, i, err := m.Page(Page{Current: 1, PageSize: 2})
	if err != nil {
		t.Log(err)
	} else {
		t.Log(i)
		for i2 := range page {
			t.Log(page[i2])
		}
	}
}

func TestMapper_PageByCondition(t *testing.T) {
	m := Mapper[Users]{}
	m.SetDb(connection())

	i := uint(0)
	u := Users{Status: &i}

	page, i2, err := m.PageByCondition(u, Page{1, 2})
	if err != nil {
		t.Log(err)
	} else {
		t.Log(i2)
		for i3 := range page {
			t.Log(page[i3])
		}
	}
}
