package user

import (
	"github.com/jinzhu/gorm"
)

const (
	NORMAL= "normal"
	VIP  = "vip"
)


type User struct {
	gorm.Model

	UserName string `gorm:"column:username"`
	Password string
	Level  string
}


//创建新的user
func CreateNewUser(username string,password string) error {
	mod:=&User{
		UserName: username,
		Password: password,
		Level: NORMAL,
	}
	if err:=DB.Create(mod).Error;err!=nil{
		return err
	}
	return nil
}

//通过username查询
func FindUserByUserName(username string)( *User,error) {

	mod:=&User{}

	if err:=DB.Table("users").Where("username = ?",username).First(mod).Error;err!=nil{
		return nil,err
	}
	return mod,nil
}

