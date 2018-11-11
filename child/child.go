package child

import "rapidengine/camera"
import "rapidengine/material"
import "rapidengine/geometry"
import "rapidengine/physics"

type ChildCopy struct {
	X        float32
	Y        float32
	Z        float32
	Material *material.Material
	Darkness float32

	ID string
}
type Child interface {
	PreRender(camera.Camera)

	RenderCopy(ChildCopy, camera.Camera)
	CheckCopyingEnabled() bool

	GetShaderProgram() uint32
	GetVertexArray() *geometry.VertexArray
	GetNumVertices() int32

	GetCollider() *physics.Collider
	GetCopies() *[]ChildCopy
	GetNumCopies() int

	GetX() float32
	GetY() float32

	AddCurrentCopy(ChildCopy)
	GetCurrentCopies() []ChildCopy
	RemoveCurrentCopies()

	SetSpecificRenderDistance(float32)
	GetSpecificRenderDistance() float32

	CheckCollision(Child) int
	CheckCollisionRaw(otherX float32, otherY float32, otherCollider *physics.Collider) int

	MouseCollisionFunc(bool)

	Update(camera.Camera, float64, float64)
}
