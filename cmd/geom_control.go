package cmd

import (
	"rapidengine/geometry"
	"rapidengine/material"

	"github.com/krux02/assimp"
)

type GeometryControl struct {
	engine *Engine
}

func NewGeometryControl() GeometryControl {
	return GeometryControl{}
}

func (gm *GeometryControl) Initialize(e *Engine) {
	gm.engine = e
}

func (gm *GeometryControl) LoadModel(path string, mat material.Material) geometry.Model {
	scene := assimp.ImportFile(path, uint(assimp.Process_Triangulate|assimp.Process_FlipUVs))
	model := geometry.Model{
		Materials: make(map[int]material.Material),
	}

	model.Materials[0] = mat

	// Recursively process all nodes in the scene
	gm.processNode(&model, scene.RootNode(), scene)

	return model
}

func (gm *GeometryControl) processNode(m *geometry.Model, node *assimp.Node, scene *assimp.Scene) {
	for _, mesh := range node.Meshes() {
		m.Meshes = append(m.Meshes, gm.createMesh(m, scene.Meshes()[mesh], scene))
	}
	for _, childNode := range node.Children() {
		gm.processNode(m, childNode, scene)
	}
}

func (gm *GeometryControl) createMesh(m *geometry.Model, mesh *assimp.Mesh, scene *assimp.Scene) geometry.Mesh {
	vertices, uvs, normals, tangents, bitangents := gm.loadMeshData(mesh)
	indices := gm.loadMeshIndices(mesh)

	ms := geometry.Mesh{
		ID:         "assimp",
		VAO:        geometry.NewVertexArray(vertices, indices),
		TexCoords:  uvs,
		Normals:    normals,
		Tangents:   tangents,
		Bitangents: bitangents,

		ModelMaterial: 0,
	}

	ms.NumVertices = int32(len(ms.VAO.GetIndices()))

	if len(ms.TexCoords) > 0 {
		ms.VAO.AddVertexAttribute(ms.TexCoords, 1, 3)
		ms.TexCoordsEnabled = true
	}

	if len(ms.Normals) > 0 {
		ms.VAO.AddVertexAttribute(ms.Normals, 2, 3)
		ms.NormalsEnabled = true
	}

	if len(ms.Tangents) > 0 {
		ms.VAO.AddVertexAttribute(ms.Tangents, 3, 3)
		ms.TangentsEnabled = true
	}

	if len(ms.Bitangents) > 0 {
		ms.VAO.AddVertexAttribute(ms.Bitangents, 4, 3)
		ms.BitangentsEnabled = true
	}

	ms.ModelMaterial = gm.loadMeshMaterial(m, mesh, scene)

	return ms
}

func (gm *GeometryControl) loadMeshData(mesh *assimp.Mesh) ([]float32, []float32, []float32, []float32, []float32) {
	a_vertices := mesh.Vertices()
	a_uvs := mesh.TextureCoords(0)
	a_normals := mesh.Normals()
	a_tangents := mesh.Tangents()
	a_bitangents := mesh.Bitangents()

	vertices := []float32{}
	uvs := []float32{}
	normals := []float32{}
	tangents := []float32{}
	bitangents := []float32{}

	for i := 0; i < mesh.NumVertices(); i++ {
		vertices = append(vertices, a_vertices[i].X(), a_vertices[i].Y(), a_vertices[i].Z())

		if a_uvs != nil {
			uvs = append(uvs, a_uvs[i].X(), a_uvs[i].Y(), a_uvs[i].Z())
		}

		if a_normals != nil {
			normals = append(normals, a_normals[i].X(), a_normals[i].Y(), a_normals[i].Z())
		}

		if a_tangents != nil {
			tangents = append(tangents, a_tangents[i].X(), a_tangents[i].Y(), a_tangents[i].Z())
		}

		if a_bitangents != nil {
			bitangents = append(bitangents, a_bitangents[i].X(), a_bitangents[i].Y(), a_bitangents[i].Z())
		}
	}

	return vertices, uvs, normals, tangents, bitangents
}

func (gm *GeometryControl) loadMeshIndices(mesh *assimp.Mesh) []uint32 {
	indices := []uint32{}

	for i := 0; i < mesh.NumFaces(); i++ {
		indices = append(indices, mesh.Faces()[i].CopyIndices()...)
	}

	return indices
}

func (gm *GeometryControl) loadMeshMaterial(model *geometry.Model, mesh *assimp.Mesh, scene *assimp.Scene) int {
	if mesh.MaterialIndex() >= 0 {
		material := scene.Materials()[mesh.MaterialIndex()]

		diffuse := gm.loadMaterialTexture(material, assimp.TextureMapping_Diffuse)
		//specular := gm.loadMaterialTexture(material, assimp.TextureMapping_Diffuse)
		//normal := gm.loadMaterialTexture(material, assimp.TextureMapping_Diffuse)

		if _, ok := model.Materials[mesh.MaterialIndex()]; !ok {
			newMat := gm.engine.MaterialControl.NewStandardMaterial()

			if diffuse != "" {
				newMat.DiffuseLevel = 1
				newMat.AttachDiffuseMap(gm.engine.TextureControl.GetTexture(diffuse))
			} else {
				newMat.DiffuseLevel = 0
				newMat.Hue = [4]float32{0, 0, 100, 255}
			}

			model.Materials[mesh.MaterialIndex()] = newMat
		}

		return mesh.MaterialIndex()
	}

	return 0
}

func (gm *GeometryControl) loadMaterialTexture(mat *assimp.Material, tm assimp.TextureMapping) string {
	textureType := assimp.TextureType(tm)
	path, _, _, _, _, _, _, _ := mat.GetMaterialTexture(textureType, 0)

	if path != "" {
		gm.engine.TextureControl.NewTexture(path, path, "mipmap")
	}

	return path
}
