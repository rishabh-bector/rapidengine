package control

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

type TextureControl struct {
	TexMap map[string]*uint32
}

func NewTextureControl() TextureControl {
	return TextureControl{
		make(map[string]*uint32),
	}
}

func (textureControl *TextureControl) GetTexture(name string) *uint32 {
	return textureControl.TexMap[name]
}

func (textureControl *TextureControl) NewTexture(path string, name string) error {
	rgba, err := loadImage(path)
	if err != nil {
		return err
	}

	var texture uint32

	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	textureControl.TexMap[name] = &texture
	return nil
}

func (textureControl *TextureControl) NewCubeMap(right, left, top, bottom, front, back, name string) {
	var cubeMap uint32

	gl.GenTextures(1, &cubeMap)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeMap)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	paths := []string{right, left, top, bottom, front, back}

	for i, path := range paths {
		rgba, err := loadImage(path)
		if err != nil {
			panic(err)
		}

		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0, gl.RGBA,
			int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	}

	textureControl.TexMap[name] = &cubeMap
}

func loadImage(path string) (*image.RGBA, error) {
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
