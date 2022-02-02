package main

import (
	_ "Cube-back/cube"
	_ "Cube-back/database"
	_ "Cube-back/elasticsearch"
	_ "Cube-back/gron"
	_ "Cube-back/log"
	_ "Cube-back/message"
	_ "Cube-back/models"
	_ "Cube-back/rabbitmq"
	_ "Cube-back/redis"
	_ "Cube-back/routers"
	_ "Cube-back/snowflake"
	_ "Cube-back/ssh"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/beego/beego/v2/server/web/session"
)

var GlobalSessions *session.Manager

func insertFilter() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-requested-with"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:*", "http://127.0.0.1:*", "http://www.cube.fan:*"},
	}))
}

func sessionInit() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "CubeSessionId",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "./tmp",
	}
	GlobalSessions, _ = session.NewManager("memory", sessionConfig)
	go GlobalSessions.GC()
}

func main() {
	insertFilter()
	sessionInit()
	beego.Run()
}
