package gsworm

import "fmt"

type SQLType string

type GswType interface {
	SqlType() SQLType
}

type VChar struct {
	sqlType  SQLType
	Size     int32
	Value    *string
	capacity int32
}

type Int struct {
	sqlType SQLType
	Value   *int32
}

type BInt struct {
	sqlType SQLType
	Value   *int64
}

func VARCHAR(size int32) VChar {
	return VChar{sqlType: SQLType(fmt.Sprintf("VARCHAR(%v)", size)), Size: 0, Value: nil, capacity: size}
}

func (vc VChar) SqlType() SQLType {
	return vc.sqlType
}

func INT() Int {
	return Int{sqlType: "INT", Value: nil}
}

func (i Int) SqlType() SQLType {
	return i.sqlType
}

func BIGINT() BInt {
	return BInt{sqlType: "BIGINT", Value: nil}
}

func (bi BInt) SqlType() SQLType {
	return bi.sqlType
}
