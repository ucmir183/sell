package main

import (
	"sell/models"
	_ "sell/routers"

	"github.com/astaxie/beego"
)

func main() {
	models.Init()
	beego.Run()
}
