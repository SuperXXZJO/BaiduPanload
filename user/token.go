package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

const jwtKey = "PanLoad"

type Claims struct {
	jwt.StandardClaims

	UserID  uint
	UserLevel    string
}


func SetJWT(userid uint,level string) (string,error) {
	expireTime := time.Now().Add(10 * time.Hour)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},

		UserID: userid ,
		UserLevel: level,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return tokenString,nil
}


func GetJWT(c *gin.Context)  {
	tokenstring:=c.GetHeader("Authorization")

	claims := &Claims{}
	token,err := jwt.ParseWithClaims(tokenstring,claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(jwtKey),nil
	})
	if err == nil && token.Valid{
		c.Set("userid",claims.UserID)
		c.Set("level",claims.UserLevel)
		c.Next()
	}else {
		c.JSON(300,"验证失败！请先登录！")
		return
	}

}

