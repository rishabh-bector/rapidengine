package rapidengine

type Child interface {
	PreRender(Camera)

	GetShaderProgram() uint32

	GetVertexArray() *VertexArray

	GetNumVertices() int32

	GetTexture() uint32

	GetCollider() *Collider

	GetX() float32
	GetY() float32

	CheckCollision(Child) bool

	Update(Camera)
}
