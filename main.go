package main

import (
	"log"
	"fmt"
	"flag"
	"github.com/xmonv/ssh-tunnel-manager/server"
	"github.com/gin-gonic/gin"
)


func main() {
	server.DbInit()
	app := gin.Default()
	server.Init(app)
	server.Static(app)

	port := flag.Int("port", 8080, "please input the server port")
	flag.Parse()
	log.Printf("listen on port %d", *port)
	log.Fatal(app.Run(fmt.Sprintf(":%d", *port)))
}