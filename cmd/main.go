package main

import (
	"BaiDuPanLoad/file"
	"BaiDuPanLoad/user"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	r:=gin.Default()

	r.POST("/signup",user.Signup) //注册
	r.POST("/login",user.Login) //登录

	v1:=r.Group("/v1")
	v1.Use(user.GetJWT)
	{
		v1.POST("/upload",file.UploadFile)  //上传文件
		v1.POST("/load/*file",file.Loadfile)  //下载文件
		v1.POST("/normallink",file.CreateLink)  //生成分享链接
		v1.POST("/normallink_qcode",file.CreateNewLinkCode) //生成分享二维码
		v1.POST("/secretlink",file.CreateSecretLink)  //生成加密链接

	}

	r.Run("59.110.23.117:8080")
	//r.Run()
}