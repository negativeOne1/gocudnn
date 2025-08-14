//Package crtutil allows cudart to work with Go's io Reader and Writer interfaces.
//
//This package only works with devices that have Compute Capability 6.1 and up.
package crtutil

import (
	"errors"
	"io"
	"unsafe"

	"github.com/dereklstinson/cutil"
	"github.com/negativeOne1/gocudnn/cudart"
	"github.com/negativeOne1/gocudnn/gocu"
)

//ReadWriter is made to work with the golang io packages
type ReadWriter struct {
	p      unsafe.Pointer
	i      uint
	size   uint
	cpyflg cudart.MemcpyKind
	s      gocu.Streamer
}

//Allocator alocates memory to the current device
type Allocator struct {
	s gocu.Streamer
	d gocu.Device
	w *gocu.Worker
}

//CreateAllocator creates an allocator whose memory it creates does async mem copies.
//This creates its own worker for device passed.
func CreateAllocator(s gocu.Streamer, d gocu.Device) (a *Allocator) {
	a = new(Allocator)
	a.s = s
	a.d = d
	a.w = gocu.NewWorker(d)
	return a
}

//AllocateMemory allocates memory on the current device.  Allocater allocates memory for device that was passed when
//created onto the devices default context.
func (a *Allocator) AllocateMemory(size uint) (r *ReadWriter, err error) {

	err = a.w.Work(func() error {
		r = new(ReadWriter)
		r.size = size
		err = cudart.MallocManagedGlobal(r, size)
		r.s = a.s
		r.cpyflg.Default()
		return err
	})

	return r, err
}

//NewReadWriter returns ReadWriter from already allocated memory passed in p.  It just needs to know the size of the memory.
//If s is nil. Then it will do a non async copy.  If it is not nil then it will do a async copy.
func NewReadWriter(p cutil.Pointer, size uint, s gocu.Streamer) *ReadWriter {
	r := &ReadWriter{
		p:    p.Ptr(),
		size: size,
		s:    s,
	}
	r.cpyflg.Default()
	return r
}

//Ptr satisfies cutil.Pointer interface.
func (r *ReadWriter) Ptr() unsafe.Pointer {
	return r.p
}

//DPtr satisfies cutil.DPointer interface.
func (r *ReadWriter) DPtr() *unsafe.Pointer {
	return &r.p
}

//Reset resets the index to 0.
func (r *ReadWriter) Reset() {
	r.i = 0
}

//Len returns the remaining bytes that are not read.
func (r *ReadWriter) Len() int {
	if r.i >= r.size {
		return 0
	}
	return (int)(r.size) - (int)(r.i)
}

//Size returns the total size in bytes of the memory the readwriter holds
func (r *ReadWriter) Size() uint {
	return r.size
}

//Read satisfies the io.Reader interface
func (r *ReadWriter) Read(b []byte) (n int, err error) {
	if r.i >= r.size {
		r.Reset()
		return 0, io.EOF
	}
	if len(b) == 0 {
		return 0, nil
	}
	var size = r.size - r.i
	if uint(len(b)) < size {
		size = uint(len(b))
	}
	bwrap, err := cutil.WrapGoMem(b)
	if err != nil {
		return 0, err
	}
	if r.s != nil {
		err = cudart.MemcpyAsync(bwrap, cutil.Offset(r, r.i), size, r.cpyflg, r.s)
	} else {
		err = cudart.Memcpy(bwrap, cutil.Offset(r, r.i), size, r.cpyflg)
	}
	if err != nil {
		return 0, nil
	}
	r.i += size
	n = int(size)
	if len(b) == int(r.size) {
		r.Reset()
		return n, io.EOF
	}
	return n, nil

}

//Write satisfies the io.Writer interface
func (r *ReadWriter) Write(b []byte) (n int, err error) {
	if r.i >= r.size {
		r.Reset()
		return 0, errors.New("(r *ReadWriter) Write()" +
			"Write Location Out of Memory")
	}
	if len(b) == 0 {
		return 0, nil
	}
	var size = r.size - r.i
	if uint(len(b)) < size {
		size = uint(len(b))
	}
	bwrap, err := cutil.WrapGoMem(b)
	if err != nil {
		return 0, err
	}
	if r.s != nil {
		err = cudart.MemcpyAsync(cutil.Offset(r, r.i), bwrap, size, r.cpyflg, r.s)
	} else {
		err = cudart.Memcpy(cutil.Offset(r, r.i), bwrap, size, r.cpyflg)
	}
	r.i += size
	n = int(size)
	return n, err
}

/*
func (r *ReadWriter) ReadAt(b []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("crtutil.ReadWriter.ReadAt: negative offset")
	}
	if off >= int64(r.size) {
		return 0, io.EOF
	}
	cudart.MemCpy()
}
*/
