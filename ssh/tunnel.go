package ssh

import (
	"fmt"
	"time"
	"log"
	"context"
	"sort"
	"errors"
	"github.com/rgzr/sshtun"
	"github.com/xmonv/ssh-tunnel-manager/utils"
	"github.com/xmonv/ssh-tunnel-manager/database"
)

// Tunnels
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

// 全局变量用于存储所有配置
var ConfigMap = make(map[string]*TunnelConfig)
// 全局变量用于存储所有隧道
var TunnelMap = make(map[string]*sshtun.SSHTun)
// 全局变量用于存储所有隧道的日志
var TunnelLogMap = make(map[string]string)

func appendLogById(id string, msg string) {
	TunnelLogMap[id] += fmt.Sprintf("%s %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}


func (t *TunnelConfig) String() string {
	var left, right string
	mode := "<?>"
	switch t.Mode {
	case ">":
		left, mode, right = fmt.Sprintf("%s:%d", t.BindAddr, t.LocalPort), "->", fmt.Sprintf("%s:%d", t.DialAddr, t.RemotePort)
	case "<":
		left, mode, right = fmt.Sprintf("%s:%d", t.DialAddr, t.RemotePort), "<-", fmt.Sprintf("%s:%d", t.BindAddr, t.LocalPort)
	}
	return fmt.Sprintf("%s@%s | %s %s %s", t.User, t.Host, left, mode, right)
}


func DefaultTunnelConfig() *TunnelConfig {
	return &TunnelConfig {
		Id: utils.RandString(6),
		Name: "test",
		Mode: ">",
		User: "root",
		Host: "localhost",
		Port: 22,
		Password: "",
		BindAddr: "0.0.0.0",
		DialAddr: "0.0.0.0",
		LocalPort: 15244,
		RemotePort: 8080,
		CreateTime: time.Now().Unix(),
		Status: 0,
		Retry: 0,
		Toggle: 1,
	}
}

// TunnelConfig.Connect
// @Description     开启连接
// @Create             XMoNv
// @Param              TODO:过滤器
// @Return             隧道列表
func (t *TunnelConfig) Connect() {
	tun := sshtun.New(8080, t.Host, 80)
	tun.SetUser(t.User)
	tun.SetPort(t.Port)
	tun.SetLocalEndpoint(sshtun.NewTCPEndpoint(t.BindAddr, t.LocalPort))
	tun.SetRemoteEndpoint(sshtun.NewTCPEndpoint(t.DialAddr, t.RemotePort))
	if t.Password != "" {
		tun.SetPassword(t.Password)
	}
	t.Password = ""
	TunnelMap[t.Id] = tun
	ConfigMap[t.Id] = t
	TunnelLogMap[t.Id] = ""


	// We print each tunneled state to see the connections status
	// tun.SetTunneledConnState(func(tun *sshtun.SSHTun, state *sshtun.TunneledConnState) {
	// 	log.Printf("%+v", state)
	// })

	// We set a callback to know when the tunnel is ready
	tun.SetConnState(func(tun *sshtun.SSHTun, state sshtun.ConnState) {
		var msg string
		switch state {
		case sshtun.StateStarting:
			msg = fmt.Sprintf("%s STATE is Starting", t)
			t.Status = 1
		case sshtun.StateStarted:
			msg = fmt.Sprintf("%s STATE is Started", t)
			t.Status = 2
		case sshtun.StateStopped:
			msg = fmt.Sprintf("%s STATE is Stopped", t)
			t.Status = 3
			t.Retry = t.Retry + 1
		}
		log.Println(msg)
		appendLogById(t.Id, msg)
	})

	// We start the tunnel (and restart it every time it is stopped)
	go func() {
		for {
			if err := tun.Start(context.Background()); err != nil {
				msg := fmt.Sprintf("SSH tunnel error: %v", err)
				log.Println(msg)
				appendLogById(t.Id, msg)
				if t.Toggle == 1 {       // toggle == on
					if t.Retry < RetryMax {
						time.Sleep(time.Second) 
					} else {
						break          // don't flood if there's a start error :)
					}
				} else {                      // toggle = off
					break       
				}
			} else {
				if t.Toggle == 0 {
					break
				}
			}
		}
	}()
}


func CreateTunnel(tc *TunnelConfig, insert bool) {
	ConfigMap[tc.Id] = tc
	TunnelLogMap[tc.Id] = ""

	if insert {
		var db = database.GetDb()
		db.Create(tc)
	}
}

// GetTunnelList
// @Description     获取隧道列表，按照创建时间排序，方便展示
// @Create             XMoNv
// @Param              TODO:过滤器
// @Return             隧道列表
func GetTunnelList() [] *TunnelConfig {
	tclist := make([]*TunnelConfig, 0, len(ConfigMap))
	for _, v := range ConfigMap { 
		tclist = append(tclist, v)
	}
	sort.Slice(tclist, func(i, j int) bool {
		return tclist[i].CreateTime < tclist[j].CreateTime
	})
	return tclist
}


// GetTunnelById
// @Description     获取隧道相关信息
// @Create             XMoNv
// @Param              TODO:过滤器
// @Return             error, 隧道, 日志
func GetTunnelById(id string) (error, *TunnelConfig, string) {
	if tc, ok := ConfigMap[id]; ok {
		return nil, tc, TunnelLogMap[id]
	} else {
		msg := fmt.Sprintf("Id %s not find!", id)
		return errors.New(msg), nil, ""
	}
}

// StartTunnelById
// @Description     获取隧道相关信息
// @Create             XMoNv
// @Param              TODO:过滤器
// @Return             error
func StartTunnelById(id string) error {
	if _, ok := ConfigMap[id]; ok {
		if ConfigMap[id].Status != 2 {
			ConfigMap[id].Toggle = 1
			ConfigMap[id].Retry = 0
			ConfigMap[id].Connect()
		}
		return nil
	} else {
		msg := fmt.Sprintf("Id %s not find!", id)
		return errors.New(msg)
	}
}

// StopTunnelById
// @Description     获取隧道相关信息
// @Create             XMoNv
// @Param              TODO:过滤器
// @Return             error
func StopTunnelById(id string) error {
	if tun, ok := TunnelMap[id]; ok {
		if ConfigMap[id].Status == 2 {
			ConfigMap[id].Toggle = 0
			tun.Stop()
		}
		return nil
	} else {
		msg := fmt.Sprintf("Id %s not find!", id)
		return errors.New(msg)
	}
}

// DeleteTunnelById
// @Description     删除隧道相关信息
// @Create             XMoNv
// @Param              TODO:过滤器
// @Return             error
func DeleteTunnelById(id string) error {
	if _, ok := ConfigMap[id]; ok {
		if ConfigMap[id].Status == 2 {
			StopTunnelById(id)
		}
		delete(TunnelMap, id)
		delete(ConfigMap, id)
		delete(TunnelLogMap, id)

		var bean TunnelConfig
		var db = database.GetDb()
		db.Where("id = ?", id).Delete(&bean)
		return nil
	} else {
		msg := fmt.Sprintf("Id %s not find!", id)
		return errors.New(msg)
	}
}

// ModifyTunnelById
// @Description     修改隧道相关信息
// @Create             XMoNv
// @Param              TODO:过滤器
// @Return             error
func ModifyTunnelById(id string, tc *TunnelConfig) error {
	if tc, ok := ConfigMap[id]; ok {
		go func() {
			tc.Toggle = 0
			DeleteTunnelById(id)
			time.Sleep(time.Second) // wait to disconnect
			tc.Toggle = 1
			tc.Retry = 0

			var db = database.GetDb()
			db.Save(tc)
			tc.Connect()
		}()
		return nil
	} else {
		msg := fmt.Sprintf("Id %s not find!", id)
		return errors.New(msg)
	}
}