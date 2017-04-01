package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/ralphg6/goshua/models"
)

type Account struct {
	Id     int    `orm:"column(id);auto"`
	Title  string `orm:"column(title);size(250)"`
	Name   string `orm:"column(name);size(250)"`
	Type   string `orm:"column(type);size(20)"`
	Domain string `orm:"column(domain);size(100)"`
}

func (t *Account) TableName() string {
	return "account"
}

func init() {
	orm.RegisterModel(new(Account))
}

type AccountModel struct {
	models.BaseCRUDModel
}

func NewAccountModel() (model AccountModel) {
	model.NewInstance = func(id int) interface{} {
		return &Account{Id: id}
	}
	return
}
