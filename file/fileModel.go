package file

import (
	"github.com/jinzhu/gorm"
	"time"
)

type File struct {
	gorm.Model
	UPUserID   uint
	FileName   string
	FileURL    string
	FileMD5    string
	ISUPLoad    string
}

type FiileSecretLink struct {
	FileURL  string
	Password string
}


type FileLog struct {
	gorm.Model
	FileName string
	FileSize  int
	FileReadSize int
	Updatetime  time.Time
	ISDown   string    //是否下载完成
}


//
func CreateNewFile(file *fileItem) error {

	f:=&File{
		UPUserID: file.UPUserID,
		FileName: file.FileName,
		FileURL: file.FileURL,
		FileMD5: file.FileMD5,
		ISUPLoad: "1",
	}
	if err:=DB.Table("files").Create(f).Error;err!=nil{
		return err
	}
	return nil
}

//
func SecCreateNewFile(f *File) error {
	if err:=DB.Table("files").Create(f).Error;err!=nil{
		return err
	}
	return nil
}

//
func FindFileByURL(Url string) (res *File,err error) {
	if err :=DB.Table("files").Where("file_url = ?",Url).First(res).Error;err!=nil{
		return nil,err
	}
	return res, nil
}

//
func FindFileByMD5(filemd5 string) (res *File,err error) {
	if err:=DB.Table("files").Where("file_md5 = ?",filemd5).First(res).Error;err!=nil{
		return nil,err
	}
	return res,nil
}

//
func FindFileByFilename(filename string)(*File,error)  {
	res:=&File{}
	if err:=DB.Table("files").Where("file_name = ?",filename).First(res).Error;err!=nil{
		return nil,err
	}
	return res,nil
}

//生成加密连接
func CreateNewShareLink(link string,password string) error {
	f :=&FiileSecretLink{
		FileURL:  link,
		Password: password,
	}
	if err:=DB.Table("fiile_secret_links").Create(f).Error;err!=nil{
		return err
	}
	return nil
}


//创建文件日志数据
func CreateFileLog(f *Filemeta) error {

	flog :=&FileLog{
		FileName:     f.Filename,
		FileSize:     f.Filesize,
		FileReadSize: f.FileReadSize,
		Updatetime:   f.Updatetime,
		ISDown:       f.ISDown,
	}
	if err:=DB.Table("file_log").Create(flog).Error;err!=nil{
		return err
	}
	return nil
}

//查询文件日志数据
func FindFileLogByFileName(filename string) (filelog *Filemeta ,err error ) {
	if err:=DB.Table("file_log").Where("file_name = ?",filename).First(filelog).Error;err!=nil{
		return nil,err
	}
	return filelog,nil
}

//更新数据

func UpdeteFileLog(f *Filemeta) error  {
	if err:=DB.Table("file_log").Where("file_name = ?",f.Filename).Save(f).Error;err!=nil{
		return err
	}
	return nil
}