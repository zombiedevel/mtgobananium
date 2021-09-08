package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DB    *gorm.DB
	err   error
	DBErr error
)

type Users struct {
	*gorm.Model
	ChatID int64
	UserID int32
	IsProtected bool
	State string
}

type Database struct {
	*gorm.DB
}

//initialize database engine
func Setup() {
	var db = DB
	db, err = gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		DBErr = err
		fmt.Println("db err: ", err)
	}

	db.LogMode(true)
    db.AutoMigrate(&Users{})
	DB = db

}

func GetDB() *gorm.DB {
	return DB
}

func GetDBErr() error {
	return DBErr
}

