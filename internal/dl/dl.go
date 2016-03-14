package dl

// #include <stdlib.h>
// #include <dlfcn.h>
import "C"

import (
	"fmt"
	"unsafe"
)

type Flags int

const (
	Lazy     Flags = C.RTLD_LAZY
	Now      Flags = C.RTLD_NOW
	Global   Flags = C.RTLD_GLOBAL
	Local    Flags = C.RTLD_LOCAL
	NoLoad   Flags = C.RTLD_NOLOAD
	NoDelete Flags = C.RTLD_NODELETE
	// First Flags = C.RTLD_FIRST
)

type Handle struct {
	Addr unsafe.Pointer
}

func Open(fname string, flags Flags) (Handle, error) {
	c_str := C.CString(fname)
	defer C.free(unsafe.Pointer(c_str))

	h := C.dlopen(c_str, C.int(flags))
	if h == nil {
		c_err := C.dlerror()
		return Handle{}, fmt.Errorf("dl: %s", C.GoString(c_err))
	}
	return Handle{h}, nil
}

func (h Handle) Close() error {
	o := C.dlclose(h.Addr)
	if o != C.int(0) {
		c_err := C.dlerror()
		return fmt.Errorf("dl: %s", C.GoString(c_err))
	}
	return nil
}

func (h Handle) Symbol(symbol string) (unsafe.Pointer, error) {
	sym := C.CString(symbol)
	defer C.free(unsafe.Pointer(sym))

	addr := C.dlsym(h.Addr, sym)
	if addr == nil {
		err := C.dlerror()
		return nil, fmt.Errorf("dl: %s", C.GoString(err))
	}
	return addr, nil
}

// /* Portable libltdl versions of the system dlopen() API. */
// LT_SCOPE lt_dlhandle lt_dlopen          (const char *filename);
// LT_SCOPE lt_dlhandle lt_dlopenext       (const char *filename);
// LT_SCOPE lt_dlhandle lt_dlopenadvise    (const char *filename,
//                                          lt_dladvise advise);
// LT_SCOPE void *     lt_dlsym            (lt_dlhandle handle, const char *name);
// LT_SCOPE const char *lt_dlerror         (void);
// LT_SCOPE int        lt_dlclose          (lt_dlhandle handle);

// EOF
