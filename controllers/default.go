package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "chegongcai.go"
	c.Data["Email"] = "chegc@szbdy.org"
	c.Data["PhoneNumber"] = "18689465458"
	c.TplName = "index.tpl"
	//c.Ctx.WriteString("hello")
}
