package modelloader

import (
	"fmt"
	"reflect"
	"strconv"
)

const ptr = "ptr"

//ModelResult 解析模型结果
type ModelResult struct {
	modelName string
	appName   string

	fields   []ModelFeild
	numField uint

	PkValue ModelFeild
}

//ModelFeild 模型中的模板
type ModelFeild struct {
	isPk bool

	feildType reflect.Type

	rawValue      reflect.Value
	formatedValue string

	feildIndex    uint
	feildName     string
	feildShowName string

	maxSize       uint
	htmlInputType string

	autoGenerate bool
	generateWay  string

	sha256Hash bool
	editAble   bool
}

//HTMLModel for template render
type HTMLModel struct {
	ShowName      string
	Name          string
	HTMLType      string
	FormatedValue string

	MaxLen uint
}

//NewModel 通过给定的
func NewModel(target interface{}, appName string) ModelResult {
	var models ModelResult
	targetType := typeElemOutPtr(reflect.TypeOf(target))
	targetValue := valueElemOutPtr(reflect.ValueOf(target))

	models.modelName = targetType.Name()
	models.appName = appName
	models.numField = uint(targetType.NumField())

	for i := 0; i < int(models.numField); i++ {
		feildType := targetType.Field(i)
		feildValue := targetValue.Field(i)
		model := newModelFeild(feildType, feildValue, uint(i))

		if model.isPk {
			models.PkValue = model
		}

		models.fields = append(models.fields, model)
	}

	return models
}

func typeElemOutPtr(t reflect.Type) reflect.Type {
	if t.Kind().String() != ptr {
		return t
	}
	return typeElemOutPtr(t.Elem())
}

func valueElemOutPtr(t reflect.Value) reflect.Value {
	if t.Kind().String() != ptr {
		return t
	}
	return valueElemOutPtr(t.Elem())
}

func newModelFeild(t reflect.StructField, v reflect.Value, index uint) ModelFeild {
	var model ModelFeild
	model.feildType = t.Type
	model.rawValue = v

	model.feildName = t.Type.Name()
	model.feildIndex = index

	gormTagLoad(t, &model)
	adminTagLoad(t, &model)
	pkFind(&model)

	model.formatedValue = dataFormat(model.rawValue)

	return model
}

func dataFormat(t reflect.Value) string {
	dataType := t.Kind()

	if dataType == reflect.Bool {
		return strconv.FormatBool(t.Bool())
	} else if dataType >= reflect.Int && dataType <= reflect.Int64 {
		return strconv.Itoa(int(t.Int()))
	} else if dataType >= reflect.Uint && dataType <= reflect.Uint64 {
		return strconv.Itoa(int(t.Uint()))
	} else if dataType == reflect.Float32 {
		f:= fmt.Sprintf("%.2f",t.Float())
		return f
	} else if dataType == reflect.Float64 {
		f:= fmt.Sprintf("%.2f",t.Float())
		return f
	}else {
		return t.String()
	}
}

func pkFind(model *ModelFeild) {
	if model.feildName == "Model" {
		model.feildName = "id"
		model.isPk = true

		mod := model.rawValue.Interface()

		model.rawValue = reflect.ValueOf(mod).FieldByName("ID")
		model.feildType = reflect.TypeOf(model.rawValue.Uint())
	}

}

func (model *ModelFeild)isRequired()bool{
	return !(model.autoGenerate ||model.isPk)
}

func (model *ModelFeild) toHTMLData() (bool, HTMLModel) {
	if !model.isRequired() {
		return false, HTMLModel{}
	}
	return true, HTMLModel{
		Name:          model.feildName,
		ShowName:      model.feildShowName,
		HTMLType:      model.htmlInputType,
		FormatedValue: model.formatedValue,
		MaxLen:        model.maxSize,
	}
}

//HtMLTemplateData generlate for template
func (model *ModelResult) HtMLTemplateData() map[string]interface{} {
	var data map[string]interface{} = make(map[string]interface{})

	data["Name"] = model.modelName
	data["App"] = model.appName
	data["Pk"] = model.PkValue.formatedValue

	var feilds []HTMLModel
	for _, v := range model.fields {
		isOk, target := v.toHTMLData()
		if isOk {
			feilds = append(feilds, target)
		}
	}
	data["feilds"] = feilds

	return data
	
}
