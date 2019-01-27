package geometry

//  --------------------------------------------------
//	Mesh.go allows easy creation of basic shapes
//  by automatically creating and binding a VAO to the
//  Mesh struct, which can be passed to a Child.
//  --------------------------------------------------

import (
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {

	// Mesh type
	id string

	// VAO containing vertices & indices
	vao *VertexArray

	// Texture coordinates
	texCoords []float32

	// Face normals
	normals []float32

	// Number of vertices
	numVertices int32

	// Normal Mapping
	tangents   []float32
	bitangents []float32

	// Material
	Material material.Material
}

func (p *Mesh) Render(viewMtx *float32, modelMtx *float32) {
	gl.UniformMatrix4fv(
		p.Material.GetShader().GetUniform("viewMtx"),
		1, false, viewMtx,
	)

	gl.UniformMatrix4fv(
		p.Material.GetShader().GetUniform("modelMtx"),
		1, false, modelMtx,
	)

}

func (p *Mesh) GetID() string {
	return p.id
}

func (p *Mesh) GetVAO() *VertexArray {
	return p.vao
}

func (p *Mesh) GetNumVertices() int32 {
	return p.numVertices
}

func (p *Mesh) GetTexCoords() []float32 {
	return p.texCoords
}

func (p *Mesh) GetNormals() []float32 {
	return p.normals
}

// NormalizeSizes takes in a size in pixels and normalizes to [0, 1]
func NormalizeSizes(x, y, sw, sh float32) (float32, float32) {
	return x / sw, y / sh
}

func normalizeX(x, sw float32) float32 {
	return x / sw
}

func normalizeY(y, sh float32) float32 {
	return y / sh
}

// GetMeshCoords returns the appropriate texture coordinates
// for each mesh.
func GetMeshCoords(id string) []float32 {
	switch id {
	case "rectangle":
		return RectTextures
	}
	return RectTextures
}

// ComputeTangents calculates the tangents and bitangents
// of a mesh, based on vertices and UVs.
func (m *Mesh) ComputeTangents() {
	for i := 0; i < len(m.vao.vertices); i++ {
		m.tangents = append(m.tangents, 0)
		m.bitangents = append(m.bitangents, 0)
	}

	for i := 0; i < int(len(m.vao.indices)); i += 3 {
		v0 := mgl32.Vec4{
			m.vao.vertices[m.vao.indices[i]*3],
			m.vao.vertices[m.vao.indices[i]*3+1],
			m.vao.vertices[m.vao.indices[i]*3+2],
		}
		v1 := mgl32.Vec4{
			m.vao.vertices[m.vao.indices[i+1]*3],
			m.vao.vertices[m.vao.indices[i+1]*3+1],
			m.vao.vertices[m.vao.indices[i+1]*3+2],
		}
		v2 := mgl32.Vec4{
			m.vao.vertices[m.vao.indices[i+2]*3],
			m.vao.vertices[m.vao.indices[i+2]*3+1],
			m.vao.vertices[m.vao.indices[i+2]*3+2],
		}

		texCoords := m.texCoords
		uv0 := mgl32.Vec2{
			texCoords[m.vao.indices[i]*3],
			texCoords[m.vao.indices[i]*3+1],
		}
		uv1 := mgl32.Vec2{
			texCoords[m.vao.indices[i+1]*3],
			texCoords[m.vao.indices[i+1]*3+1],
		}
		uv2 := mgl32.Vec2{
			texCoords[m.vao.indices[i+2]*3],
			texCoords[m.vao.indices[i+2]*3+1],
		}

		e1 := v1.Sub(v0)
		e2 := v2.Sub(v0)

		deltaUV1 := uv1.Sub(uv0)
		deltaUV2 := uv2.Sub(uv0)

		r := float32(1) / (deltaUV1.X()*deltaUV2.Y() - deltaUV1.Y()*deltaUV2.X())

		tangent := (e1.Mul(deltaUV2.Y()).Sub(e2.Mul(deltaUV1.Y()))).Mul(r)
		bitangent := (e2.Mul(deltaUV1.X()).Sub(e1.Mul(deltaUV2.X()))).Mul(r)

		m.tangents[m.vao.indices[i]*3] += tangent.X()
		m.tangents[m.vao.indices[i]*3+1] += tangent.Y()
		m.tangents[m.vao.indices[i]*3+2] += tangent.Z()

		m.tangents[m.vao.indices[i+1]*3] += tangent.X()
		m.tangents[m.vao.indices[i+1]*3+1] += tangent.Y()
		m.tangents[m.vao.indices[i+1]*3+2] += tangent.Z()

		m.tangents[m.vao.indices[i+2]*3] += tangent.X()
		m.tangents[m.vao.indices[i+2]*3+1] += tangent.Y()
		m.tangents[m.vao.indices[i+2]*3+2] += tangent.Z()

		m.bitangents[m.vao.indices[i]*3] += bitangent.X()
		m.bitangents[m.vao.indices[i]*3+1] += bitangent.Y()
		m.bitangents[m.vao.indices[i]*3+2] += bitangent.Z()

		m.bitangents[m.vao.indices[i+1]*3] += bitangent.X()
		m.bitangents[m.vao.indices[i+1]*3+1] += bitangent.Y()
		m.bitangents[m.vao.indices[i+1]*3+2] += bitangent.Z()

		m.bitangents[m.vao.indices[i+2]*3] += bitangent.X()
		m.bitangents[m.vao.indices[i+2]*3+1] += bitangent.Y()
		m.bitangents[m.vao.indices[i+2]*3+2] += bitangent.Z()
	}

	m.GetVAO().AddVertexAttribute(m.tangents, 3, 3)
	m.GetVAO().AddVertexAttribute(m.bitangents, 4, 3)

}
