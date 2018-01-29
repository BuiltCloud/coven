package coven

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var (
	createdConvertersMu sync.Mutex
	createdConverters   = make(map[convertType]*delegateConverter)
)

type convertType struct {
	dstTyp reflect.Type
	srcTyp reflect.Type
}

type converter interface {
	convert(dPtr, sPtr unsafe.Pointer)
}

type delegateConverter struct {
	*convertType
	converter
}

func (d *delegateConverter) Convert(dst, src interface{}) {
	dv := dereferencedValue(dst)
	if !dv.CanSet() {
		panic(fmt.Sprintf("target should be a pointer. [actual:%v]", dv.Type()))
	}

	if dv.Type() != d.dstTyp {
		panic(fmt.Sprintf("invalid target type. [expected:%v] [actual:%v]", d.dstTyp, dv.Type()))
	}

	sv := dereferencedValue(src)
	if !sv.CanAddr() {
		panic(fmt.Sprintf("source should be a pointer. [actual:%v]", dv.Type()))
	}

	if sv.Type() != d.srcTyp {
		panic(fmt.Sprintf("invalid source type. [expected:%v] [actual:%v]", d.srcTyp, sv.Type()))
	}

	d.converter.convert(unsafe.Pointer(dv.UnsafeAddr()), unsafe.Pointer(sv.UnsafeAddr()))
}

func NewConverter(dst, src interface{}) *delegateConverter {
	dstTyp := reflect.TypeOf(dst)
	srcTyp := reflect.TypeOf(src)

	if c := newConverter(dstTyp, srcTyp, true); c == nil {
		panic(fmt.Sprintf("can't convert source type %s to target type %s", srcTyp, dstTyp))
	} else {
		return c
	}
}

func newConverter(dstTyp, srcTyp reflect.Type, lock bool) *delegateConverter {
	if lock {
		createdConvertersMu.Lock()
		defer createdConvertersMu.Unlock()
	}

	dstTyp = dereferencedType(dstTyp)
	srcTyp = dereferencedType(srcTyp)

	cTyp := &convertType{dstTyp, srcTyp}
	if dc, ok := createdConverters[*cTyp]; ok {
		return dc
	}

	var c converter
	if c = newGeneralConverter(cTyp); c == nil {
		switch sk, dk := srcTyp.Kind(), dstTyp.Kind(); {

		case sk == reflect.Struct && dk == reflect.Struct:
			c = newStructConverter(cTyp)

		case sk == reflect.Slice && dk == reflect.Slice:
			c = newSliceConverter(cTyp)

		case sk == reflect.Map && dk == reflect.Map:
			c = newMapConverter(cTyp)

		default:
			return nil
		}
	}
	if c != nil {
		dc := &delegateConverter{cTyp, c}
		createdConverters[*cTyp] = dc
		return dc
	}

	return nil
}
