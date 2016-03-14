package plugin

//#include <string.h>
import "C"

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sbinet/go-plugin/internal/dl"
)

var (
	errNilPtr = errors.New("plugin: nil pointer")
	errNotPtr = errors.New("plugin: expected a pointer to a value")
)

// A plugin opened at runtime.
type Plugin struct {
	handle dl.Handle
}

// Open the plugin name.
func Open(name string) (Plugin, error) {
	lib, err := dl.Open(name, dl.Now)
	if err != nil {
		return Plugin{}, err
	}
	return Plugin{lib}, nil
}

// Lookup looks up a symbol in a Go style plugin.
// You can only look up functions and global variables.
func (p Plugin) Lookup(name string) (interface{}, error) {
	panic("not implemented")
}

// LookupC looks up a symbol in a C style plugin, passing in a pointer to
// a value with the type it is expected to have.
// The value must be a function type with a C style API, or a C variable type.
func (p Plugin) LookupC(name string, valptr interface{}) error {
	addr, err := p.handle.Symbol(name)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(valptr)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr {
		return errNotPtr
	}
	if rv.IsNil() {
		return errNilPtr
	}

	val := rv.Elem()
	switch val.Kind() {
	case reflect.Int:
		val.SetInt(int64(*(*int)(addr)))
	case reflect.Int8:
		val.SetInt(int64(*(*int8)(addr)))
	case reflect.Int16:
		val.SetInt(int64(*(*int16)(addr)))
	case reflect.Int32:
		val.SetInt(int64(*(*int32)(addr)))
	case reflect.Int64:
		val.SetInt(int64(*(*int64)(addr)))
	case reflect.Uint:
		val.SetUint(uint64(*(*uint)(addr)))
	case reflect.Uint8:
		val.SetUint(uint64(*(*uint8)(addr)))
	case reflect.Uint16:
		val.SetUint(uint64(*(*uint16)(addr)))
	case reflect.Uint32:
		val.SetUint(uint64(*(*uint32)(addr)))
	case reflect.Uint64:
		val.SetUint(uint64(*(*uint64)(addr)))
	case reflect.Uintptr:
		val.SetUint(uint64(*(*uintptr)(addr)))
	case reflect.Float32:
		val.SetFloat(float64(*(*float32)(addr)))
	case reflect.Float64:
		val.SetFloat(float64(*(*float64)(addr)))
	case reflect.Complex64:
		// FIXME(sbinet) C layout of complex may differ from Go's
		val.SetComplex(complex128(*(*complex64)(addr)))
	case reflect.Complex128:
		// FIXME(sbinet) C layout of complex may differ from Go's
		val.SetComplex(complex128(*(*complex128)(addr)))
	case reflect.String:
		val.SetString(C.GoString(*(**C.char)(addr)))
	case reflect.UnsafePointer:
		val.SetPointer(addr)

	default:
		return fmt.Errorf("plugin: invalid type %T", valptr)
	}

	return err
}

// Close the plugin.
func (p Plugin) Close() error {
	return p.handle.Close()
}
