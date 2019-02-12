package cmd

import (
	"fmt"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/material"
	"rapidengine/terrain"

	"github.com/go-gl/mathgl/mgl32"
)

type TerrainControl struct {
	engine *Engine

	terrainEnabled bool
	root           *terrain.Terrain

	foliages []*terrain.Foliage

	waters []*terrain.Water
}

func NewTerrainControl() TerrainControl {
	return TerrainControl{}
}

func (tc *TerrainControl) Initialize(engine *Engine) {
	tc.engine = engine
}

func (tc *TerrainControl) Update() {
	if tc.terrainEnabled {
		tc.engine.Renderer.RenderTerrainChild(tc.root.TChild)

		for _, f := range tc.foliages {
			tc.engine.Renderer.RenderChild(f.FChild)
		}

		for _, w := range tc.waters {
			tc.engine.Renderer.RenderChild(w.WChild)
		}
	}
}

func (tc *TerrainControl) InstanceFoliage(f *terrain.Foliage) {
	tc.foliages = append(tc.foliages, f)
}

func (tc *TerrainControl) InstanceWater(w *terrain.Water) {
	tc.waters = append(tc.waters, w)
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

func (tc *TerrainControl) NewPlanetaryTerrain(width int, height int, vertices int) *terrain.Terrain {
	t := terrain.NewTerrain(width, height)

	t.TChild = tc.engine.ChildControl.NewChild3D()

	t.TChild.AttachMaterial(tc.engine.MaterialControl.NewTerrainMaterial())
	//t.TChild.AttachMesh(geometry.NewPlane(width, height, vertices, nil, 1))
	t.TChild.AttachMesh(geometry.LoadObj("../rapidengine/assets/obj/sphere.obj", 10000))
	t.TChild.SetInstanceRenderDistance(1000000000)

	t.TChild.PreRender(tc.engine.Renderer.MainCamera)

	tc.terrainEnabled = true
	tc.root = &t

	return &t
}

func (tc *TerrainControl) NewFoliage(width int, height int, instances int) *terrain.Foliage {
	f := terrain.NewFoliage(width, height)

	f.FChild = tc.engine.ChildControl.NewChild3D()

	f.FChild.AttachMaterial(tc.engine.MaterialControl.NewFoliageMaterial())
	f.FChild.AttachMesh(geometry.LoadObj("./billboard.obj", 1))

	f.FChild.EnableGLInstancing(instances)
	f.FChild.SetInstanceRenderDistance(100000)

	f.FChild.PreRender(tc.engine.Renderer.MainCamera)

	return &f
}

func (tc *TerrainControl) NewWater(width int, height int, vertices int) *terrain.Water {
	w := terrain.NewWater(width, height)

	w.WChild = tc.engine.ChildControl.NewChild3D()

	w.WChild.AttachMaterial(tc.engine.MaterialControl.NewWaterMaterial())
	w.WChild.AttachMesh(geometry.NewPlane(width, height, vertices, nil, 1))

	w.WChild.PreRender(tc.engine.Renderer.MainCamera)

	return &w
}

func (terrainControl *TerrainControl) NewSkyBox(
	path string,
	ext string,

	shaderControl *ShaderControl,
	textureControl *TextureControl,

	config *configuration.EngineConfig,
) *terrain.SkyBox {

	shaderControl.GetShader("skybox").Bind()

	textureControl.NewCubeMap(
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_LF.%s", path, path, ext),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_RT.%s", path, path, ext),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_UP.%s", path, path, ext),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_DN.%s", path, path, ext),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_FR.%s", path, path, ext),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_BK.%s", path, path, ext),
		"skybox",
	)

	cmaterial := terrainControl.engine.MaterialControl.NewCubemapMaterial()
	cmaterial.CubeDiffuseMap = textureControl.GetTexture("skybox")

	indices := []uint32{}
	for i := 0; i < len(terrain.SkyBoxVertices); i++ {
		indices = append(indices, uint32(i))
	}

	vao := geometry.NewVertexArray(terrain.SkyBoxVertices, indices)
	vao.AddVertexAttribute(geometry.CubeTextures, 1, 2)

	return terrain.NewSkyBox(
		cmaterial,
		vao,
		mgl32.Perspective(
			mgl32.DegToRad(45),
			float32(config.ScreenWidth)/float32(config.ScreenHeight),
			0.1, 100,
		),
		mgl32.Ident4(),
		[]*material.ShaderProgram{
			terrainControl.engine.ShaderControl.GetShader("standard"),
		},
	)
}
