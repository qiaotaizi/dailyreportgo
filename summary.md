# 总结一下踩过的坑和知识点

1.数组/slice不能标为const  
2.os.Stat()可以用于判断文件是否存在  
3.reflect.Type()用于获取结构体变量的类型,进而获取成员的tag信息,
reflect.Value()用于获取变结构体量的值,进而获取成员的值
如果这两个函数的入参是指针变量,那么在进一步获取对象信息之前,需要调用.Elem()函数  
4.通过type.NumField()函数获取结构体成员的数量,结合for循环遍历成员  
5.可以对结构体成员自定义标签,在反射时使用  
6.必须使用以下方法发送post请求:
```go
...
form := url.Values{
    "param1": {param1},
    "param2": {param2},
    "param3": {param3},
    ...
}
formString := form.Encode()
req, err := http.NewRequest(http.MethodPost, reportConfig.JiraLoginUrl, strings.NewReader(formString))
...
```
需要有Header
```go
req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
```
7.时间格式化,格式字符串必须使用2006-01-02 15:04:05这个时间点  
8.测试文件以_test.go为后缀,测试方法以Test为前缀,入参为t *testing.T  
9.变长参数需要作为另外一个函数的入参时,使用args...的形式传入  
10.init()方法在同一包中的调用顺序是按照文件名字母序从a到z进行调用的  
