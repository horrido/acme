package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type AuthUser struct {
	Id       int
	First    string
	Last     string
	Email    string `orm:"unique"`
	Password string
	Reg_key  string
	Reg_date time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(AuthUser))
}
