package file

import "mime/multipart"

const FILEURL = "/load/"


type fileItem struct {
	UPUserID   uint
	FileName   string
	FileURL    string
	FileMD5    string
	ISUPload   string
}


func SaveFile(file *multipart.FileHeader,form *FileForm) error  {
	f :=&fileItem{
		UPUserID: form.userID,
		FileName: file.Filename,
		FileURL: FILEURL+file.Filename,
		FileMD5: form.FileMD5,
	}
	if err:=CreateNewFile(f);err!=nil{
		return err
	}
	return nil
}

func SecUPLoad(filemd5 string) (*File,error)  {
	res,err:=FindFileByMD5(filemd5)
	if err != nil {
		return nil,err
	}
	return res,nil
}

