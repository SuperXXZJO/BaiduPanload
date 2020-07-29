package file

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

//md5加密
func FileMD5(fiilename string)(linkstr string) {
	m:=md5.New()
	res:=m.Sum([]byte(fiilename+time.Now().String()))
	return hex.EncodeToString(res)
}