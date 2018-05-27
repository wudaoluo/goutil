package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"messagePush/glog"
	"bytes"
	"messagePush/conf"
)

var db *sql.DB


//username:password@tcp(dbhost:dbport)/dbname?charset=utf8
func Init() {
	var buf bytes.Buffer
	var err error

	cfg := conf.GetInstance()
	buf.WriteString(cfg.Cfg.Mysql.DBuser)
	buf.WriteString(":")
	buf.WriteString(cfg.Cfg.Mysql.DBpasswd)
	buf.WriteString("@tcp(")
	buf.WriteString(cfg.Cfg.Mysql.DBaddr)
	buf.WriteString(":")
	buf.WriteString(cfg.Cfg.Mysql.DBport)
	buf.WriteString(")/")
	buf.WriteString(cfg.Cfg.Mysql.DBname)
	buf.WriteString("?charset=utf8")

	glog.Info(buf.String())
	db, err = sql.Open("mysql",buf.String())
	if err != nil {
		glog.Fatal("mysql连接失败",err)
	}

	//设置连接池
	db.SetMaxOpenConns(cfg.Cfg.Mysql.BDmaxconn)
	db.SetMaxIdleConns(cfg.Cfg.Mysql.DBidleconn)

	err = db.Ping()
	if err != nil {
		glog.Fatal("mysql ping失败",err)
	}
	glog.Info("mysql连接成功")

	initService()
}


//设置数据库前缀
func tableName(name string) string {
	return name
}




var DBuser *userService


func initService() {
	DBuser = &userService{}
}