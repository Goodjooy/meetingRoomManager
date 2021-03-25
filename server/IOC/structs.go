package IOC

import (
	"reflect"
)

//FuncHandler 接管函数的相关信息
type FuncHandler struct {
	//call able
	fn reflect.Value
	//func type
	fnType reflect.Type

	inNum uint
	//参数列表
	inArray []InHandler
}

type InHandler struct {
	//参数类型
	parmType reflect.Type
	//是否为结构体
	structType bool
	//结构体变量
	insideFeild []InFeildHandler
}

type InFeildHandler struct {
	name string
	//类型 Value 和基础类型 http request Respond
	feildType reflect.Type

	pkgPath string

	tag reflect.StructTag
	//目标转换类型
	targetType string
	from       string
	nameFrom   string

	reqire      bool
	defaultData string
}


