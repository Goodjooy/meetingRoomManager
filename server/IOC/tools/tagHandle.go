package tools

import (
	"fmt"
	"reflect"
	"strings"
)

func splitTagString(tag reflect.StructTag) []string {
	str, ok := tag.Lookup(tagName)
	if ok && str != "-" {
		return strings.Split(str, ";")
	} else {
		return []string{}
	}
}
func splitIntoMap(tag reflect.StructTag) map[string]string {
	tags := splitTagString(tag)
	var result map[string]string = make(map[string]string)
	for _, v := range tags {
		nAndV := strings.Split(v, ":")
		if len(nAndV) == 2 {
			result[nAndV[0]] = nAndV[1]
		} else if len(nAndV) == 1 {
			result[nAndV[0]] = "-"
		}
	}
	return result
}
func LoadTargetTypeTag(tag reflect.StructTag) (targetType string, from string, nameFrom string, reqire bool,defau string) {
	t := splitIntoMap(tag)
	typeName, tOK := t[targetTypeName]
	fromName, fOK := t[fromName]
	nameName, nOK := t[name]
	defData, require := t[defaultName]
	if !tOK {
		typeName = "raw"
	}
	if !fOK {
		panic(fmt.Errorf("data Source Not Special"))
	}
	if !nOK {
		panic(fmt.Errorf("data Name not Special"))
	}
	return typeName, fromName, nameName, !require,defData
}
