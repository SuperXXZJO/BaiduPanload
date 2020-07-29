package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserBind struct {
	UserName  string `json:"user_name" binding:"required,min=5,max=8"`
	Password  string  `json:"password" binding:"required,min=6,max=16"`
}



//登录
func Login(c *gin.Context) {

	u:=&UserBind{}
	if err:=c.BindJSON(u);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}

	id,level,err:=CheckPsw(u)
	if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"err":err.Error(),
		})
		return
	}

	token,err :=SetJWT(id,level)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"登录成功",
		"token":token,
	})

}

//注册
func Signup(c *gin.Context)  {
	u:=&UserBind{}
	if err:=c.BindJSON(u);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}
	if err:=CreateNewUser(u.UserName,u.Password);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"注册失败",
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"注册成功",
	})
}