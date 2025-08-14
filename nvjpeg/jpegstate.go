package nvjpeg

/*
#include <nvjpeg.h>
#include <cuda_runtime_api.h>
*/
import "C"
import (
    "runtime"

    "github.com/negativeOne1/gocudnn/gocu"
)

//JpegState is Opaque jpeg decoding state handle identifier - used to store intermediate information between deccding phases
type JpegState struct {
	j C.nvjpegJpegState_t
}

//CreateJpegState creates an initialized decode state
func CreateJpegState(h *Handle) (*JpegState, error) {
	j := new(JpegState)
	err := status(C.nvjpegJpegStateCreate(h.h, &j.j)).error()
	if err != nil {
		return nil, err
	}
	runtime.SetFinalizer(j, nvjpegJpegStateDestroy)
	return j, nil
}

func nvjpegJpegStateDestroy(j *JpegState) error {
	err := status(C.nvjpegJpegStateDestroy(j.j)).error()
	if err != nil {
		return err
	}
	j = nil
	return nil
}

/*Decode does the nvjpegDecode

  Decodes single image. Destination buffers should be large enough to be able to store
  output of specified format. For each color plane sizes could be retrieved for image using nvjpegGetImageInfo()
  and minimum required memory buffer for each plane is nPlaneHeight*nPlanePitch where nPlanePitch >= nPlaneWidth for
  planar output formats and nPlanePitch >= nPlaneWidth*nOutputComponents for interleaved output format.

  Function will perform an s.Sync() before returning.

  IN/OUT     h             : Library handle

  IN         data          : Pointer to the buffer containing the jpeg image to be decoded.

  IN         fmt           : Output data format. See nvjpegOutputFormat_t for description

  IN/OUT     dest	 	 : Pointer to structure with information about output buffers. See nvjpegImage_t description.

  IN/OUT     s             : gocu.Streamer where to submit all GPU work

*/
func (j *JpegState) Decode(h *Handle, data []byte, frmt OutputFormat, dest *Image, s gocu.Streamer) error {
	d := (*C.uchar)(&data[0])
	length := len(data)
	err := status(C.nvjpegDecode(h.h, j.j, d, C.size_t(length), frmt.c(), dest.cptr(), stream(s))).error()
	if err != nil {
		return err
	}
	return s.Sync()
}
func (j *JpegState) AttachPinnedBuffer(p *PinnedBuffer) error {
	return status(C.nvjpegStateAttachPinnedBuffer(j.j, p.b)).error()
}
func (j *JpegState) AttachDeviceBuffer(d *DeviceBuffer) error {
	return status(C.nvjpegStateAttachDeviceBuffer(j.j, d.b)).error()
}
