package modelloader

import (
	"reflect"
	"strconv"
	"todo-web/server/manage"

	"github.com/gin-gonic/gin"
)

//GetNeedFeilds 获取需要外部提供的变量
func (model *ModelResult) GetNeedFeilds() []ModelFeild {
	var needFeilds []ModelFeild
	for _, feild := range model.fields {
		if feild.isRequired() {
			needFeilds = append(needFeilds, feild)
		}
	}
	return needFeilds
}


//GetUUIDGenerateFeilds 获取需要使用uuid生成器的变量
func (model *ModelResult) GetUUIDGenerateFeilds() []ModelFeild {
	var needFeilds []ModelFeild
	for _, feild := range model.fields {
		if feild.generateWay == uuidGenerateFeild {
			feild.formatedValue = manage.UUIDGenerate()
			feild.rawValue = typeTransform(feild.formatedValue, feild.feildType)
			needFeilds = append(needFeilds, feild)
		}
	}
	return needFeilds
}
//CheckPostPram 检查post发送的内容，并对配对的模型赋值
func (model *ModelResult) CheckPostPram(c *gin.Context, targetModel interface{}) bool {
	needFeilds := model.GetNeedFeilds()

	for i := 0; i < len(needFeilds); i++ {
		feild := needFeilds[i]
		targetValue := c.PostForm(feild.feildName)
		if targetValue != "" {
			feild.formatedValue = targetValue

			if feild.sha256Hash {
				targetValue = manage.DateSHA256Hash(targetValue)
			}

			feild.rawValue = typeTransform(targetValue, feild.feildType)

			needFeilds[i] = feild
		} else {
			return false
		}
	}
	//set data
	target := reflect.ValueOf(targetModel)
	if target.Kind() != reflect.Ptr {
		return false
	}

	target = target.Elem()

	//given values
	for _, feild := range needFeilds {
		targetFeild := target.FieldByName(feild.feildName)
		if !targetFeild.IsValid() && targetFeild.CanSet() {
			targetFeild.Set(feild.rawValue)
		}
	}
	//uuid gengerate values
	for _, feild := range model.GetUUIDGenerateFeilds() {
		targetFeild := target.FieldByName(feild.feildName)
		if !targetFeild.IsValid() && targetFeild.CanSet() {
			targetFeild.Set(feild.rawValue)
		}
	}

	//set pk to nil
	pkFeild:=target.FieldByName(model.PkValue.feildName)
	if  !pkFeild.IsValid()&&pkFeild.CanSet(){
		pkFeild.Set(reflect.ValueOf(nil))
	}

	return true
}

func typeTransform(str string, typeOf reflect.Type) reflect.Value {
	kind := typeOf.Kind()
	var result reflect.Value
	if kind == reflect.Bool {
		r, _ := strconv.ParseBool(str)
		result = reflect.ValueOf(r)
	} else if kind >= reflect.Int && kind <= reflect.Int64 {
		i, _ := strconv.Atoi(str)
		result = reflect.ValueOf(i)
	} else if kind >= reflect.Uint && kind <= reflect.Uint64 {
		ui, _ := strconv.Atoi(str)
		result = reflect.ValueOf(uint64(ui))
	} else if kind == reflect.Float32 {
		f, _ := strconv.ParseFloat(str, 32)
		result = reflect.ValueOf(float32(f))
	} else if kind == reflect.Float64 {
		f, _ := strconv.ParseFloat(str, 64)
		result = reflect.ValueOf(f)
	} else {
		result = reflect.ValueOf(str)
	}
	return result
}
