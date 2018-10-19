package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type FileVo struct {
	Id              int64
	FileName        string `orm:size(400)`
	FileFullPath    string `orm:size(4000)`
	FileMd5         string `orm:size(40)`
	LastVersionTime time.Time
}

func init() {
	orm.RegisterModel(new(FileVo))
}

func FindFileByMd5(md5 string) (*FileVo, error) {
	o := orm.NewOrm()
	fileVo := make([]*FileVo, 0)
	_, err := o.QueryTable("file_vo").Filter("file_md5", md5).Limit(1).All(&fileVo)
	if len(fileVo) > 0 {
		return fileVo[0], err
	} else {
		return nil, err
	}
}

func FindFileByFullPath(fileFullPath string) (*FileVo, error) {
	o := orm.NewOrm()
	fileVo := make([]*FileVo, 0)
	_, err := o.QueryTable("file_vo").Filter("file_full_path", fileFullPath).Limit(1).All(&fileVo)
	if len(fileVo) > 0 {
		return fileVo[0], err
	} else {
		return nil, err
	}
}
