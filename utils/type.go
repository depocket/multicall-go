package utils

import "reflect"

func TypeCast(src []interface{}, dst interface{}) interface{} {
	dstType := reflect.TypeOf(dst)
	dstValue := reflect.New(dstType)
	for i := 0; i < dstType.NumField(); i++ {
		dstValue.Elem().FieldByName(dstType.Field(i).Name).Set(reflect.ValueOf(src[i]))
	}
	return dstValue.Interface()
}
