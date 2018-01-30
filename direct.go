package coven

import (
	"github.com/petersunbag/coven/basic"
	"reflect"
	"unsafe"
)

// directConverter handles converting among convertible basic types,
// and of the identical struct type.
type directConverter struct {
	*convertType
	cvtOp basic.CvtOp
	size  uintptr
}

var intAlign = unsafe.Alignof(int(1))

func newDirectConverter(convertType *convertType) (c converter) {
	st := convertType.srcTyp
	dt := convertType.dstTyp

	if cvtOp := basic.GetCvtOp(st, dt); cvtOp != nil {
		c = &directConverter{
			convertType: convertType,
			cvtOp:       cvtOp,
		}
		return
	}

	if st == dt && st.Kind() == reflect.Struct {
		c = &directConverter{
			convertType: convertType,
			size:        st.Size(),
		}
		return
	}

	return
}

// convert assigns converted source value to target.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by delegateConverter.Convert() and elemConverter.convert()
func (g *directConverter) convert(dPtr, sPtr unsafe.Pointer) {
	if g.cvtOp != nil {
		g.cvtOp(sPtr, dPtr)
	} else { // same struct type
		size := g.size
		align := uintptr(0)
		for ; align < size-intAlign; align += intAlign {
			*(*int)(unsafe.Pointer(uintptr(dPtr) + align)) = *(*int)(unsafe.Pointer(uintptr(sPtr) + align))
		}
		for ; align < size; align++ {
			*(*byte)(unsafe.Pointer(uintptr(dPtr) + align)) = *(*byte)(unsafe.Pointer(uintptr(sPtr) + align))
		}
	}
}
