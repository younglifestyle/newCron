package util

import (
	"reflect"
	"runtime"
)

//InStringArray if str in arr, return index. or return -1
func InStringArray(arr []string, str string) int {
	for idx, item := range arr {
		if item == str {
			return idx
		}
	}

	return -1
}

// FunctionName ...
func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
