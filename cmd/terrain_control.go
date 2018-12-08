package cmd

import (
	"fmt"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/terrain"

	"github.com/go-gl/mathgl/mgl32"
)

type TerrainControl struct {
	engine *Engine

	terrainEnabled bool

	root *terrain.Terrain

	foliage *terrain.Foliage
}

func NewTerrainControl() TerrainControl {
	return TerrainControl{}
}

func (tc *TerrainControl) Initialize(engine *Engine) {
	tc.engine = engine
}

func (tc *TerrainControl) Update() {
	if tc.terrainEnabled {
		tc.engine.Renderer.RenderChild(tc.root.TChild)
		tc.engine.Renderer.RenderChild(tc.foliage.FChild)
	}
}

func (tc *TerrainControl) NewTerrain(width int, height int, vertices int) *terrain.Terrain {
	t := terrain.NewTerrain(width, height)

	t.TChild = tc.engine.ChildControl.NewChild3D()

	t.TChild.AttachMaterial(tc.engine.MaterialControl.NewTerrainMaterial())
	t.TChild.AttachMesh(geometry.NewPlane(width, height, vertices, nil, 1))

	t.TChild.PreRender(tc.engine.Renderer.MainCamera)

	tc.terrainEnabled = true
	tc.root = &t

	return &t
}

func (tc *TerrainControl) NewFoliage(width int, height int) *terrain.Foliage {
	f := terrain.NewFoliage(width, height)

	f.FChild = tc.engine.ChildControl.NewChild3D()

	f.FChild.AttachMaterial(tc.engine.MaterialControl.NewFoliageMaterial())
	f.FChild.AttachMesh(geometry.LoadObj("./billboard.obj"))

	f.FChild.PreRender(tc.engine.Renderer.MainCamera)

	tc.foliage = &f
	return &f
}

func (terrainControl *TerrainControl) NewSkyBox(
	path string,

	shaderControl *ShaderControl,
	textureControl *TextureControl,

	config *configuration.EngineConfig,
) *terrain.SkyBox {

	shaderControl.GetShader("skybox").Bind()

	textureControl.NewCubeMap(
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_LF.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_RT.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_UP.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_DN.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_FR.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_BK.png", path, path),
		"skybox",
	)

	material := terrainControl.engine.MaterialControl.NewCubemapMaterial()
	material.CubeDiffuseMap = textureControl.GetTexture("skybox")

	indices := []uint32{}
	for i := 0; i < len(terrain.SkyBoxVertices); i++ {
		indices = append(indices, uint32(i))
	}

	vao := geometry.NewVertexArray(terrain.SkyBoxVertices, indices)
	vao.AddVertexAttribute(geometry.CubeTextures, 1, 2)

	return terrain.NewSkyBox(
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
