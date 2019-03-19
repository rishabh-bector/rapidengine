package cmd

import (
	"encoding/json"
	"io/ioutil"
	"rapidengine/configuration"
	"rapidengine/material"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type TextureControl struct {
	TexMap map[string]*material.Texture `json:"textures"`

	config *configuration.EngineConfig
}

func NewTextureControl(config *configuration.EngineConfig) TextureControl {
	return TextureControl{
		TexMap: make(map[string]*material.Texture),
		config: config,
	}
}

func (textureControl *TextureControl) GetTexture(name string) *material.Texture {
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
	textureControl.TexMap[name] = &material.Texture{
		Name:   name,
		Path:   path,
		Filter: filter,
		Addr:   &texture,
	}
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

	textureControl.TexMap[name] = &material.Texture{
		Name: right,
		Path: right,
		Addr: &cubeMap,
	}
}

//   --------------------------------------------------
//   Disk
//   --------------------------------------------------

func (textureControl *TextureControl) Save(path string) {
	blob, err := json.Marshal(&textureControl.TexMap)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, blob, 0644)
	if err != nil {
		panic(err)
	}
}

func (textureControl *TextureControl) Load(path string) {
	blob, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	temp := TextureControl{
		TexMap: make(map[string]*material.Texture),
	}

	err = json.Unmarshal(blob, &temp.TexMap)
	if err != nil {
		panic(err)
	}

	println("Loading textures from file...")
	for n, t := range temp.TexMap {
		textureControl.NewTexture(t.Path, n, t.Filter)
	}
}
