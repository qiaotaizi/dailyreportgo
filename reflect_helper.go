package main

import (
	"fmt"
	"reflect"
)

//反射帮助

//递归遍历一个结构体Value的所有字段
//structValue:结构体变量指针的Element
//hook: 对每个字段进行操作的函数,参数分别为成员属性的类型与值,当该函数返回error时,遍历会终止,并返回这个字段的名称与这个,error
//在钩子函数中,你只需要处理基本类型参数
func iterateStructField(structValue reflect.Value, hook func(fieldType reflect.StructField, fieldValue *reflect.Value) (string, error)) (string, error) {

	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		ft := structType.Field(i)
		fv := structValue.Field(i)
		fmt.Println("----", fv.Kind())
		if fv.Kind() == reflect.Struct {
			if failField, err := iterateStructField(fv, hook); err != nil {
				return failField, err
			}
			continue
		}
		if failField, err := hook(ft, &fv); err != nil {
			return failField, err
		}
	}
	return "", nil

}
