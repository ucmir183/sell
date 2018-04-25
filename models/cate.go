package models

import (
	"github.com/astaxie/beego/orm"
	"fmt"
)

type Cate struct {
	Id int
	Name string
	Pid int
	CreateId int
	UpdateId int
	Status   int
	CreateTime int64
	UpdateTime int64
}

func init()  {
	orm.RegisterModel(new(Cate))

}

func (a *Cate) TableName() string {
	return TableName("cate")
}


func CateGetList(page, pageSize int, filters ...interface{})([]*Cate, int64) {
	offset := (page -1) * pageSize
	list := make([]*Cate,0)

	query := orm.NewOrm().QueryTable(TableName("cate"))
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

func CateGetAllList() []*Cate {
	cate := make([]*Cate,0)
	orm.NewOrm().QueryTable(TableName("cate")).Filter("status",1).All(&cate)
	return cate

}

func CateGetById(id int) *Cate {
	cate := new(Cate)
	orm.NewOrm().QueryTable(TableName("cate")).Filter("id",id).One(cate)

	return cate

}

func CateAdd(c * Cate) (int64,error) {
	return orm.NewOrm().Insert(c)

}

func ( cate *Cate ) Update(fields ...string)(int64,error) {
	return orm.NewOrm().Update(cate)
}