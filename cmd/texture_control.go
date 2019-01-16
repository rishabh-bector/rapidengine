package cmd

import (
	"rapidengine/configuration"
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type TextureControl struct {
	TexMap map[string]*uint32

	config *configuration.EngineConfig
}

func NewTextureControl(config *configuration.EngineConfig) TextureControl {
	return TextureControl{
		TexMap: make(map[string]*uint32),
		config: config,
	}
}

func (textureControl *TextureControl) GetTexture(name string) *uint32 {
	if tx, ok := textureControl.TexMap[name]; ok {
		return tx
	} else {
		panic("couldn't find texture: " + name)
	}
}

func (textureControl *TextureControl) NewTexture(path string, name string, filter string) error {
	rgba, err := material.LoadImage(path)
	if err != nil {
		panic(err)
	}

	var texture uint32

	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	switch filter {

	case "pixel":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	case "mipmap":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	case "linear":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	case "anisotropic":
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	}

	texFormat := gl.RGBA
	if textureControl.config.GammaCorrection {
		texFormat = gl.SRGB_ALPHA
	}

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		int32(texFormat),
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
		rgba, err := material.LoadImage(path)
		if err != nil {
			panic(err)
		}

		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0, gl.RGBA,
			int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	}

	textureControl.TexMap[name] = &cubeMap
}
