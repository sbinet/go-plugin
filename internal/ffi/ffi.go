package ffi

// #cgo pkg-config: libffi
//
// #include <stdlib.h>
// #include <string.h>
// #include "go-ffi.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

var (
	ptrSize = unsafe.Sizeof((*uintptr)(nil))
)

// ABI is the FFI ABI for the local platform
type ABI C.ffi_abi

const (
	FirstABI   ABI = C.FFI_FIRST_ABI
	DefaultABI ABI = C.FFI_DEFAULT_ABI
	LastABI    ABI = C.FFI_LAST_ABI
)

const (
	TrampolineSize = C.FFI_TRAMPOLINE_SIZE
	NativeRawAPI   = C.FFI_NATIVE_RAW_API

	//Closures          = C.FFI_CLOSURES
	//TypeSmallStruct1B = C.FFI_TYPE_SMALL_STRUCT_1B
	//TypeSmallStruct2B = C.FFI_TYPE_SMALL_STRUCT_2B
	//TypeSmallStruct4B = C.FFI_TYPE_SMALL_STRUCT_4B
)

type Status uint32

const (
	Ok         Status = C.FFI_OK
	BadTypedef Status = C.FFI_BAD_TYPEDEF
	BadABI     Status = C.FFI_BAD_ABI
)

func (sc Status) String() string {
	switch sc {
	case Ok:
		return "FFI_OK"
	case BadTypedef:
		return "FFI_BAD_TYPEDEF"
	case BadABI:
		return "FFI_BAD_ABI"
	}
	panic("ffi: unreachable")
}

// CIF is the ffi call interface
type CIF struct {
	c    C.ffi_cif
	ret  *C.ffi_type
	args **C.ffi_type
}

// New creates a new ffi call interface object
func New(rtype Type, args []Type) (*CIF, error) {
	cif := &CIF{
		ret:  rtype.c,
		args: nil,
	}

	n := C.size_t(len(args))
	if n > 0 {
		size := n * C.size_t(ptrSize)
		cif.args = (**C.ffi_type)(unsafe.Pointer(C.malloc(size)))
		for i, arg := range args {
			C._go_ffi_type_array(cif.args, C.int(i), arg.c)
		}
	}
	sc := C.ffi_prep_cif(&cif.c, C.FFI_DEFAULT_ABI, C.uint(n), cif.ret, cif.args)
	if sc != C.FFI_OK {
		return nil, fmt.Errorf("error while preparing cif (%s)",
			Status(sc))
	}
	return cif, nil
}

func NewFrom(rt reflect.Type) (*CIF, error) {
	switch rt.Kind() {
	case reflect.Func:

		var cif CIF

		for i := 0; i < rt.NumIn(); i++ {
			typ := rt.In(i)
			err := validateReflectType(typ)
			if err != nil {
				return nil, err
			}
		}
		for i := 0; i < rt.NumOut(); i++ {
			if i > 1 {
				return nil, fmt.Errorf("ffi: too many return types in function signature")
			}
			typ := rt.Out(i)
			err := validateReflectType(typ)
			if err != nil {
				return nil, err
			}
		}

		if rt.NumOut() == 1 {
			cif.ret = ffiTypeFrom(rt.Out(0)).c
		}

		n := C.size_t(rt.NumIn())
		if n > 0 {
			size := n * C.size_t(ptrSize)
			cif.args = (**C.ffi_type)(unsafe.Pointer(C.malloc(size)))
			for i := 0; i < int(n); i++ {
				arg := ffiTypeFrom(rt.In(i))
				C._go_ffi_type_array(cif.args, C.int(i), arg.c)
			}
		}
		sc := C.ffi_prep_cif(&cif.c, C.FFI_DEFAULT_ABI, C.uint(n), cif.ret, cif.args)
		if sc != C.FFI_OK {
			return nil, fmt.Errorf("error while preparing cif (%s)", Status(sc))
		}
		return &cif, nil

	default:
		return nil, fmt.Errorf("ffi: invalid type kind %v", rt.Kind())
	}

	return nil, fmt.Errorf("ffi: invalid type kind %v", rt.Kind())
}

func (cif *CIF) Release() {
	C.free(unsafe.Pointer(cif.args))
	cif.args = nil
}

func (cif *CIF) Call(fptr unsafe.Pointer, ret unsafe.Pointer, args ...unsafe.Pointer) {
	var cargs *unsafe.Pointer

	if len(args) > 0 {
		cargs = (*unsafe.Pointer)(unsafe.Pointer(C.malloc(C.size_t(len(args)) * C.size_t(ptrSize))))
		for i, arg := range args {
			C._go_ffi_void_array(cargs, C.int(i), arg)
		}
		defer C.free(unsafe.Pointer(cargs))
	}

	C.ffi_call(&cif.c, C._go_ffi_func(fptr), ret, cargs)
}

func validateReflectType(typ reflect.Type) error {
	switch typ.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return nil
	case reflect.Float32, reflect.Float64:
		return nil
	}
	return fmt.Errorf("ffi: invalid type kind %v", typ.Kind())
}

func ffiTypeFrom(typ reflect.Type) Type {
	switch typ.Kind() {
	case reflect.Int:
		return Int
	case reflect.Int8:
		return Int8
	case reflect.Int16:
		return Int16
	case reflect.Int32:
		return Int32
	case reflect.Int64:
		return Int64
	case reflect.Uint:
		return Uint
	case reflect.Uint8:
		return Uint8
	case reflect.Uint16:
		return Uint16
	case reflect.Uint32:
		return Uint32
	case reflect.Uint64:
		return Uint64
	case reflect.Float32:
		return Float32
	case reflect.Float64:
		return Float64
	default:
		panic(fmt.Errorf("ffi: invalid type kind %v", typ.Kind()))
	}
}
