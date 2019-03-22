package lib

import (
	"container/list"
	"net/http"
)

// TOperation ternary operation
func TOperation(condition bool, trueOperation func() interface{}, falseOperation func() interface{}) interface{} {
	if condition {
		if trueOperation == nil {
			return nil
		}
		return trueOperation()
	}
	if falseOperation == nil {
		return nil
	}
	return falseOperation()
}

// TValue ternary operation
func TValue(condition bool, trueValue interface{}, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// WalkList walk a list
// walker return value as break signal
// if it is true, break walking
func WalkList(ls *list.List, walker func(item interface{}) bool) {
	if ls == nil {
		return
	}
	for ele := ls.Front(); ele != nil; ele = ele.Next() {
		breakWalk := walker(ele.Value)
		if breakWalk {
			break
		}
	}
}

// Try simulate try catch
func Try(f func(), catcher func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			catcher(err)
		}
	}()
	f()
}

// write http response
func WriteHttpResponse(writer http.ResponseWriter, status int, content string, headers map[string]string) {
	if headers != nil {
		for k, v := range headers {
			writer.Header().Add(k, v)
		}
	}
	writer.WriteHeader(status)
	if content != "" {
		writer.Write([]byte(content))
	}
}
