package gsworm

import (
	"fmt"
	"reflect"
)

type GswTblDeclarationErr struct {
	Table string
}

type GswTypeAssertionErr struct {
	Column       string
	ExpectedType reflect.Type
	ActualType   reflect.Type
}

func (ge GswTblDeclarationErr) Error() string {
	return fmt.Sprintf("Failed to declare %v\n", ge.Table)
}

func (ge GswTypeAssertionErr) Error() string {
	return fmt.Sprintf("Type assertion error has been occured.\nExpected Type:%v\nActualType:%v\n", ge.ExpectedType, ge.ActualType)
}
