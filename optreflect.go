package dbmodel

import (
	"reflect"
	"fmt"
	"unsafe"
)

type OptReflect struct{
	structName string
	fieldsMap  map[string]field
}
type field struct {
	offset uintptr
	fieldType string
}
type empty struct {
	etype  *struct{}
	ptr unsafe.Pointer
}
/**
获取表名称
 */
func (t *OptReflect) GetStructName() string{
	return t.structName
}
/**
获取值
 */
func (t *OptReflect) Get(obj interface{}, key string) (interface{},error) {
	if t.fieldsMap == nil{
		return nil,error(fmt.Errorf("对象尚未初始化"))
	}
	v,exist:=t.fieldsMap[key]
	if !exist{
		return nil,error(fmt.Errorf("%s字段不存在",key))
	}
	// 非指
	on:=reflect.TypeOf(obj).Name()
	if on==""{
		// 如果传入引用类型
		on=reflect.TypeOf(obj).Elem().Name()
	}
	if on!=t.structName{
		return nil,error(fmt.Errorf("给出的类型是%s,要求的类型为%s",on,t.structName))
	}
	ptr := (*empty)(unsafe.Pointer(&obj)).ptr
	ptr = unsafe.Pointer(uintptr(ptr) + v.offset)
	return getInterfaceType(uintptr(ptr),v.fieldType),nil
}
/**
设置值
 */
func (t *OptReflect) Set(obj interface{}, key string, value interface{}) error{
	if t.fieldsMap == nil{
		return error(fmt.Errorf("对象尚未初始化"))
	}
	p:=reflect.TypeOf(obj)
	if p.Kind().String()!="ptr"{
		return error(fmt.Errorf("要求传入指针类型,请检查是否忽略了&引用"))
	}
	n:=p.Name()
	if n==""{
		n=p.Elem().Name()
	}
	if n!=t.structName{
		return error(fmt.Errorf("当前传入类型%s,需要传入%s类型结构体的引用",n,t.structName))
	}
	// todo
	v,exist:=t.fieldsMap[key]
	if !exist{
		return error(fmt.Errorf("%s字段不存在",key))
	}
	// 非指针
	on:=reflect.TypeOf(obj).Name()
	if on==""{
		// 如果传入引用类型
		on=reflect.TypeOf(obj).Elem().Name()
	}
	if on!=t.structName{
		return error(fmt.Errorf("给出的类型是%s,要求的类型为%s",on,t.structName))
	}
	ptr := (*empty)(unsafe.Pointer(&obj)).ptr
	ptr = unsafe.Pointer(uintptr(ptr) + v.offset)
	switch v.fieldType {
	// 基本类型
	case "string":
		v,e:=value.(string)
		if !e{
			return error(fmt.Errorf("断言错误"))
		}
		*(* string)(unsafe.Pointer(ptr)) = v
	case "int":
		v,e:=value.(int)
		if !e{
			return error(fmt.Errorf("断言错误"))
		}
		*(* int)(unsafe.Pointer(ptr)) = v
	case "uint8":
		v,e:=value.(byte)
		if !e{
			return error(fmt.Errorf("断言错误"))
		}
		*(* byte)(unsafe.Pointer(ptr)) = v
		// 切片
	case "slice-string":
		//*(unsafe.Pointer(ptr)) = value.(uintptr)
	case "slice-int":
		//*(unsafe.Pointer(ptr)) = value.(uintptr)
	case "slice-uint8":
		//*(unsafe.Pointer(ptr)) = value.(uintptr)
		// 指针
	case "ptr-string":
		v,e:=value.(string)
		if !e{
			return error(fmt.Errorf("断言错误"))
		}
		*(* string)(unsafe.Pointer(ptr)) = v
	case "ptr-int":
		v,e:=value.(int)
		if !e{
			return error(fmt.Errorf("断言错误"))
		}
		*(* int)(unsafe.Pointer(ptr)) = v
	case "ptr-uint8":
		v,e:=value.(byte)
		if !e{
			return error(fmt.Errorf("断言错误"))
		}
		*(* byte)(unsafe.Pointer(ptr)) = v
	default:
		return error(fmt.Errorf("暂不支持对%s进行赋值",v.fieldType))
	}
	return nil
}
/**
使用前初始化,字段别名可以通过tag中的alias设置
如:
type Test struct {
	field1 string
	field2 []string `alias:"oo"`
	field3 int
}
 */
func (t *OptReflect) Init(model interface{}) {
	p:=reflect.TypeOf(model)
	if p.Kind().String()!="ptr"{
		panic(fmt.Errorf("要求传入指针类型,请检查是否忽略了&引用"))
	}
	elem:=p.Elem()
	if elem.Kind().String()!="struct"{
		panic(fmt.Errorf("给出的类型是%s,要求的类型为%s",elem.Kind().String(),"struct"))
	}
	if elem.NumField()==0{
		panic(fmt.Errorf("%s不存在可用字段",elem.Kind().String()))
	}
	t.fieldsMap = make(map[string]field,elem.NumField())
	for i:=0;i< elem.NumField();i++{
		f:=elem.Field(i)
		key:=f.Name
		if _,exist:=t.fieldsMap[key];exist{
			t.fieldsMap = nil
			panic(fmt.Errorf("字段名%s被多次定义.请检查结构体%s中tag及field是否存在重复命名",key,elem.Name()))
		}
		kind:=f.Type.Kind().String()
		if kind=="slice"{
			kind = "slice-"+f.Type.Elem().Name()
		}else if kind=="ptr"{
			//kind = "ptr-"+f.Type.Elem().Name()
			panic(fmt.Errorf("暂不支持&%s等指针类型",f.Type.Elem().Name()))
		}
		t.fieldsMap[key] = field{f.Offset,kind}
		if alias:=f.Tag.Get("alias");alias!="" && key!=alias{
			if _,exist:=t.fieldsMap[alias];exist{
				t.fieldsMap = nil
				panic(fmt.Errorf("字段名%s被多次定义.请检查结构体%s中tag及field是否存在重复命名",alias,elem.Name()))
			}
			t.fieldsMap[alias] = field{f.Offset,kind}
		}
		//fmt.Println(fmt.Sprintf("field:name=%s,tag=%s,type=%s,kind=%s,offset=%d,pkgPath=%s",f.Name,f.Tag.Get("db"),f.Type.Name(),kind,f.Offset,f.PkgPath))
	}
	t.structName = elem.Name()
}
func getInterfaceType(ptr uintptr, t string) interface{}{
	switch t {
	// 基本类型
	case "string":
		return *(* string)(unsafe.Pointer(ptr))
	case "int":
		return *(* int)(unsafe.Pointer(ptr))
	case "uint8":
		return *(* byte)(unsafe.Pointer(ptr))
	// 切片
	case "slice-string":
		return *(* []string)(unsafe.Pointer(ptr))
	case "slice-int":
		return *(* []int)(unsafe.Pointer(ptr))
	case "slice-uint8":
		return *(* []byte)(unsafe.Pointer(ptr))
	// 指针
	case "ptr-string":
		return *(* string)(unsafe.Pointer(ptr))
	case "ptr-int":
		return *(* int)(unsafe.Pointer(ptr))
	case "ptr-uint8":
		return *(* byte)(unsafe.Pointer(ptr))
	}
	return nil
}
//
//switch t.Kind() {
//case reflect.Bool:
//return boolEncoder
//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//return intEncoder
//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
//return uintEncoder
//case reflect.Float32:
//return float32Encoder
//case reflect.Float64:
//return float64Encoder
//case reflect.String:
//return stringEncoder
//case reflect.Interface:
//return interfaceEncoder
//case reflect.Struct:
//return newStructEncoder(t)
//case reflect.Map:
//return newMapEncoder(t)
//case reflect.Slice:
//return newSliceEncoder(t)
//case reflect.Array:
//return newArrayEncoder(t)
//case reflect.Ptr:
//return newPtrEncoder(t)
//default:
//return unsupportedTypeEncoder
//}