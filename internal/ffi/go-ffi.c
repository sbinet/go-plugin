#include "go-ffi.h"

void _go_ffi_type_array(ffi_type **args, int i, ffi_type *arg) {
	args[i] = arg;
}

void _go_ffi_void_array(void **args, int i, void *arg) {
	args[i] = arg;
}
