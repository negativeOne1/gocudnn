package xtrakerns

import (
    "github.com/negativeOne1/gocudnn/cuda"
    "runtime"
    "testing"
)

func TestCreateModule(t *testing.T) {
	runtime.LockOSThread()
	devs, err := cuda.GetDeviceList()
	if err != nil {
		t.Error(err)
	}

	ctx, err := cuda.CtxCreate(-1, devs[0])
	if err != nil {
		t.Error(err)
	}

	ctx.Set()
	got, err := CreateModule(SwapEveryOther(), devs[0])
	if err != nil {
		t.Error(err)
	}
	print(got)
}
