package util

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/godfather1103/xlsxImportTools/models"
	"github.com/tealeg/xlsx"
	"log"
	"strings"
	"time"
)

type Config struct {
	XlsToDbSkipOneLine int
	XlsToDbKey         string
	XlsToDbField       string
	XlsToDbTabls       string
}

func initParam() (*Config, error) {
	config := &Config{
		XlsToDbSkipOneLine: beego.AppConfig.DefaultInt("xls_to_db_skipOneLine", 0),
		XlsToDbKey:         beego.AppConfig.DefaultString("xls_to_db_key", ""),
		XlsToDbField:       beego.AppConfig.String("xls_to_db_field"),
		XlsToDbTabls:       beego.AppConfig.String("xls_to_db_tabls"),
	}
	if len(config.XlsToDbTabls) < 1 || len(config.XlsToDbField) < 1 {
		return nil, errors.New("参数xls_to_db_field或xls_to_db_tabls异常")
	}
	return config, nil
}

func ImportData(fileVo models.FileVo) error {
	config, err := initParam()
	if err != nil {
		return err
	}

	var fields = strings.Split(config.XlsToDbField, ",")
	var keyIndex = -1
	for index, item := range fields {
		if strings.EqualFold(item, config.XlsToDbKey) {
			keyIndex = index
		}
	}

	if len(config.XlsToDbKey) > 0 && keyIndex < 0 {
		return errors.New("主键必须要在xls_to_db_field中存在")
	}

	xlFile, err := xlsx.OpenFile(fileVo.FileFullPath)
	if err == nil {
		var item = make([]string, len(fields))
		for _, sheet := range xlFile.Sheets {
			for rowIndex, row := range sheet.Rows {
				if rowIndex == 0 && config.XlsToDbSkipOneLine == 1 {
					continue
				}
				for cellIndex, cell := range row.Cells {
					if cellIndex >= len(fields) {
						break
					}
					text := cell.String()
					item[cellIndex] = text
				}
				if keyIndex >= 0 && CheckDataExistsInTableByKey(config.XlsToDbTabls, config.XlsToDbKey, item[keyIndex]) {
					log.Printf("数据[%s=%s]已经存在数据中了", config.XlsToDbKey, item[keyIndex])
					continue
				} else {
					err := SaveData(item, config)
					if err != nil {
						return err
					}
				}
			}
		}
		o := orm.NewOrm()
		file, _ := models.FindFileByFullPath(fileVo.FileFullPath)
		if file != nil {
			file.FileName = fileVo.FileName
			file.FileFullPath = fileVo.FileFullPath
			file.FileMd5 = fileVo.FileMd5
			file.LastVersionTime = time.Now()
			id, err := o.Update(file)
			log.Printf("文件[id=%d,Name=%s]更新导入数据成功！", id, file.FileName)
			return err
		} else {
			id, err := o.Insert(&fileVo)
			log.Printf("文件[id=%d,Name=%s]导入数据成功！", id, fileVo.FileName)
			return err
		}
	} else {
		return err
	}
}

func CheckDataExistsInTableByKey(tableName string, key string, value string) bool {
	o := orm.NewOrm()
	var maps []orm.Params
	var sql = fmt.Sprintf("select %s from %s where %s=?", key, tableName, key)
	num, _ := o.Raw(sql, value).Values(&maps)
	return num > 0
}

func SaveData(data []string, config *Config) error {
	o := orm.NewOrm()
	var it = ""
	for index := range data {
		if index == 0 {
			it = "?"
		} else {
			it += ",?"
		}
	}
	if len(it) < 1 {
		return errors.New("数据字段不能为空！")
	}
	var sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", config.XlsToDbTabls, config.XlsToDbField, it)
	res, err := o.Raw(sql, data).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Printf("成功添加%d条数据", num)
		return nil
	} else {
		return err
	}
}
