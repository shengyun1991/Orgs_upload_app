package routers

import (
	"Orgs_upload_app/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/ogaj",&controllers.AuthControllers{},"post:RecodeAuth")
	beego.Router("/ofgj",&controllers.CertControllers{},"post:RecodeHouse")
	beego.Router("/ozxzx",&controllers.CreditControllers{},"post:RecodeCredit")
}
