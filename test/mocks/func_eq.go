package mocks

import (
	"github.com/golang/mock/gomock"
	"reflect"
)

func FuncEq(x interface{}) gomock.Matcher { return funEq{x} }

type funEq struct {
	x interface{}
}

func (f funEq) Matches(x interface{}) bool {
	xv := reflect.ValueOf(x)
	fxv := reflect.ValueOf(f.x)
	return xv.Pointer() == fxv.Pointer()
}

func (funEq) String() string {
	return "is same func"
}
