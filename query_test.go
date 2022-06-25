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
	Id         int       `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	CreateTime time.Time `gorm:"column:createTime" json:"createTime"`
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

func TestInsert(t *testing.T) {
	db := connection()
	m := Mapper[Users]{db}
	u := Users{3, "insert", time.Now()}
	err := m.Insert(u)
	if err != nil {
		t.Log(err)
	}
}

func TestSelectById(t *testing.T) {
	db := connection()
	m := Mapper[Users]{db}
	u, err := m.GetById(1)
	if err != nil {
		t.Log(err)
	}
	t.Log(u)
}

func TestUpdate(t *testing.T) {
	db := connection()
	m := Mapper[Users]{db}
	u := Users{Id: 1, Name: "李四"}
	err := m.UpdateById(u)
	if err != nil {
		t.Log(err)
	}
}

func TestList(t *testing.T) {
	db := connection()
	m := Mapper[Users]{db}
	list, err := m.List()
	if err != nil {
		t.Log(err)
	}
	t.Log(list)
}

func TestExec(t *testing.T) {
	db := connection()
	u := Users{}
	db.Raw("select * from user limit 1").Scan(&u)
	t.Log(u)
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

func TestPage(t *testing.T) {
	c := connection()
	m := Mapper[Users]{c}

	page, total, _ := m.Page(Page{Current: 1, PageSize: 1})
	t.Log(page)
	t.Log(total)
}
