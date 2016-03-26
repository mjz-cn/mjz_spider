package models

import (
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql" // import your used driver

    "mjz_spider/config"
)

type Proxy struct {
	Id int 				`orm:"auto"`
	Host string			`orm:"size(16)"`
	Port string
	Type int8
	Anonymity int8
	Country string 
}

func SaveProxy(h, p string, t, a int8, c string) {
    o := orm.NewOrm()
    o.Using("default")
    proxy := Proxy{
        Host: h,
        Port: p,
        Type: t,
        Anonymity: a,
        Country: c,
    }
    o.Insert(&proxy)
}

func init() {
    // set default database
    orm.RegisterDataBase("default", "mysql", config.GlobalConfig.MysqlConn, 30)
    
    // register model
    orm.RegisterModel(new(Proxy))
}