package plugin_test

import (
	"testing"
	"unsafe"

	"github.com/sbinet/go-plugin"
)

func TestOpen(t *testing.T) {
	for _, test := range []string{libcName, libmName} {
		p, err := plugin.Open(test)
		if err != nil {
			t.Error(err)
			continue
		}
		defer p.Close()

		err = p.Close()
		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestLookupC(t *testing.T) {
	for _, test := range []struct {
		lib string
		sym string
	}{
		{
			lib: libcName,
			sym: "puts",
		},
		{
			lib: libmName,
			sym: "fabs",
		},
	} {
		p, err := plugin.Open(test.lib)
		if err != nil {
			t.Error(err)
			continue
		}
		defer p.Close()

		var val unsafe.Pointer

		err = p.LookupC(test.sym, &val)
		if err != nil {
			t.Error(err)
			continue
		}

		err = p.Close()
		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestLibmFabs(t *testing.T) {
	p, err := plugin.Open(libmName)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	{
		var fabs func(float64) float64
		err = p.LookupC("fabs", &fabs)
		if err != nil {
			t.Fatal(err)
		}

		v := fabs(-42)
		if v != 42 {
			t.Fatalf("fabs(-42)=%v\n", v)
		}
	}

	{
		var fabs func(float32) float32
		err = p.LookupC("fabsf", &fabs)
		if err != nil {
			t.Fatal(err)
		}

		v := fabs(-42)
		if v != 42 {
			t.Fatalf("fabsf(-42)=%v\n", v)
		}
	}

	err = p.Close()
	if err != nil {
		t.Fatal(err)
	}
}
