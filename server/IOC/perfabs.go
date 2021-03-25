package IOC

import (
	"meetroom/err"
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

//部分预制件
//pathValue 路径参数，自动转换
type Value struct {
	value reflect.Value
}

//Get send the value into given data pointer
func (pv *Value) Get(target interface{}) err.Exception {
	inValue := reflect.ValueOf(target)
	//一级指针
	rawValue := inValue.Elem()

	//试图赋值
	if rawValue.Type() == pv.value.Type() {
		rawValue.Set(pv.value)
		return err.NoExcetion
	}
	return err.UnSupportData(
		rawValue.Kind().String() +
			" not thr same with " +
			pv.value.String())
}

func setValue(data reflect.Value) reflect.Value {
	t := Value{value: data}
	return reflect.ValueOf(t)
}

type TenmplateData struct {
	data map[string]interface{}
}

func newTemplateMap() *TenmplateData {
	t := TenmplateData{}
	t.data = make(map[string]interface{})
	return &t
}
func (t *TenmplateData) Set(key string, data interface{}) {
	t.data[key] = data
}

type ConxtextSeter struct {
	data map[string]interface{}
}

func newConxtextSeter() *ConxtextSeter {
	t := ConxtextSeter{}
	t.data = make(map[string]interface{})
	return &t
}
func (t *ConxtextSeter) Set(key string, data interface{}) {
	t.data[key] = data
}

type CookieSeter struct {
	handle CookieSetFunc
}

func newCookieSeter(c *gin.Context) *CookieSeter {
	t := CookieSeter{}
	t.handle = c.SetCookie
	return &t
}
func (t *CookieSeter) Set() CookieSetFunc {
	return t.handle
}

type RedisHandle struct {
	conn redis.Conn
}

func (r *RedisHandle) GetAsString(key string) (string, error) {
	return redis.String(r.conn.Do("GET", key))
}
func (r *RedisHandle) GetAsStringMap(key string) (map[string]string, error) {
	return redis.StringMap(r.conn.Do("GET", key))
}
func (r *RedisHandle) SetRecord(userID uint, info string) {
	//if user not exist create it

	_, e := r.conn.Do("LPUSH", userID, info)
	if e != nil {
		panic(e)
	}
	_, e = r.conn.Do("LTRIM",userID, 0, 9)
	if e != nil {
		panic(e)
	}

}
func (r *RedisHandle) GetAllRecord(userID uint) []string {
	var records []string
	records, e := redis.Strings(r.conn.Do("LRANGE", userID,0,-1))
	if e != nil {
		panic(e)
	}
	return records
}

type HtmlRender struct {
	htmlName string

	
}