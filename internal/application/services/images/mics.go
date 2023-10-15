package images

import (
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"math"
)

func getMagic(blob []byte, size int) ([]byte, error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	//mw.GetImage()
	err := mw.ReadImageBlob(blob)
	if err != nil {
		return nil, err
	}
	fmt.Println("Dimensions", mw.GetImageWidth(), mw.GetImageWidth())
	//mw.MergeImageLayers(imagick.IMAGE_LAYER_OPTIMIZE_IMAGE)

	////imagick.NewMagickImage()
	//ret, err := mw.ConvertImageCommand([]string{"convert", "logo:", "-resize", "100x100", "/tmp/out.png"})
	//if err != nil {return nil, err}
	//fmt.Println("Meta:", ret.Meta)
	result := resizeImage(mw, size, false)

	return result.GetImageBlob(), nil
}

func getDimensions(wand *imagick.MagickWand) (w, h uint) {
	h = wand.GetImageHeight()
	w = wand.GetImageWidth()
	return
}
func resizeRatio(im *imagick.MagickWand, width, height int) float64 {
	return math.Abs((float64)(width*height) / (float64)(im.GetImageWidth()*im.GetImageHeight()))
}
func heightToWidthRatio(im *imagick.MagickWand) float64 {
	return (float64)(im.GetImageHeight()) / (float64)(im.GetImageWidth())
}
func resizeImage(wand *imagick.MagickWand, size int, mozaic bool) *imagick.MagickWand {
	width := float32(wand.GetImageWidth())
	height := float32(wand.GetImageHeight())
	var rate float32
	if width > height {
		rate = float32(size) / width
	} else {
		rate = float32(size) / height
	}
	if mozaic {
		wand.ResizeImage(uint(width*rate/20), uint(height*rate/20), imagick.FILTER_LANCZOS, 1)
		wand.ResizeImage(uint(width*rate), uint(height*rate), imagick.FILTER_POINT, 1)
	} else {
		wand.ResizeImage(uint(width*rate), uint(height*rate), imagick.FILTER_LANCZOS, 1)
	}
	return wand.GetImage()
}

/*
func glitchImage(wand *imagick.MagickWand, q url.Values) ([]byte, error) {
	data := wand.GetImage().GetImageBlob()
	jpgHeaderLength := getJpegHeaderSize(data)
	maxIndex := len(data) - jpgHeaderLength - 4
	params := getParams(q)
	length := int(params["iterations"])
	for i := 0; i < length; i++ {
		pxMin := math.Floor(float64(maxIndex) / params["iterations"] * float64(i))
		pxMax := math.Floor(float64(maxIndex) / params["iterations"] * float64((i + 1)))
		delta := pxMax - pxMin
		pxI := math.Floor(pxMin + delta*params["seed"])
		if int(pxI) > maxIndex {
			pxI = float64(maxIndex)
		}
		index := math.Floor(float64(jpgHeaderLength) + pxI)
		data[int(index)] = byte(math.Floor(params["amount"] * float64(256)))
	}
	wand2 := imagick.NewMagickWand()
	if err := wand2.ReadImageBlob(data); err != nil {
		return nil, err
	}
	if err := wand2.SetImageFormat("PNG"); err != nil {
		return nil, err
	}
	return wand2.GetImage().GetImageBlob(), nil
}
*/

func crop(mw *imagick.MagickWand, x, y int, cols, rows uint) error {
	var result error
	result = nil
	imCols := mw.GetImageWidth()
	imRows := mw.GetImageHeight()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if uint(x) >= imCols || uint(y) >= imRows {
		result = fmt.Errorf("x, y more than image cols, rows")
		return result
	}
	if cols == 0 || imCols < uint(x)+cols {
		cols = imCols - uint(x)
	}
	if rows == 0 || imRows < uint(y)+rows {
		rows = imRows - uint(y)
	}
	fmt.Print(fmt.Sprintf("wi_crop(im, %d, %d, %d, %d)\n", x, y, cols, rows))
	result = mw.CropImage(cols, rows, x, y)
	return result
}

func shrinkImage(wand *imagick.MagickWand, maxSize int) (w, h uint) {
	w, h = getDimensions(wand)
	shrinkBy := 1
	if w >= h {
		shrinkBy = int(w) / maxSize
	} else {
		shrinkBy = int(h) / maxSize
	}
	wand.AdaptiveResizeImage(
		uint(int(w)/shrinkBy),
		uint(int(h)/shrinkBy),
	)
	// Sharpen the image to bring back some of the color lost in the shrinking
	//wand.AdaptiveSharpenImage(0, AdaptiveSharpenVal)
	w, h = getDimensions(wand)
	return
}

func rotate(blob []byte, degrees float64) ([]byte, error) {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	err := mw.ReadImageBlob(blob)
	if err != nil {
		return nil, err
	}
	err = mw.RotateImage(mw.GetBackgroundColor(), degrees)
	result := mw.GetImageBlob()
	return result, err
}

//info https://zalinux.ru/?p=7544
