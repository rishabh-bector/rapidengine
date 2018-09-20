package rapidengine

import "rapidengine/camera"

type ChildCopy struct {
	X        float32
	Y        float32
	Z        float32
	Material *Material
}
type Child interface {
	PreRender(camera.Camera)

	RenderCopy(ChildCopy, camera.Camera)
	CheckCopyingEnabled() bool

	GetShaderProgram() uint32
	GetVertexArray() *VertexArray
	GetNumVertices() int32

	GetCollider() *Collider
	GetCopies() []ChildCopy

	GetX() float32
	GetY() float32

	AddCurrentCopy(ChildCopy)
	GetCurrentCopies() []ChildCopy
	RemoveCurrentCopies()

	CheckCollision(Child) int
	CheckCollisionRaw(otherX float32, otherY float32, otherCollider *Collider) int

	Update(camera.Camera, float64, float64)
}
