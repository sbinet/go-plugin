package ffi

// #include "go-ffi.h"
// #include <stdint.h>
// #include <stdlib.h>
import "C"
import (
	"reflect"
	"unsafe"
)

var (
	ffiTypePtrSize = unsafe.Sizeof(C.uintptr_t(0))
)

type Type struct {
	name string
	c    *C.ffi_type
}

var (
	Void = Type{"void", &C.ffi_type_void}

	Uchar  = Type{"unsigned char", &C.ffi_type_uchar}
	Ushort = Type{"unsigned short", &C.ffi_type_ushort}
	Uint   = Type{"unsigned int", &C.ffi_type_uint}
	Ulong  = Type{"unsigned long", &C.ffi_type_ulong}

	Char  = Type{"char", &C.ffi_type_schar}
	Short = Type{"short", &C.ffi_type_sshort}
	Int   = Type{"int", &C.ffi_type_sint}
	Long  = Type{"long", &C.ffi_type_slong}

	Uint8  = Type{"uint8_t", &C.ffi_type_uint8}
	Uint16 = Type{"uint16_t", &C.ffi_type_uint16}
	Uint32 = Type{"uint32_t", &C.ffi_type_uint32}
	Uint64 = Type{"uint64_t", &C.ffi_type_uint64}

	Int8  = Type{"int8_t", &C.ffi_type_sint8}
	Int16 = Type{"int16_t", &C.ffi_type_sint16}
	Int32 = Type{"int32_t", &C.ffi_type_sint32}
	Int64 = Type{"int64_t", &C.ffi_type_sint64}

	Float32 = Type{"float32", &C.ffi_type_float}
	Float64 = Type{"float64", &C.ffi_type_double}

	UnsafePointer = Type{"void*", &C.ffi_type_pointer}
)

func NewReturnFrom(typ reflect.Type) unsafe.Pointer {
	err := validateReflectType(typ)
	if err != nil {
		panic(err)
	}
	return nil
}

type Args struct {
	c    *unsafe.Pointer
	args []reflect.Value
}

func (args *Args) C() unsafe.Pointer {
	if args.c != nil {
		return *args.c
	}
	return nil
}

func (args *Args) Release() {
	C.free(unsafe.Pointer(args.c))
	args.c = nil
}

func NewArgsFrom(args []reflect.Value) Args {
	n := len(args)
	if n <= 0 {
		return Args{}
	}

	out := Args{
		args: make([]reflect.Value, n),
		c:    (*unsafe.Pointer)(unsafe.Pointer(C.malloc(C.size_t(n) * C.size_t(ptrSize)))),
	}
	for i := range args {
		out.args[i] = reflect.New(args[i].Type()).Elem()
		out.args[i].Set(args[i])
		C._go_ffi_void_array(out.c, C.int(i), unsafe.Pointer(out.args[i].UnsafeAddr()))
	}
	return out
}
