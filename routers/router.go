package routers

import (
	"github.com/astaxie/beego"
	"github.com/godfather1103/xlsxImportTools/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
