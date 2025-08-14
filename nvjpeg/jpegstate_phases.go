//go:build nvjpeg_phases
// +build nvjpeg_phases

package nvjpeg

/*
#cgo CFLAGS: -I/opt/cuda/include -I/opt/cuda/targets/x86_64-linux/include
#cgo LDFLAGS:-L/opt/cuda/lib64 -L/opt/cuda/targets/x86_64-linux/lib -lnvjpeg -lcudart -lcuda
#include <nvjpeg.h>
#include <cuda_runtime_api.h>
*/
import "C"
import (
    "github.com/negativeOne1/gocudnn/gocu"
)

// DecodePhase1 - CPU processing
func (j *JpegState) DecodePhase1(h *Handle, data []byte, frmt OutputFormat, s gocu.Streamer) error {
	d := (*C.uchar)(&data[0])
	length := len(data)
	return status(C.nvjpegDecodePhaseOne(h.h, j.j, d, C.size_t(length), frmt.c(), stream(s))).error()
}

// DecodePhase2 - Mixed processing
func (j *JpegState) DecodePhase2(h *Handle, s gocu.Streamer) error {
	return status(C.nvjpegDecodePhaseTwo(h.h, j.j, stream(s))).error()
}

// DecodePhase3 - GPU processing
func (j *JpegState) DecodePhase3(h *Handle, dest *Image, s gocu.Streamer) error {
	return status(C.nvjpegDecodePhaseThree(h.h, j.j, dest.cptr(), stream(s))).error()
}

// DecodeBatchedInitialize - init batch decoder
func (j *JpegState) DecodeBatchedInitialize(h *Handle, batchsize, maxCPUthreads int, frmt OutputFormat) error {
	return status(C.nvjpegDecodeBatchedInitialize(h.h, j.j, C.int(batchsize), C.int(maxCPUthreads), frmt.c())).error()
}

// DecodeBatched - decode batch
func (j *JpegState) DecodeBatched(h *Handle, data [][]byte, dest []*Image, s gocu.Streamer) error {
	x := make([]*C.uchar, len(data))
	y := make([]C.size_t, len(data))
	z := make([]*C.nvjpegImage_t, len(dest))
	var length int
	for i := range data {
		length = len(data[i])
		x[i] = (*C.uchar)(&data[i][0])
		y[i] = C.size_t(length)
		z[i] = dest[i].cptr()
	}

	return status(C.nvjpegDecodeBatched(h.h, j.j, &x[0], &y[0], z[0], stream(s))).error()
}

// DecodeBatchedPhase1 - CPU step
func (j *JpegState) DecodeBatchedPhase1(h *Handle, data []byte, imageidx, threadidx int, s gocu.Streamer) error {
	return status(C.nvjpegDecodeBatchedPhaseOne(h.h, j.j, (*C.uchar)(&data[0]), C.size_t(len(data)), C.int(imageidx), C.int(threadidx), stream(s))).error()
}

// DecodeBatchedPhase2 - Mixed step
func (j *JpegState) DecodeBatchedPhase2(h *Handle, s gocu.Streamer) error {
	return status(C.nvjpegDecodeBatchedPhaseTwo(h.h, j.j, stream(s))).error()
}

// DecodeBatchedPhase3 - GPU step
func (j *JpegState) DecodeBatchedPhase3(h *Handle, dest []*Image, s gocu.Streamer) error {
	z := make([]*C.nvjpegImage_t, len(dest))
	for i := range z {
		z[i] = dest[i].cptr()
	}
	return status(C.nvjpegDecodeBatchedPhaseThree(h.h, j.j, z[0], stream(s))).error()
}
