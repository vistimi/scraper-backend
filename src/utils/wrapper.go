package utils

// import (
// 	"reflect"

// 	"github.com/foolin/pagser"

// 	"errors"
// )

// // res := WrapperError(<fctName>, <param1>, <param2>)
// func WrapperParser(fn interface{}, parser *pagser.Pagser, dataType interface{}, params ...interface{}) (reflect.Type, error) {
// 	f := reflect.ValueOf(fn)
// 	if f.Type().NumIn() != len(params) {
// 		panic("incorrect number of parameters!")
// 	}
// 	inputs := make([]reflect.Value, len(params))
// 	for k, in := range params {
// 		inputs[k] = reflect.ValueOf(in)
// 	}
// 	res := f.Call(inputs)
// 	if err, ok := res[1].Interface().(error); ok && err != nil {
// 		return nil, err
// 	}
// 	dataTypeExtract := reflect.TypeOf(dataType)
// 	switch dataTypeExtract.Kind() {
// 	case reflect.Struct:
// 		data := reflect.New(reflect.TypeOf(dataType))
// 		err := parser.Parse(&data, res[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		if data.Stat != "ok" {
// 			return nil, err
// 		}
// 		return data, nil
// 	default:
// 		return nil, errors.New("Wrong dataType, should a be struct!")
// 	}
// }