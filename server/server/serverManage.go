package server

import (
	"fmt"
	"todo-web/server/IOC"
	"todo-web/server/manage"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Controller *gin.RouterGroup

var serverController *gin.Engine
var database *gorm.DB
var redisDB redis.Conn

func NewServer() {
	serverController = gin.Default()
}
func InitDatabase() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charsetutf8&parseTime=True&loc=Local", manage.SQLUser, manage.SQLPasswd, manage.DatabaseName))
	if err != nil {
		fmt.Printf("数据库连接失败，%s", err.Error())
		database = nil
	}
	database = db.Debug()

	return database
}
func InitRedis() redis.Conn {
	c, e := redis.Dial("tcp", "localhost:6379")
	if e != nil {
		fmt.Printf("Redis数据库连接失败，%s", e.Error())

		redisDB = nil
	}
	redisDB = c
	return redisDB
}
func UseIOCMid(fns ...interface{}) {
	var handles []gin.HandlerFunc
	for _, f := range fns {
		fn := IOC.ToIOC(f)
		handle := IOC.DoIOC(fn, database,&redisDB)
		handles = append(handles, handle)
	}
	UseMid(handles...)
}

func UseMid(handle ...gin.HandlerFunc) {
	serverController.Use(handle...)
}

func AddApplication(app manage.Application) {
	app.AsignApplication(serverController, database)
}

func NewController(path string) Controller {
	return serverController.Group(path)
}

func Build(address ...string) error {
	return serverController.Run(address...)
}
