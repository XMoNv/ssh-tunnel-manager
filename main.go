package main

import (
	"log"
	"github.com/xmonv/ssh-tunnel-manager/server"
	"github.com/gin-gonic/gin"
)


func main() {
	server.DbInit()
	app := gin.Default()
	server.Init(app)
	server.Static(app)
	log.Fatal(app.Run(":8080"))
}