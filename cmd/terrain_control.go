package cmd

import (
	"fmt"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/material"
	"rapidengine/terrain"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type TerrainControl struct {
}

func NewTerrainControl() TerrainControl {
	return TerrainControl{}
}

func (terrainControl *TerrainControl) NewSkyBox(
	path string,

	shaderControl *ShaderControl,
	textureControl *TextureControl,

	config *configuration.EngineConfig,
) *terrain.SkyBox {

	gl.UseProgram(shaderControl.GetShader("skybox"))
	gl.BindAttribLocation(shaderControl.GetShader("skybox"), 0, gl.Str("position\x00"))

	textureControl.NewCubeMap(
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_LF.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_RT.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_UP.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_DN.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_FR.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_BK.png", path, path),
		"skybox",
	)

	material := material.NewMaterial(shaderControl.GetShader("skybox"), config)
	material.BecomeCubemap(textureControl.GetTexture("skybox"))

	indices := []uint32{}
	for i := 0; i < len(terrain.SkyBoxVertices); i++ {
		indices = append(indices, uint32(i))
	}

	vao := geometry.NewVertexArray(terrain.SkyBoxVertices, indices)
	vao.AddVertexAttribute(geometry.CubeTextures, 1, 2)

	return terrain.NewSkyBox(
		shaderControl.GetShader("skybox"),
		material,
		vao,
		mgl32.Perspective(
			mgl32.DegToRad(45),
			float32(config.ScreenWidth)/float32(config.ScreenHeight),
			0.1, 100,
		),
		mgl32.Ident4(),
	)
}
