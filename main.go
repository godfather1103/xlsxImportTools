package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/godfather1103/utils"
	"github.com/godfather1103/xlsxImportTools/models"
	_ "github.com/godfather1103/xlsxImportTools/routers"
	"github.com/godfather1103/xlsxImportTools/util"
	"log"
	"path/filepath"
	"time"
)

func init() {
	models.RegisterDB()
	orm.RunSyncdb("default", false, false)
}
func main() {
	var XlsToDbScantime = beego.AppConfig.DefaultInt("xls_to_db_scantime", 2)
	ticker := time.NewTicker(time.Duration(XlsToDbScantime) * time.Second)
	for t := range ticker.C {
		log.Println(t)
		AutoImportData()
	}
	//beego.Run()
}

func AutoImportData() {
	var XlsToDbScandir = beego.AppConfig.DefaultString("xls_to_db_scandir", ".")
	flag, _ := utils.PathExists(XlsToDbScandir)
	if XlsToDbScandir == "." || flag {
		files, err := utils.GetAllFile(XlsToDbScandir, ".xlsx", false)
		if err == nil {
			rootPath, _ := filepath.Abs(XlsToDbScandir)
			for _, f := range files {
				fileFullPath := rootPath + string(filepath.Separator) + f.Name()
				md5 := utils.GetFileMd5(fileFullPath)
				if len(md5) <= 0 {
					continue
				}
				fileVo, err := models.FindFileByMd5(md5)
				if fileVo == nil && err == nil {
					fileVo = &models.FileVo{
						FileName:        f.Name(),
						FileFullPath:    fileFullPath,
						FileMd5:         md5,
						LastVersionTime: time.Now(),
					}
					err := util.ImportData(*fileVo)
					if err != nil {
						log.Fatalf("导入数据出错，err=%s", err)
					}
				} else {
					continue
				}
			}
		} else {
			log.Fatalf("获取文件失败，err=%s", err)
		}
	} else {
		log.Fatal("扫描路径不存在！")
	}
}
