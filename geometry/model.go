package geometry

import (
	"github.com/krux02/assimp"
)

// A model can be imported from a 3D object file
// format such as OBJ or STL, and can contain multiple
// meshes / shapes.

type Model struct {
	Meshes   []Mesh
	textures map[string]*uint32
}

func LoadModel(path string) {
	scene := assimp.ImportFile(path, uint(assimp.Process_Triangulate|assimp.Process_FlipUVs))
	model := Model{}

	// Recursively process all nodes in the scene
	model.processNode(scene.RootNode(), scene)
}

func (m *Model) processNode(node *assimp.Node, scene *assimp.Scene) {
	for _, mesh := range node.Meshes() {
		m.Meshes = append(m.Meshes, createMesh(scene.Meshes()[mesh]))
	}
	for _, childNode := range node.Children() {
		m.processNode(childNode, scene)
	}
}

func createMesh(mesh *assimp.Mesh) Mesh {
	vertices, uvs, normals, tangents, bitangents := processMeshData(mesh)
	indices := processMeshIndices(mesh)

	ms := Mesh{
		id:         "assimp",
		vao:        NewVertexArray(vertices, indices),
		texCoords:  uvs,
		normals:    normals,
		tangents:   tangents,
		bitangents: bitangents,
	}

	ms.numVertices = int32(len(ms.vao.indices))

	ms.vao.AddVertexAttribute(ms.texCoords, 1, 3)
	ms.vao.AddVertexAttribute(ms.normals, 2, 3)

	return ms
}

func processMeshData(mesh *assimp.Mesh) ([]float32, []float32, []float32, []float32, []float32) {
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

func processMeshIndices(mesh *assimp.Mesh) []uint32 {
	indices := []uint32{}

	for i := 0; i < mesh.NumFaces(); i++ {
		indices = append(indices, mesh.Faces()[i].CopyIndices()...)
	}

	return indices
}

func processMeshTextures(mesh *assimp.Mesh) {

}
