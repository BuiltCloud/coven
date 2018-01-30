package coven

import (
	"reflect"
	"unsafe"
)

type sliceConverter struct {
	*convertType
	*elemConverter
}

func newSliceConverter(convertType *convertType) (s converter) {
	if elemConverter, ok := newElemConverter(convertType.dstTyp.Elem(), convertType.srcTyp.Elem()); ok {
		s = &sliceConverter{
			convertType,
			elemConverter,
		}
	}
	return
}

// convert will overwrite the whole target slice.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by delegateConverter.Convert() and elemConverter.convert()
func (s *sliceConverter) convert(dPtr, sPtr unsafe.Pointer) {
	dSlice := (*sliceHeader)(dPtr)
	sSlice := (*sliceHeader)(sPtr)

	length := sSlice.Len
	dSlice.Len = length

	if dSlice.Cap < length {
		newVal := reflect.MakeSlice(s.dstTyp, 0, length)
		dSlice.Data = unsafe.Pointer(newVal.Pointer())
		dSlice.Cap = length
	}

	for dOffset, sOffset, i := uintptr(0), uintptr(0), 0; i < length; i++ {
		dElemPtr := unsafe.Pointer(uintptr(dSlice.Data) + dOffset)
		sElemPtr := unsafe.Pointer(uintptr(sSlice.Data) + sOffset)
		s.elemConverter.convert(dElemPtr, sElemPtr)
		dOffset += s.elemConverter.dDereferSize
		sOffset += s.elemConverter.sDereferSize
	}
}
