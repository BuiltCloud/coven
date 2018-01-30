def methodname(stype, dtype):
	return "cvt" + str.capitalize(stype) + str.capitalize(dtype)

def method(stype, dtype):
	if dtype == 'string':
		tempelate = '''
func %(methodname)s(sPtr, dPtr unsafe.Pointer) {
	*(*%(dtype)s)(dPtr) = fmt.Sprintf("%%v", *(*%(stype)s)(sPtr))
}
'''
	else:
		tempelate = '''
func %(methodname)s(sPtr, dPtr unsafe.Pointer) {
	*(*%(dtype)s)(dPtr) = (%(dtype)s)(*(*%(stype)s)(sPtr))
}
'''
	return tempelate % {'methodname': methodname(stype, dtype), 'stype':stype,'dtype':dtype}

def addCvtOP(stype, dtype):
	return "	cvtOps[convertKind{reflect.%(stype)s, reflect.%(dtype)s}] = %(methodname)s\n" % {'stype': str.capitalize(stype), 'dtype': str.capitalize(dtype), 'methodname': methodname(stype, dtype)}

def newT(type):
	return \
'''func new%(type)sPtr() unsafe.Pointer {
	v := %(v)s
	return unsafe.Pointer(&v)
}

''' % {'type':str.capitalize(type), 'v': defaultvalue[type]}

def new():
	s = '''
func newBasicValuePtr(k reflect.Kind) unsafe.Pointer {
	switch k {'''
	for t in typelist:
		s += '''
	case reflect.%(type)s:
		return new%(type)sPtr()''' % {'type':str.capitalize(t)}
	s +='''
	default:
		return nil
	}
}

'''
	for t in typelist:
		s += newT(t)
	return s

typelist = ['bool', 'int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'complex64', 'complex128', 'uintptr', 'string']
typeconvert = {
	'bool': ['bool', 'string'],
	'int' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'uint' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'int8' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'uint8' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'int16' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'uint16' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'int32' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'uint32' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'int64' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'uint64' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'float32' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'float64' : ['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'complex64' : ['complex64', 'complex128', 'string'],
	'complex128' : ['complex64', 'complex128', 'string'],
	'uintptr' :['int', 'uint', 'int8', 'uint8', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'uintptr', 'string'],
	'string' : ['string'],
}
defaultvalue = {
	'bool' : 'false',
	'int' : '0',
	'uint': '0',
	'int8': '0',
	'uint8': '0',
	'int16': '0',
	'uint16': '0',
	'int32': '0',
	'uint32': '0',
	'int64': '0',
	'uint64': '0',
	'float32': '0',
	'float64': '0',
	'complex64': '0',
	'complex128': '0',
	'uintptr': '0',
	'string': '""',
}

s = \
'''package coven

//this file is generated by gen.py

import (
	"unsafe"
	"reflect"
	"fmt"
)

type convertKind struct {
	srcTyp reflect.Kind
	dstTyp reflect.Kind
}

type cvtOp func(unsafe.Pointer, unsafe.Pointer)

var cvtOps = make(map[convertKind]cvtOp)

func init() {
'''

for t in typelist:
	for tt in typeconvert[t]:
		s += addCvtOP(t,tt)

s += '}\n\n'

for t in typelist:
	for tt in typeconvert[t]:
		s += method(t, tt)
s += '\n'

s += new()

with open('basicPtrCvt.go', 'w') as f:
	f.write(s)
