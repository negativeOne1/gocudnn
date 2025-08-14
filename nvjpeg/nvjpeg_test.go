package nvjpeg

import (
    "bytes"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "testing"
)

func TestImageInfo(t *testing.T) {
	check := func(err error, t *testing.T) {
		if err != nil {
			t.Error(err)
		}
	}
    var backendflag Backend
    handle, err := CreateEx(backendflag.Default())
    check(err, t)

    // Generate a simple 2x2 JPEG in-memory
    img := image.NewRGBA(image.Rect(0, 0, 2, 2))
    img.Set(0, 0, color.RGBA{0, 0, 0, 255})
    img.Set(1, 0, color.RGBA{255, 0, 0, 255})
    img.Set(0, 1, color.RGBA{0, 255, 0, 255})
    img.Set(1, 1, color.RGBA{0, 0, 255, 255})

    var buf bytes.Buffer
    err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
    check(err, t)
    imgbytes := buf.Bytes()

    subsampletype, ws, hs, err := GetImageInfo(handle, imgbytes)
	check(err, t)
	fmt.Println("Subsample Type: ", subsampletype.String())
	fmt.Println("Widths :", ws)
	fmt.Println("Heights :", hs)
	//decoder, err := nvjpeg.JpegStateCreate(handle)
	//	check(err, t)

}
