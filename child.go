package rapidengine

import "rapidengine/camera"

type ChildCopy struct {
	X   float32
	Y   float32
	Tex *uint32
}
type Child interface {
	PreRender(camera.Camera)

	RenderCopy(ChildCopy, camera.Camera)

	CheckCopyingEnabled() bool

	GetShaderProgram() uint32

	GetVertexArray() *VertexArray

	GetNumVertices() int32

	GetTexture() *uint32

	GetCollider() *Collider

	GetCopies() []ChildCopy

	GetX() float32
	GetY() float32

	CheckCollision(Child) bool

	Update(camera.Camera)
}
