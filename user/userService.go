package user

import (
	"errors"
	"fmt"
)

//密码or用户名验证
func CheckPsw(mod *UserBind) (uint,string,error) {
	res,err:=FindUserByUserName(mod.UserName)
	if err != nil {
		return 0,"",fmt.Errorf("db error:%s",err.Error())
	}
	if mod.Password !=res.Password {
		return 0,"",errors.New("用户名或密码错误")
	}
	return  res.ID,res.Level, nil
	
}

//VIP验证
func CheckVIP(mod *UserBind) error  {
	res,err:=FindUserByUserName(mod.UserName)
	if err != nil {
		return fmt.Errorf("db error:%s",err.Error())
	}
	if res.Level != VIP {
		return errors.New("您不是VIP")
	}
	return nil
}

