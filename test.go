package main


import (
	"strings"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type TunnelConfig struct {
	Id         string `json:"id", gorm="index:idx_id"`
	Name       string `form:"name" json:"name" binding:"required"`
	Mode        string `form:"mode" json:"mode" binding:"required,oneof=< >"`
	User        string `form:"user" json:"user" binding:"required" `
	Host        string `form:"host" json:"host" binding:"required"`
	Port        int  `form:"port" json:"port" binding:"required,min=1,max=65535"`
	Password string 
	BindAddr 	string `form:"bindAddr" json:"bindAddr" binding:"required"`
	DialAddr 	string `form:"dialAddr" json:"dialAddr" binding:"required"`
	LocalPort  int  `form:"localPort" json:"localPort" binding:"required,min=10000,max=65535", gorm="index:unique"`
	RemotePort int `form:"remotePort" json:"remotePort" binding:"required,min=1,max=65535"`
	CreateTime int64 `form:"createTime" json:"createTime"`
	Status int `form:"status" json:"status" `
	Retry  int `form:"retry" json:"retry"`
	Toggle int `form:"toggle" json:"toggle"`
}

var db *gorm.DB


func main() {
	file := "data.db"
	var err error

	if !(strings.HasSuffix(file, ".db") && len(file) > 3) {
		log.Fatalf("db name error.")
	}

	
	db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",file)), &gorm.Config{})

	//defer Close()

	if err != nil {
		log.Fatalf("failed to connect database:%s", err.Error())
	}

	err = db.AutoMigrate(&TunnelConfig{})
	if err != nil {
		log.Fatalf("failed migrate database: %s", err.Error())
	}

	var confs [] TunnelConfig
	db.Find(&confs)
	for _, tc := range confs {
		log.Println(tc)
	}
}