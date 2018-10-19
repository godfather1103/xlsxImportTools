package routers

import (
	"github.com/astaxie/beego"
	"github.com/godfather1103/xlsImportTools/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
