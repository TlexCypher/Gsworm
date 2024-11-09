package gsworm

import (
	"fmt"
	"log"
	"reflect"
)

type SQLType string

type GswType interface {
	SqlType() SQLType
	RefType() reflect.Type
	validate() bool
}

type VChar struct {
	refType  reflect.Type
	sqlType  SQLType
	Size     int32
	Value    *string
	Capacity int32
}

type Int struct {
	sqlType SQLType
	refType reflect.Type
	Value   *int32
}

type BInt struct {
	sqlType SQLType
	refType reflect.Type
	Value   *int64
}

func VCHAR(capacity int32) VChar {
	return VChar{sqlType: SQLType(fmt.Sprintf("VARCHAR(%v)", capacity)),
		refType: reflect.TypeOf(string("")), Size: 0, Value: nil, Capacity: capacity}
}

func (v VChar) val(val string) VChar {
	v.Value = &val
	return v
}

func (vc VChar) validate() bool {
	if vc.Capacity < vc.Size {
		log.Printf("Capacity should be larger than size of VChar. Capacity:%v, Size:%v\n", vc.Capacity, vc.Size)
		return false
	}
	if len(*vc.Value) > int(vc.Size) {
		log.Println("Length of value should be less than the capacity and size of VChar.")
	}
	return true
}

func (vc VChar) SqlType() SQLType {
	return vc.sqlType
}

func (vc VChar) RefType() reflect.Type {
	return vc.refType
}

func INT() Int {
	return Int{sqlType: SQLType(fmt.Sprintf("INT")), refType: reflect.TypeOf(int32(0)), Value: nil}
}

func (i Int) val(value int32) Int {
	i.Value = &value
	return i
}

func (i Int) validate() bool {
	return true
}

func (i Int) SqlType() SQLType {
	return i.sqlType
}

func (i Int) RefType() reflect.Type {
	return i.refType
}

func BINT() BInt {
	return BInt{sqlType: SQLType(fmt.Sprintf("BInt")), refType: reflect.TypeOf(int64(0)), Value: nil}
}

func (bi BInt) val(val int64) BInt {
	bi.Value = &val
	return bi
}

func (bi BInt) validate() bool {
	return true
}

func (bi BInt) SqlType() SQLType {
	return bi.sqlType
}

func (bi BInt) RefType() reflect.Type {
	return bi.refType
}
