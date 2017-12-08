## 前言
    该项目旨在利用指针偏移量记录类型结构，以提高反射操作效率。  
## 使用方法

首先将用于反射的结构对optreflect对象进行初始化，optreflect记录该结构属性特征。  
但是，使用调用Init的时候要保证传入的类型是可以被成功解析的，否则panic.
```go
// 用于反射的结构
type TestObj struct {
	field1 string	`alias:"f1"`
	field2 []string
	field3 int
	field4 []int
	field5 byte
	field6 []byte
}
optObj:=&optreflect.OptReflect{}
optObj.Init(&TestObj{})
// 
```
对给出的对象我们可以通过优化后的反射机制获取其属性值
```go
testObj:=TestObj{
        field1:"test string",
        field2:[]string{"a","b","c"},
        field3:4523,
        field4:[]int{1,2,3,5},
        field5:byte(4),
        field6:[]byte{2,2,2},
    }
    if v,err:=optObj.Get(&testObj,"field1");err==nil{
        fmt.Println(v)
    }
    if v,err:=optObj.Get(&testObj,"f1");err==nil{
        fmt.Println(v)
    }
    if v,err:=optObj.Get(&testObj,"field2");err==nil{
        fmt.Println(v)
    }
    if v,err:=optObj.Get(&testObj,"field3");err==nil{
        fmt.Println(v)
    }
    if v,err:=optObj.Get(&testObj,"field5");err==nil{
        fmt.Println(v)
    }
    if err:=optObj.Set(&testObj,"f1","alert data");err==nil{
        fmt.Println(testObj.field1)
    }
    /** 输出
    test string
    test string
    [a b c]
    4523
    4
     */
```
同时，我们可以设置对象的属性。
```go
if err:=optObj.Set(&testObj,"f1","alert data");err==nil{
    fmt.Println(testObj.field1)
}
/** 输出
    alert data
     */
```

## plus

该项目作为学习的尝试，存在着许多不足，包括其支持的数据类型并不多（后面会添加）.欢迎讨论指导.