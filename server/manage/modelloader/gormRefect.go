package modelloader

import (
	"reflect"
	"regexp"
	"strconv"
)

const sizePattern = `^.*size:(\d+).*$`
const sqlcharSizePattern = `^.*type:(?:varchar|char)\((\d+)\).*$`

const columnNamePattern = `^.*column:([^;]+)`

const defaultMaxLen = 128

func gormTagLoad(target reflect.StructField, model *ModelFeild) {
	model.maxSize = defaultMaxLen
	sizeP := regexp.MustCompile(sizePattern)
	typeP := regexp.MustCompile(sqlcharSizePattern)

	gormTags := target.Tag.Get("gorm")

	if sizeP.MatchString(gormTags) {
		num := sizeP.FindStringSubmatch(gormTags)
		if len(num) != 0 {
			number, err := strconv.Atoi(num[1])
			if err == nil {
				model.maxSize = uint(number)
			}
		}
	} else if typeP.MatchString(gormTags) {
		num := typeP.FindStringSubmatch(gormTags)
		if len(num) != 0 {
			number, err := strconv.Atoi(num[1])
			if err == nil {
				model.maxSize = uint(number)
			}
		}
	}

	model.feildName = target.Name
	columnNameP := regexp.MustCompile(columnNamePattern)
	if columnNameP.MatchString(gormTags) {
		values := columnNameP.FindStringSubmatch(gormTags)
		if len(values) != 0 {
			model.feildName = values[1]
		}
	}

	if !model.isPk {
		isOK, _ := regexp.MatchString(`^.*?primary_key.*?$`, gormTags)
		if isOK {
			model.isPk = true
		}
	}
}
