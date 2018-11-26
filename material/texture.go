package material

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

func LoadImage(path string) (*image.RGBA, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	ct, err := detectContentType(imgFile)
	if err != nil {
		return nil, err
	}

	img, err := convert(imgFile, ct)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func convert(file io.Reader, contentType string) (*image.RGBA, error) {
	var decode func(io.Reader) (image.Image, error)
	switch contentType {
	case "image/jpeg":
		decode = jpeg.Decode
	case "image/png":
		decode = png.Decode
	default:
		return nil, fmt.Errorf("unrecognized image format: %s", contentType)
	}
	src, err := decode(file)
	if err != nil {
		return nil, err
	}

	img := &notOpaqueRGBA{image.NewRGBA(src.Bounds())}
	draw.Draw(img, img.Bounds(), src, image.ZP, draw.Src)

	return img.RGBA, nil
}

type notOpaqueRGBA struct {
	*image.RGBA
}

func (i *notOpaqueRGBA) Opaque() bool {
	return false
}

func detectContentType(file io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return "", err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf), nil
}

func LoadImageFullDepth(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	src, _, err := image.Decode(imgFile)

	if err != nil {
		return nil, err
	}

	return src, nil
}
