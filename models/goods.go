package models

import (
	"github.com/astaxie/beego/orm"
	"fmt"
)

type Goods struct {
	Id int
	Name string
	CateId int
	Desc string
	Price float64 `orm:"digits(12);decimals(2)"`
	Inventory int
	UpdateId int
	CreateId int
	CreateTime int64
	UpdateTime int64
	Status int8

}


func init() {
	orm.RegisterModel(new(Goods))
}

func (a *Goods) TableName() string {
	return TableName("goods")
}

func GoodsGetList(page, pageSize int, filters ...interface{})([]*Goods, int64) {
	offset := (page -1) * pageSize
	list := make([]*Goods,0)

	query := orm.NewOrm().QueryTable(TableName("goods"))
	if len(filters) > 0 {
		l := len(filters)
		for k:=0; k < l ; k+=2 {
			fmt.Println(filters[k].(string),filters[k+1])
			query = query.Filter(filters[k].(string),filters[k+1])

		}
	}


	total,_ := query.Count()
	query.OrderBy("-id").Limit(pageSize,offset).All(&list)

	return list,total
}


func GoodsGetById(id int) *Cate {
	cate := new(Cate)
	orm.NewOrm().QueryTable(TableName("cate")).Filter("id",id).One(cate)

	return cate

}

func GoodsAdd(c * Goods) (int64,error) {
	return orm.NewOrm().Insert(c)

}


func ( cate *Goods ) Update(fields ...string)(int64,error) {
	return orm.NewOrm().Update(cate)
}