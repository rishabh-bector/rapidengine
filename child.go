package rapidengine

type ChildCopy struct {
	X float32
	Y float32
}
type Child interface {
	PreRender(Camera)

	RenderCopy(ChildCopy, Camera)

	CheckCopyingEnabled() bool

	GetShaderProgram() uint32

	GetVertexArray() *VertexArray

	GetNumVertices() int32

	GetTexture() uint32

	GetCollider() *Collider

	GetCopies() []ChildCopy

	GetX() float32
	GetY() float32

	CheckCollision(Child) bool

	Update(Camera)
}
