package file

import (
	"BaiDuPanLoad/user"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	NORMALSIZE = 1024*1024
	VIPSIZE =1024*1024*3
)

var (
	UPuserID uint
	filesize int64
	fileloadsize int
)

type FileForm struct {
	userID uint `json:"userid"`
	FileMD5  string `json:"filemd5"`
	ISUpload  string  `json:"isupload"`  //上传过--1 ，没有上传过--0
}

type FileLinkItem struct {
	FileName string `json:"file_name" binding:"required"`
	Password string  `json:"password" binding:"required,len=4"`
}

type Filemeta struct {
	Filename string
	Filesize  int
	FileReadSize int
	Updatetime time.Time
	ISDown   string   // 下载完成 -- 1 未完成 --0
}



func UploadFile(c *gin.Context)  {
	UPuserIDstr,ok:=c.Get("userid")
	level:=c.GetString("level")
	if !ok{
		c.String(http.StatusInternalServerError,errors.New("上下文错误").Error())
	}

	res,_:=json.Marshal(UPuserIDstr)
	json.Unmarshal(res,UPuserID)

	filemd5:=c.PostForm("filemd5")

	f :=&FileForm{
		userID:  UPuserID,
		FileMD5: filemd5,
		ISUpload: "1",
	}


	fileheader,err:=c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"上传文件失败",
			"data":err,
		})
		return
	}


	//不能上传空文件
	if fileheader.Size <=0 {
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"不能上传空文件",
		})
	}

	//对vip和普通用户限制上传文件大小
	if level == user.NORMAL {
		filesize = NORMALSIZE
	}else if level == user.VIP {
		filesize = VIPSIZE
	}
	if fileheader.Size > filesize {
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"文件太大了",
		})
		return
	}



	//康康是否上传过
	result,err:=SecUPLoad(filemd5)
	if err != nil  {
		if err:=SaveFile(fileheader,f);err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"上传文件失败",
				"data":err,
			})
			return
		}
		if err:=c.SaveUploadedFile(fileheader,"./files/"+fileheader.Filename);err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"上传文件失败",
				"data":err,
			})
			return
		}
	}else {
		//秒传
		result.UPUserID = UPuserID
		if err:=SecCreateNewFile(result);err!=nil{

			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"秒传失败",
				"data":err,
			})

		}
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"上传成功",
	})

}



func Loadfile(c *gin.Context) {

	var (
		p     = make([]byte, fileloadsize)
		idint int
	)

	level := c.GetString("level")
	UserId, _ := c.Get("userid")

	id, _ := json.Marshal(UserId)
	json.Unmarshal(id, idint)

	fileurl := c.Request.URL.Path
	fmt.Println(fileurl)
	res, err := FindFileByURL(fileurl)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "文件不存在",
		})
		return
	}

	//查询文件是否下载过
	result, err := FindFileLogByFileName(res.FileName)
	if err != nil {

	//没有下载过
	file, err := os.Open(res.FileName)
	defer file.Close()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "文件不存在",
			"data":    err,
		})
		return
	}

	//限制上传文件大小
	if level == user.NORMAL {
		fileloadsize = 512
	}
	c.Request.Header.Add("Range", fmt.Sprintf("byte=%d", fileloadsize))

	r, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "文件不存在",
			"data":    err,
		})
		return
	}

	//文件长度大于限制大小
	//新建临时文件
	if len(r) > fileloadsize {   //todo
		file.Read(p)
		tfile, err := os.Create(fmt.Sprintf("./temporary_file/%d%s", idint, res.FileName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "sever error",
				"data":    err,
			})
			return
		}
		written, _ := io.Copy(tfile, file)

		f := &Filemeta{
			Filename:     string(idint)+res.FileName,
			Filesize:     len(r),
			FileReadSize: len(r) - int(written),
			Updatetime:   time.Now(),
			ISDown:       "0",
		}
		if err := CreateFileLog(f); err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"sever error",
				"data":err,
			})
			return
		}
		c.JSON(http.StatusPartialContent,gin.H{
			"message":"下载失败",
		})
		return
	}else {
		
		//下载过
		var tfilepath = fmt.Sprintf("./temporary_file/%d%s", idint, res.FileName)

		tfile,_:=os.Open(tfilepath)
		defer tfile.Close()
		lastsize := result.Filesize-result.FileReadSize //剩下的字节
		readsize :=int64(result.FileReadSize)
		if lastsize < fileloadsize {
			tfile.Seek(readsize,0)
			_,err:=tfile.Read(p)   //todo
			if err == io.EOF {

				//更新filelog
				f :=&Filemeta{
					Filename:     string(idint)+res.FileName,
					FileReadSize: len(r),
					Updatetime:   time.Now(),
					ISDown:       "1",
				}
				if err:=UpdeteFileLog(f);err!=nil{

					c.JSON(http.StatusInternalServerError,gin.H{
						"message":"server error",
						"data":err,
					})
					return

				}



			}
		}
	}
}

	c.Request.Header.Set("Content-type","application/octet-stream")
	c.Request.Header.Add("Content-Disposition",fmt.Sprintf("attachment; filename=%s/", fmt.Sprintf("./temporary_file/%d%s", idint, res.FileName)))
	c.File(fmt.Sprintf("./files/%s",fmt.Sprintf("./temporary_file/%d%s", idint, res.FileName)))
}


//生成分享链接
func CreateLink(c *gin.Context)  {

	filaname:=c.PostForm("filename")

	res,err:=FindFileByFilename(filaname)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"文件不存在",
			"data":err.Error(),
		})
		return
	}


	c.JSON(http.StatusOK,gin.H{
		"message":c.Request.Host+res.FileURL,
	})

}


//生成分享链接二维码
func CreateNewLinkCode(c *gin.Context)  {
	filename:=c.PostForm("filename")

	res,err:=FindFileByFilename(filename)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"文件不存在",
			"data":err.Error(),

		})
		return
	}

	codefile :=fmt.Sprintf("%s.png",filename)
	//result,_:=os.Create("./code/"+codefile)
	//defer result.Close()
	err2:=qrcode.WriteFile(c.Request.Host+res.FileURL,qrcode.Medium,256,fmt.Sprintf("./code/%s",codefile))
	if err2 != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"生成二维码失败",
			"data":err2.Error(),
		})
		return
	}

	c.File(fmt.Sprintf("./code/%s",codefile))
}


//生成加密分享链接
func CreateSecretLink(c *gin.Context)  {

	var  LINKHEAD = c.Request.Host +"/share/"

	f :=&FileLinkItem{}
	if err:=c.BindJSON(f);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"生成链接错误",
			"data":err.Error(),
		})
		return
	}

	_,err:=FindFileByFilename(f.FileName)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"文件不存在",
			"data":err.Error(),
		})
		return
	}

	linkstr:=FileMD5(f.FileName)
	if err:=CreateNewShareLink(linkstr,f.Password);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"生成分享链接错误",
			"data":err.Error(),
		})
		return
	}


	c.JSON(http.StatusOK,gin.H{
		"sharelink":LINKHEAD+linkstr,
		"password":f.Password,
	})
}


