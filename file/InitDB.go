package file

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB()  {
	db, err := gorm.Open("mysql", "root:root@/panLoad?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(fmt.Errorf("db error :%s",err.Error()))
	}
	db.AutoMigrate(&File{},&Filemeta{},&FiileSecretLink{})
	DB=db

}

func init() {
	InitDB()
}

