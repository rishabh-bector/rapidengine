package rapidengine

type Child interface {
	PreRender(Camera)

	GetShaderProgram() uint32

	GetVertexArray() *VertexArray

	GetNumVertices() int32

	GetTexture() uint32

	Update(Camera)
}
