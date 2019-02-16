package geometry

import (
	"rapidengine/material"
)

// A model can be imported from a 3D object file
// format such as OBJ or STL, and can contain multiple
// meshes / shapes.

type Model struct {
	Meshes    []Mesh
	Materials map[int]material.Material
}

func (m *Model) Render(viewMtx *float32, modelMtx *float32, projMtx *float32) {
	for _, ms := range m.Meshes {
		ms.Render(m.Materials[ms.ModelMaterial], viewMtx, modelMtx, projMtx, 0, 0, 1)
	}
}

func (m *Model) ComputeTangents() {
	for _, ms := range m.Meshes {
		ms.ComputeTangents()
	}
}
