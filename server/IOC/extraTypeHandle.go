package IOC

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CookieSetFunc func(name, value string, maxAge int, path, domain string, secure, httpOnly bool)

var extraHandler map[reflect.Type]func(*gin.Context, *gorm.DB, redis.Conn) reflect.Value
var request = reflect.TypeOf(&http.Request{})
var database = reflect.TypeOf(&gorm.DB{})
var valueTypeName = reflect.TypeOf(Value{})
var templateMapPtr = reflect.TypeOf(&TenmplateData{})
var contextSeterPtr = reflect.TypeOf(&ConxtextSeter{})
var cookieJarPtr = reflect.TypeOf(&CookieSeter{})
var redisConnectionPtr = reflect.TypeOf(&RedisHandle{})

func initExtraHandler() {
	extraHandler = make(map[reflect.Type]func(*gin.Context, *gorm.DB, redis.Conn) reflect.Value)

	extraHandler[request] = handleReuest
	extraHandler[database] = handleDB
	extraHandler[contextSeterPtr] = handleContextSeter
	extraHandler[cookieJarPtr] = cookieSeter
	extraHandler[redisConnectionPtr] = redisHandler

}
func handleReuest(c *gin.Context, db *gorm.DB, r redis.Conn) reflect.Value {
	return reflect.ValueOf(c.Request)
}

func handleDB(c *gin.Context, db *gorm.DB, r redis.Conn) reflect.Value {
	return reflect.ValueOf(db)
}
func handleContextSeter(c *gin.Context, db *gorm.DB, r redis.Conn) reflect.Value {
	v := newConxtextSeter()
	return reflect.ValueOf(v)
}
func cookieSeter(c *gin.Context, db *gorm.DB, r redis.Conn) reflect.Value {
	var f = newCookieSeter(c)
	return reflect.ValueOf(f)
}
func redisHandler(c *gin.Context, db *gorm.DB, r redis.Conn) reflect.Value {
	f := RedisHandle{conn: r}
	return reflect.ValueOf(&f)
}

func getHandler(t reflect.Type) func(*gin.Context, *gorm.DB, redis.Conn) reflect.Value {
	f, ok := extraHandler[t]
	if !ok {
		panic(fmt.Errorf("not handler found for %s", t.String()))
	}
	return f
}
