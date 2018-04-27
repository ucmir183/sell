package controllers

import (
	"strings"
	"sell/models"
	"github.com/astaxie/beego"
	"time"
)

type GoodsController struct {
	BaseController
}

func ( self *GoodsController ) List() {
	self.Data["pageTitle"] = "商品列表"
	self.display()
}

func(self *GoodsController) Table() {
	page,err := self.GetInt("page")
	if err != nil {
		page = 1
	}

	limit,err := self.GetInt("limit")

	if err != nil {
		limit = 30
	}

	name := strings.TrimSpace(self.GetString("name"))

	self.pageSize = limit

	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)

	if name != "" {
		filters = append(filters,"name__contains",name)
	}

	result,count := models.GoodsGetList(page,self.pageSize,filters...)


	list := make([]map[string]interface{},len(result))

	for k,v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		row["create_time"] = beego.Date(time.Unix(v.CreateTime,0),"Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime,0),"Y-m-d H:i:s")

		list[k] = row

	}

	self.ajaxList("成功",MSG_OK,count,list)

}

func (self *GoodsController)Add() {
	cate_list := models.CateGetAllList()

	self.Data["pageTitle"] = "添加商品"
	self.Data["cate_list"] = cate_list
	self.display()

}


func (self *GoodsController) Edit() {
	id,_ := self.GetInt("id",0)
	cate_list := models.CateGetAllList()

	info := models.GoodsGetById(id)

	self.Data["pageTitle"] = "修改商品"

	self.Data["cate_list"] = cate_list
	self.Data["info"] = info
	self.display()

}

func (self *GoodsController) AjaxSave() {
	id,_ := self.GetInt("id",0)
	if id == 0 {
		m := new(models.Goods)
		m.Name = strings.TrimSpace(self.GetString("name"))
		m.Name = strings.TrimSpace(self.GetString("name"))
		m.Img = strings.TrimSpace(self.GetString("img"))
		m.Price,_ = self.GetFloat("price")
		m.CateId,_ = self.GetInt("cate_id",0)
		m.Inventory,_ = self.GetInt("inventory",0)
		m.Status = 1

		m.CreateId = self.userId
		m.UpdateId = self.userId
		m.CreateTime = time.Now().Unix()
		m.UpdateTime = time.Now().Unix()

		if _,err := models.GoodsAdd(m); err != nil {
			self.ajaxMsg(err.Error(),MSG_ERR)
		}

		self.ajaxMsg("保存成功",MSG_OK)

	}

	info := models.GoodsGetById(id)
	info.Name = strings.TrimSpace(self.GetString("name"))
	info.Desc = strings.TrimSpace(self.GetString("desc"))
	info.Img = strings.TrimSpace(self.GetString("img"))
	info.Price,_ = self.GetFloat("price")
	info.CateId,_ = self.GetInt("cate_id",0)
	info.Inventory,_ = self.GetInt("inventory",0)

	info.UpdateId = self.userId
	info.UpdateTime = time.Now().Unix()

	if _,err := info.Update();err != nil {
		self.ajaxMsg(err.Error(),MSG_ERR)
	}

	self.ajaxMsg("保存成功",MSG_OK)

}

func (self *GoodsController) AjaxDel() {

	id, _ := self.GetInt("id")
	info := models.GoodsGetById(id)
	if info == nil {
		self.ajaxMsg("未找到相应的信息",MSG_ERR)
	}

	info.UpdateTime = time.Now().Unix()
	info.UpdateId = self.userId
	info.Status = 0
	info.Id = id


	if _,err := info.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("删除成功", MSG_OK)
}