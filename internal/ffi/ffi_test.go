package ffi_test

import (
	"path"
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/sbinet/go-plugin/internal/dl"
	"github.com/sbinet/go-plugin/internal/ffi"
)

func eq(t *testing.T, ref, chk interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = path.Base(file)
	if !reflect.DeepEqual(ref, chk) {
		t.Errorf("%s:%d: expected [%v], got [%v]", file, line, ref, chk)
	}
}

func TestCIF(t *testing.T) {
	cif, err := ffi.New(ffi.Float64, []ffi.Type{ffi.Float64})
	if err != nil {
		t.Fatal(err)
	}

	lib, err := dl.Open(libmName, dl.Now)
	if err != nil {
		t.Fatal(err)
	}
	defer lib.Close()

	fptr, err := lib.Symbol("fabs")
	if err != nil {
		t.Fatal(err)
	}

	var (
		cval = -42.0
		cret = 0.0
	)

	cif.Call(fptr, unsafe.Pointer(&cret), unsafe.Pointer(&cval))
	if cret != 42 {
		t.Fatalf("error calling fabs(-42): got=%v. want=42\n", cret)
	}
}

/*
func TestFFIMathf(t *testing.T) {
	t.Skip("not compatible with new CGo rules")

	lib, err := ffi.NewLibrary(libm_name)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer lib.Close()

	tests := []struct {
		fct  string
		arg  float64
		want float64
	}{
		{"cos", 0., math.Cos(0.)},
		{"cos", math.Pi / 2., math.Cos(math.Pi / 2.)},
		{"sin", 0., math.Sin(0.)},
		{"sin", math.Pi / 2., math.Sin(math.Pi / 2.)},
	}

	for _, info := range tests {
		f, err := lib.Fct(info.fct, ffi.C_double, []ffi.Type{ffi.C_double})
		if err != nil {
			t.Errorf("could not locate function [%s]: %v", info.fct, err)
		}
		out := f(info.arg).Float()
		if math.Abs(out-info.want) > 1e-16 {
			t.Errorf("expected [%v], got [%v] (fct=%v(%v))", info.want, out, info.fct, info.arg)
		}

	}

	err = lib.Close()
	if err != nil {
		t.Errorf("error closing [%s]: %v", libm_name, err)
	}
}

func TestFFIMathi(t *testing.T) {
	t.Skip("not compatible with new CGo rules")

	lib, err := ffi.NewLibrary(libm_name)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer lib.Close()

	f, err := lib.Fct("abs", ffi.C_int, []ffi.Type{ffi.C_int})
	if err != nil {
		t.Errorf("could not locate function [abs]: %v", err)
	}
	{
		out := f(10).Int()
		if out != 10 {
			t.Errorf("expected [10], got [%v] (fct=abs(10))", out)
		}

	}
	{
		out := f(-10).Int()
		if out != 10 {
			t.Errorf("expected [10], got [%v] (fct=abs(-10))", out)
		}

	}

	err = lib.Close()
	if err != nil {
		t.Errorf("error closing [%s]: %v", libm_name, err)
	}
}

func TestFFIStrCmp(t *testing.T) {
	t.Skip("not compatible with new CGo rules")

	lib, err := ffi.NewLibrary(libc_name)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer lib.Close()

	//int strcmp(const char* cs, const char* ct);
	f, err := lib.Fct("strcmp", ffi.C_int, []ffi.Type{ffi.C_pointer, ffi.C_pointer})
	if err != nil {
		t.Errorf("could not locate function [strcmp]: %v", err)
	}
	{
		s1 := "foo"
		s2 := "foo"
		out := f(s1, s2).Int()
		if out != 0 {
			t.Errorf("expected [0], got [%v]", out)
		}

	}
	{
		s1 := "foo"
		s2 := "foo1"
		out := f(s1, s2).Int()
		if out == 0 {
			t.Errorf("expected [!0], got [%v]", out)
		}

	}

	err = lib.Close()
	if err != nil {
		t.Errorf("error closing [%s]: %v", libc_name, err)
	}
}

func TestFFIStrLen(t *testing.T) {
	t.Skip("not compatible with new CGo rules")

	lib, err := ffi.NewLibrary(libc_name)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer lib.Close()

	//size_t strlen(const char* cs);
	f, err := lib.Fct("strlen", ffi.C_int, []ffi.Type{ffi.C_pointer})
	if err != nil {
		t.Errorf("could not locate function [strlen]: %v", err)
	}
	{
		str := `foo-bar-\nfoo foo`
		out := int(f(str).Int())
		if out != len(str) {
			t.Errorf("expected [%d], got [%d]", len(str), out)
		}

	}

	err = lib.Close()
	if err != nil {
		t.Errorf("error closing [%s]: %v", libc_name, err)
	}
}

func TestFFIStrCat(t *testing.T) {
	t.Skip("not compatible with new CGo rules")

	lib, err := ffi.NewLibrary(libc_name)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer lib.Close()

	//char* strcat(char* s, const char* ct);
	f, err := lib.Fct("strcat", ffi.C_pointer, []ffi.Type{ffi.C_pointer, ffi.C_pointer})
	if err != nil {
		t.Errorf("could not locate function [strlen]: %v", err)
	}
	{
		s1 := "foo"
		s2 := "bar"
		out := f(s1, s2).String()
		//FIXME
		if out != "foobar" && false {
			t.Errorf("expected [foobar], got [%s] (s1=%s, s2=%s)", out, s1, s2)
		}

	}

	err = lib.Close()
	if err != nil {
		t.Errorf("error closing [%s]: %v", libc_name, err)
	}
}
*/
