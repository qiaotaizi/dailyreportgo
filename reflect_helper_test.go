package main

type s1 struct {
	s2
	F1 string `v:"f1"`
	F2 bool   `v:"true"`
	F3 int    `v:"100"`
}

type s2 struct {
	F4 string `v:"f4"`
	F5 int    `v:"200"`
}

//func TestIterateStructField(t *testing.T) {
//	s := new(s1)
//	failField, err := iterateStructField(reflect.ValueOf(s).Elem(),
//		func(fieldType reflect.StructField, fieldValue reflect.Value) (s string, e error) {
//			tagV := fieldType.Tag.Get("v")
//			switch fieldValue.Kind() {
//			case reflect.String:
//				fieldValue.SetString(tagV)
//			case reflect.Bool:
//				fieldValue.SetBool(tagV == "true")
//			case reflect.Int, reflect.Int32, reflect.Int64:
//				i, err := strconv.Atoi(tagV)
//				if err != nil {
//					return fieldType.Name, err
//				}
//				fieldValue.SetInt(int64(i))
//			default:
//				return fieldType.Name, fmt.Errorf("未能识别的类型: %s", fieldValue.Kind())
//			}
//			return "", nil
//		})
//	if err != nil {
//		fmt.Printf("遍历失败的字段: %s,异常为: %v\n", failField, err)
//	}
//	fmt.Println(s)
//}
