package typeTransfrom

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func toRaw(s interface{}) reflect.Value {
	return reflect.ValueOf(s)
}

func toString(s interface{}) reflect.Value {
	return reflect.ValueOf(s.(string))
}

func toInt(s interface{}) reflect.Value {
	i, e := strconv.ParseInt(s.(string), 10, 0)
	var r reflect.Value
	if e != nil {
		panic(fmt.Errorf("failure transfrom init int64 | %v", e))
	}
	r = reflect.ValueOf(i)

	return r
}
func toUint(s interface{}) reflect.Value {
	i, e := strconv.ParseUint(s.(string), 10, 0)
	var r reflect.Value
	if e != nil {
		panic(fmt.Errorf("failure transfrom init Uint64 | %v", e))
	}
	r = reflect.ValueOf(i)

	return r
}

func toStringSlice(s interface{}) reflect.Value {
	strs := strings.Split(s.(string), ",")
	return reflect.ValueOf(strs)
}
func toUintSlice(s interface{}) reflect.Value {
	strs := strings.Split(s.(string), ",")
	var uints []uint64
	for _, v := range strs {
		i, e := strconv.ParseUint(v, 10, 0)
		if e != nil {
			panic(fmt.Errorf("failure transfrom init Uint64 | %v", e))
		}
		uints = append(uints, i)
	}
	return reflect.ValueOf(uints)
}
func toIntSlice(s interface{}) reflect.Value {
	strs := strings.Split(s.(string), ",")
	var ints []int64
	for _, v := range strs {
		i, e := strconv.ParseInt(v, 10, 0)
		if e != nil {
			panic(fmt.Errorf("failure transfrom init Int64 | %v", e))
		}
		ints = append(ints, i)
	}
	return reflect.ValueOf(ints)
}
