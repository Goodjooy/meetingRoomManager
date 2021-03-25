package typeTransfrom

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

//将字符串转换为指定类型
type TransfromHandler func(interface{}) reflect.Value
type GoterHandler func(*gin.Context) func(key string) (data interface{}, exist bool)

var typeTransfromMap map[string]TransfromHandler
var fromGotersMap map[string]GoterHandler

func InitTypeTrasnfrom() {
	typeTransfromMap = make(map[string]TransfromHandler)
	fromGotersMap = make(map[string]GoterHandler)

	typeTransfromMap["raw"] = toRaw
	typeTransfromMap["string"] = toString
	typeTransfromMap["int"] = toInt
	typeTransfromMap["uint"] = toUint
	typeTransfromMap["stringSlice"] = toStringSlice
	typeTransfromMap["intSlice"] = toIntSlice
	typeTransfromMap["uintSlice"] = toUintSlice

	fromGotersMap["form"] = fromPostForm
	fromGotersMap["path"] = fromURLPath
	fromGotersMap["query"] = fromQuery
	fromGotersMap["context"] = fromContext
	fromGotersMap["cookie"] = fromCookie
	fromGotersMap["multi"] = fromMulit
}
func AddTypeTransfrom(tag string, handle TransfromHandler) {
	typeTransfromMap[tag] = handle
}

func GetTransfromer(name string) TransfromHandler {
	t, ok := typeTransfromMap[name]
	if !ok {
		panic(fmt.Errorf("`%s` is not a exist Transfrom", name))
	}
	return t

}
func GetDataGoter(name string) GoterHandler {
	t, ok := fromGotersMap[name]
	if !ok {
		panic(fmt.Errorf("`%s` is not a exist Goter", name))
	}
	return t
}
