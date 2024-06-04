package server

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/xmonv/ssh-tunnel-manager/ssh"
	"github.com/xmonv/ssh-tunnel-manager/database"
	"net/http"
	"github.com/gin-gonic/gin"
)


func DbInit() {
	// tc := ssh.DefaultTunnelConfig()
	// tc.Connect()

	database.Init()
	var db = database.GetDb()

	var confs [] ssh.TunnelConfig
	db.Find(&confs)

	for _, tc := range confs {
		tc.Status = 0
		tc.Retry = 0
		ssh.CreateTunnel(&tc, false)
		logrus.Info(tc.String())
		if tc.Toggle == 1 {
			tc.Connect()
		}
	}
	logrus.Info("Db init Finish!")
}


func Init(app *gin.Engine) {
	//DbInit()

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// V1 API'
	api := app.Group("/api")
	v1 := api.Group("/v1") 
	v1.GET("list",func(c *gin.Context) {
		tlist := ssh.GetTunnelList()
		c.JSON(http.StatusOK, gin.H{"configs":tlist})
	})

	v1.GET("info/:id",func(c *gin.Context) {
		err, tc, msg := ssh.GetTunnelById(c.Param("id"))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"config": tc, "log": msg, "code": 0})
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 1})
		}
	})

	v1.POST("create",func(c *gin.Context) {
		tc := ssh.DefaultTunnelConfig()
		//var tc ssh.TunnelConfig
		if err := c.ShouldBind(&tc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ssh.CreateTunnel(tc, true)

		if tc.Toggle == 1 {
			tc.Connect()
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Success",  "code": 0, "config": tc})
	})

	v1.POST("start/:id",func(c *gin.Context) {
		id := c.Param("id")
		err, tc, _ := ssh.GetTunnelById(id)
		if err == nil {
			err = ssh.StartTunnelById(id)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"config": tc, "code": 0})
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 1})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 2})
		}
	})

	v1.POST("stop/:id",func(c *gin.Context) {
		id := c.Param("id")
		err, tc, _ := ssh.GetTunnelById(id)
		if err == nil {
			err = ssh.StopTunnelById(id)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"config": tc, "code": 0})
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 1})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 2})
		}
	})

	v1.POST("modify/:id",func(c *gin.Context) {
		id := c.Param("id")
		err, tc, _ := ssh.GetTunnelById(id)
		if err == nil {
			if err := c.ShouldBind(&tc); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			err = ssh.ModifyTunnelById(id, tc)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"config": tc, "code": 0})
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 1})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 2})
		}
	})

	v1.POST("delete/:id",func(c *gin.Context) {
		id := c.Param("id")
		err, tc, _ := ssh.GetTunnelById(id)
		if err == nil {
			err = ssh.DeleteTunnelById(id)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"config": tc, "code": 0})
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 1})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": err,  "code": 2})
		}
	})
}