package reflect

import (
	"context"
	"fmt"
	"reflect"
)

/**
handle方法 第一个参数必须为context
 */
func HandleFun(handle interface{}, ctx context.Context, p ...interface{}) error {
	handleType := reflect.TypeOf(handle)
	if handleType.Kind() == reflect.Func {
		fmt.Println("input field num: ", handleType.NumIn())
		// 获取调用方法
		handleVal := reflect.ValueOf(handle)
		// 获取方法请求参数
		parameters := make([]reflect.Value, 0)
		parameters = append(parameters, reflect.ValueOf(ctx))
		for _, parameter := range p {
			parameters = append(parameters, reflect.ValueOf(parameter))
		}
		// 调用
		res := handleVal.Call(parameters)
		return checkErr(res)
	}
	return nil
}

/**
handle函数最后一个元素返回error
 */
func checkErr(result []reflect.Value) error {
	if len(result) > 0 && !result[len(result)-1].IsNil() {
		return result[len(result)-1].Interface().(error)
	}
	return nil
}
