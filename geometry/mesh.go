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
	ID string

	// VAO containing vertices & indices
	VAO *VertexArray

	// Texture coordinates
	TexCoords        []float32
	TexCoordsEnabled bool

	// Face normals
	Normals        []float32
	NormalsEnabled bool

	// Number of vertices
	NumVertices int32

	// Normal Mapping
	Tangents        []float32
	TangentsEnabled bool

	Bitangents        []float32
	BitangentsEnabled bool

	// Material
	ModelMaterial int
}

func (p *Mesh) Render(mat material.Material, viewMtx *float32, modelMtx *float32, projMtx *float32) {
	gl.BindVertexArray(p.VAO.id)
	mat.GetShader().Bind()

	gl.EnableVertexAttribArray(0)

	if p.TexCoordsEnabled {
		gl.EnableVertexAttribArray(1)
	}

	if p.NormalsEnabled {
		gl.EnableVertexAttribArray(2)
	}

	if p.TangentsEnabled {
		gl.EnableVertexAttribArray(3)
	}

	if p.BitangentsEnabled {
		gl.EnableVertexAttribArray(4)
	}

	gl.UniformMatrix4fv(
		mat.GetShader().GetUniform("viewMtx"),
		1, false, viewMtx,
	)

	gl.UniformMatrix4fv(
		mat.GetShader().GetUniform("modelMtx"),
		1, false, modelMtx,
	)

	gl.UniformMatrix4fv(
		mat.GetShader().GetUniform("projectionMtx"),
		1, false, projMtx,
	)

	mat.Render(0, 1, 0)

	gl.DrawElements(gl.TRIANGLES, p.NumVertices, gl.UNSIGNED_INT, gl.PtrOffset(0))
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

// ComputeTangents calculates the Tangents and Bitangents
// of a mesh, based on vertices and UVs.
func (m *Mesh) ComputeTangents() {
	for i := 0; i < len(m.VAO.vertices); i++ {
		m.Tangents = append(m.Tangents, 0)
		m.Bitangents = append(m.Bitangents, 0)
	}

	for i := 0; i < int(len(m.VAO.indices)); i += 3 {
		v0 := mgl32.Vec4{
			m.VAO.vertices[m.VAO.indices[i]*3],
			m.VAO.vertices[m.VAO.indices[i]*3+1],
			m.VAO.vertices[m.VAO.indices[i]*3+2],
		}
		v1 := mgl32.Vec4{
			m.VAO.vertices[m.VAO.indices[i+1]*3],
			m.VAO.vertices[m.VAO.indices[i+1]*3+1],
			m.VAO.vertices[m.VAO.indices[i+1]*3+2],
		}
		v2 := mgl32.Vec4{
			m.VAO.vertices[m.VAO.indices[i+2]*3],
			m.VAO.vertices[m.VAO.indices[i+2]*3+1],
			m.VAO.vertices[m.VAO.indices[i+2]*3+2],
		}

		texCoords := m.TexCoords
		uv0 := mgl32.Vec2{
			texCoords[m.VAO.indices[i]*3],
			texCoords[m.VAO.indices[i]*3+1],
		}
		uv1 := mgl32.Vec2{
			texCoords[m.VAO.indices[i+1]*3],
			texCoords[m.VAO.indices[i+1]*3+1],
		}
		uv2 := mgl32.Vec2{
			texCoords[m.VAO.indices[i+2]*3],
			texCoords[m.VAO.indices[i+2]*3+1],
		}

		e1 := v1.Sub(v0)
		e2 := v2.Sub(v0)

		deltaUV1 := uv1.Sub(uv0)
		deltaUV2 := uv2.Sub(uv0)

		r := float32(1) / (deltaUV1.X()*deltaUV2.Y() - deltaUV1.Y()*deltaUV2.X())

		tangent := (e1.Mul(deltaUV2.Y()).Sub(e2.Mul(deltaUV1.Y()))).Mul(r)
		bitangent := (e2.Mul(deltaUV1.X()).Sub(e1.Mul(deltaUV2.X()))).Mul(r)

		m.Tangents[m.VAO.indices[i]*3] += tangent.X()
		m.Tangents[m.VAO.indices[i]*3+1] += tangent.Y()
		m.Tangents[m.VAO.indices[i]*3+2] += tangent.Z()

		m.Tangents[m.VAO.indices[i+1]*3] += tangent.X()
		m.Tangents[m.VAO.indices[i+1]*3+1] += tangent.Y()
		m.Tangents[m.VAO.indices[i+1]*3+2] += tangent.Z()

		m.Tangents[m.VAO.indices[i+2]*3] += tangent.X()
		m.Tangents[m.VAO.indices[i+2]*3+1] += tangent.Y()
		m.Tangents[m.VAO.indices[i+2]*3+2] += tangent.Z()

		m.Bitangents[m.VAO.indices[i]*3] += bitangent.X()
		m.Bitangents[m.VAO.indices[i]*3+1] += bitangent.Y()
		m.Bitangents[m.VAO.indices[i]*3+2] += bitangent.Z()

		m.Bitangents[m.VAO.indices[i+1]*3] += bitangent.X()
		m.Bitangents[m.VAO.indices[i+1]*3+1] += bitangent.Y()
		m.Bitangents[m.VAO.indices[i+1]*3+2] += bitangent.Z()

		m.Bitangents[m.VAO.indices[i+2]*3] += bitangent.X()
		m.Bitangents[m.VAO.indices[i+2]*3+1] += bitangent.Y()
		m.Bitangents[m.VAO.indices[i+2]*3+2] += bitangent.Z()
	}

	m.VAO.AddVertexAttribute(m.Tangents, 3, 3)
	m.VAO.AddVertexAttribute(m.Bitangents, 4, 3)
	m.TangentsEnabled = true
	m.BitangentsEnabled = true
}
