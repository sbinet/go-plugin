#include "ffi.h"

typedef void (*_go_ffi_func)(void);

void _go_ffi_type_array(ffi_type **array, int index, ffi_type *value);
void _go_ffi_void_array(void **array, int index, void *value);
