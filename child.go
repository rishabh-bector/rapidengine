package main

type Child interface {
	PreRender(Camera)

	GetShaderProgram() uint32

	GetVertexArray() *VertexArray

	GetNumVertices() int32

	Update()
}
