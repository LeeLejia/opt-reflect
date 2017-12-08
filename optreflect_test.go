package dbmodel

import (
	"testing"
	"optreflect"
)
type TestObj struct {
	field1 string	`alias:"f1"`
	field2 []string
	field3 int
	field4 []int
	field5 byte
	field6 []byte
}
func TestOptReflect(t *testing.T) {
	optObj:=&optreflect.OptReflect{}
	optObj.Init(&TestObj{})
	testObj:=TestObj{
		field1:"test string",
		field2:[]string{"a","b","c"},
		field3:4523,
		field4:[]int{1,2,3,5},
		field5:byte(4),
		field6:[]byte{2,2,2},
	}
	// 测试 string
	if v,err:=optObj.Get(&testObj,"field1");err!=nil || v!=testObj.field1{
		t.Error("failed")
	}
	// 测试别名
	if v,err:=optObj.Get(&testObj,"f1");err!=nil || v!=testObj.field1{
		t.Error("failed")
	}
	// 测试数组
	v,err:=optObj.Get(&testObj,"field2")
	if err!=nil{
		t.Error("failed")
	}else {
		for i,vv:=range v.([]string){
			if vv!=testObj.field2[i]{
				t.Error("failed")
			}
		}
	}
	// 测试int
	if v,err:=optObj.Get(&testObj,"field3");err!=nil || v!=testObj.field3{
		t.Error("failed")
	}
	// 测试 byte
	if v,err:=optObj.Get(&testObj,"field5");err!=nil || v!=testObj.field5{
		t.Error("failed")
	}
	// 测试修改 基本类型
	err=optObj.Set(&testObj,"f1","alert data")
	if err!=nil{
		t.Error(err.Error())
	}
	if testObj.field1!="alert data"{
		t.Error("failed")
	}
	// 测试修改切片类型
	// todo 待更新
	t.Log("test success!!")
}