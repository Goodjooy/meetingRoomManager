package IOC

import (
	"fmt"
	"reflect"
	"todo-web/dataHandle"
	"todo-web/err"
	"todo-web/server/IOC/tools"
	"todo-web/server/IOC/typeTransfrom"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//IOC 控制反转
var inited = false

func InitIOC() {
	typeTransfrom.InitTypeTrasnfrom()
	initExtraHandler()
	inited = true
}

func ToIOC(handleFunc interface{}) FuncHandler {
	if !inited {
		InitIOC()
	}

	valueType := reflect.TypeOf(handleFunc)
	if valueType.Kind() != reflect.Func {
		panic(fmt.Errorf("traget type is not Func %v", valueType.String()))
	}
	valueValue := reflect.ValueOf(handleFunc)

	var fun FuncHandler

	fun.fn = valueValue
	fun.fnType = valueType
	fun.inNum = uint(valueType.NumIn())
	fun.inArray = make([]InHandler, 0)
	ParmCount := fun.inNum

	//循环查找全部参数
	for i := 0; uint(i) < ParmCount; i++ {
		var inHandler InHandler
		parm := valueType.In(i)

		inHandler.parmType = parm
		inHandler.structType = parm.Kind() == reflect.Struct

		if inHandler.structType {
			//循环处理全部结构体变量
			for i := 0; i < parm.NumField(); i++ {
				var feild InFeildHandler
				f := parm.Field(i)
				feild.feildType = f.Type
				feild.name = f.Name
				feild.pkgPath = f.PkgPath

				fName := f.Type.String()
				valueName := valueTypeName.String()

				if fName == valueName {
					feild.tag = f.Tag
					feild.targetType,
						feild.from,
						feild.nameFrom,
						feild.reqire,
						feild.defaultData = tools.LoadTargetTypeTag(f.Tag)
				}
				inHandler.insideFeild =
					append(inHandler.insideFeild, feild)
			}
		}

		fun.inArray = append(fun.inArray, inHandler)
	}

	return fun
}
func generateInArray(fn FuncHandler, db *gorm.DB, c *gin.Context, rC redis.Conn) (argList []reflect.Value, setterList []reflect.Value) {
	var valueList []reflect.Value = make([]reflect.Value, 0)
	var ContextSeters []reflect.Value

	for _, v := range fn.inArray {
		//遍历参数列表
		var d reflect.Value
		if v.structType {
			//参数类型为结构体
			st := reflect.New(v.parmType).Elem()

			for _, feild := range v.insideFeild {
				targetFeild := st.FieldByName(feild.name)
				if feild.feildType == reflect.TypeOf(Value{}) {
					goter := typeTransfrom.GetDataGoter(feild.from)
					targetStr, exist := goter(c)(feild.nameFrom)

					if !feild.reqire && !exist {
						targetStr = feild.defaultData
					}
					if feild.reqire && !exist {
						panic(fmt.Errorf("target parm require but not exist | from: %s | name: %s | type: %s",
							feild.from, feild.name, feild.targetType))
					}

					transfrom := typeTransfrom.GetTransfromer(feild.targetType)
					targetFeild.Set(setValue(transfrom(targetStr)))
				} else {
					var t reflect.Value = getHandler(feild.feildType)(c, db, rC)
					if t.Type() == contextSeterPtr {
						ContextSeters = append(ContextSeters, t)
					}

					targetFeild.Set(t)
				}
			}
			d = st
		} else {
			d = getHandler(v.parmType)(c, db, rC)
			if d.Type() == contextSeterPtr {
				ContextSeters = append(ContextSeters, d)
			}
		}

		valueList = append(valueList, d)
	}
	return valueList, ContextSeters
}

func DoIOC(fn FuncHandler, db *gorm.DB, rC *redis.Conn) gin.HandlerFunc {
	if !inited {
		InitIOC()
	}

	var fun gin.HandlerFunc = func(c *gin.Context) {
		//函数参数列表
		defer recoverHandle("run func", fn, c)

		valueList, contextList := generateInArray(fn, db, c, *rC)
		//todo defer recover
		results := fn.fn.Call(valueList)

		for _, v := range contextList {
			t := v.Interface().(*ConxtextSeter)
			for key, value := range t.data {
				c.Set(key, value)
			}
		}
		//todo handle results
		var r []interface{}

		for _, v := range results {
			t := v.Interface()
			r = append(r, t)
		}
		if len(r) == 0 {
		} else if len(r) == 1 {
			c.JSON(200, r[0])
		} else {
			c.JSON(200, r)
		}
	}
	return fun
}

func recoverHandle(local string, fn FuncHandler, c *gin.Context) {
	e := recover()
	if e == nil {
		return
	}

	var r dataHandle.Result = dataHandle.FailureFuncResult(
		err.FailureGenerateFunctionParm,
		fmt.Sprintf("[%s] : %s | %v", local, fn.fnType.Name(), e))

	//todo logger

	c.JSON(400, r)
	c.Abort()

}
