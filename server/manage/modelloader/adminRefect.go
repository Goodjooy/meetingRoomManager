package modelloader

import (
	"reflect"
	"regexp"
	"strings"
)

const fieldName = `^.*name:([^;]+).*$`
const htmlType=`^.*type:([^;]+).*$`

//自动生成数据
const autoGenerateFeild = "autoG"
const uuidGenerateFeild = "uuidG"
const gormGenerateFeild = "gormG"

const editableFeild = "edab"

//要使用256哈希加密算法加密
const sha256Generate = "256G"

const autoGPattern = `^.*autoG:(uuidG|gormG).*$`

const adminForeignKeyPattern=`fk`

func adminTagLoad(target reflect.StructField, model *ModelFeild) {
	model.editAble = false
	model.sha256Hash = false
	model.autoGenerate = false
	model.generateWay = ""

	adminTags := target.Tag.Get("admin")

	nameP := regexp.MustCompile(fieldName)
	if nameP.MatchString(adminTags) {
		names := nameP.FindStringSubmatch(adminTags)
		if len(names) != 0 {
			model.feildShowName = names[1]
		}
	} else {
		model.feildShowName = model.feildName
	}

	autoP := regexp.MustCompile(autoGPattern)
	if autoP.MatchString(adminTags) {
		autoGs := autoP.FindStringSubmatch(adminTags)
		if len(autoGs) != 0 {
			model.autoGenerate = true
			model.generateWay = autoGs[1]
		}
	}

	if strings.Contains(adminTags, editableFeild) {
		model.editAble = true
	}

	if strings.Contains(adminTags, sha256Generate) {
		model.sha256Hash = true
	}

	typeP:=regexp.MustCompile(htmlType)
	if typeP.MatchString(adminTags){
		types:=typeP.FindStringSubmatch(adminTags)
		if len(types)!=0{
			model.htmlInputType=types[1]
		}
	}
}
