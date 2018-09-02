package rapidengine

// --------------------------------------------------
// Child2D.go contains Child2D, the basic Object in
// rapidengine. Every game object is either a Child,
// or a copy of a Child.
// --------------------------------------------------

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"rapidengine/camera"
	"rapidengine/configuration"
)

type Child2D struct {
	vertexArray *VertexArray
	numVertices int32

	primitive string

	shaderProgram uint32
	texture       *uint32

	modelMatrix      mgl32.Mat4
	projectionMatrix mgl32.Mat4

	copies         []ChildCopy
	currentCopies  []ChildCopy
	copyingEnabled bool

	X float32
	Y float32

	VX float32
	VY float32

	Gravity float32

	Group    string
	collider Collider

	config           *configuration.EngineConfig
	collisioncontrol *CollisionControl
}

func NewChild2D(config *configuration.EngineConfig, collision *CollisionControl) Child2D {
	return Child2D{
		modelMatrix:      mgl32.Ident4(),
		projectionMatrix: mgl32.Ortho2D(-1, 1, -1, 1),
		config:           config,
		VX:               0,
		VY:               0,
		Gravity:          0,
		copyingEnabled:   false,
		collisioncontrol: collision,
	}
}

func (child2D *Child2D) PreRender(mainCamera camera.Camera) {
	child2D.config.Logger.Info("PreRendering Children...")
	child2D.BindChild()

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.modelMatrix[0],
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("projectionMtx\x00")),
		1, false, &child2D.projectionMatrix[0],
	)

	gl.BindAttribLocation(child2D.shaderProgram, 0, gl.Str("position\x00"))
	gl.BindAttribLocation(child2D.shaderProgram, 1, gl.Str("tex\x00"))

	gl.BindVertexArray(0)
}

func (child2D *Child2D) BindChild() {
	gl.BindVertexArray(child2D.vertexArray.id)
	gl.UseProgram(child2D.shaderProgram)
}

func (child2D *Child2D) Update(mainCamera camera.Camera) {
	cx, cy := mainCamera.GetPosition()
	if !child2D.collisioncontrol.CheckCollisionWithGroup(child2D, "ground", cx, cy) {
		child2D.VY -= child2D.Gravity
		child2D.X += child2D.VX
		child2D.Y += child2D.VY
	}

	child2D.Render(mainCamera)
}

func (child2D *Child2D) Render(mainCamera camera.Camera) {
	sX, sY := ScaleCoordinates(child2D.X, child2D.Y, float32(child2D.config.ScreenWidth), float32(child2D.config.ScreenHeight))
	child2D.modelMatrix = mgl32.Translate3D(sX, sY, 0)
	child2D.projectionMatrix = mgl32.Ortho2D(-1, 1, -1, 1)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.modelMatrix[0],
	)
}

func (child2D *Child2D) RenderCopy(config ChildCopy, mainCamera camera.Camera) {
	sX, sY := ScaleCoordinates(config.X, config.Y, float32(child2D.config.ScreenWidth), float32(child2D.config.ScreenHeight))
	child2D.modelMatrix = mgl32.Translate3D(sX, sY, 0)
	child2D.projectionMatrix = mgl32.Ortho2D(-1, 1, -1, 1)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.modelMatrix[0],
	)

	child2D.texture = config.Tex
}

func (child2D *Child2D) CheckCollision(other Child) bool {
	return child2D.collider.CheckCollision(child2D.X, child2D.Y, other.GetX(), other.GetY(), other.GetCollider())
}

func (child2D *Child2D) CheckCollisionRaw(otherX, otherY float32, otherCollider *Collider) bool {
	return child2D.collider.CheckCollision(child2D.X, child2D.Y, otherX, otherY, otherCollider)
}

//  --------------------------------------------------
//  Component Attachers
//  --------------------------------------------------

func (child2D *Child2D) AttachTexture(coords []float32, texture *uint32) error {
	if child2D.vertexArray == nil {
		return errors.New("Cannot attach texture without VertexArray")
	}
	if child2D.shaderProgram == 0 {
		return errors.New("Cannot attach texture without shader program")
	}

	gl.BindVertexArray(child2D.vertexArray.id)
	gl.UseProgram(child2D.shaderProgram)

	child2D.vertexArray.AddVertexAttribute(coords, 1, 2)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *texture)

	loc1 := gl.GetUniformLocation(child2D.shaderProgram, gl.Str("texture0\x00"))
	gl.Uniform1i(loc1, int32(0))

	child2D.texture = texture
	gl.BindVertexArray(0)
	return nil
}

func (child2D *Child2D) AttachTexturePrimitive(texture *uint32) {
	child2D.AttachTexture(GetPrimitiveCoords(child2D.primitive), texture)
}

func (child2D *Child2D) AttachCollider(x, y, w, h float32) {
	child2D.collider = NewCollider(x, y, w, h)
}

func (child2D *Child2D) AttachVertexArray(vao *VertexArray, numVertices int32) {
	child2D.vertexArray = vao
	child2D.numVertices = numVertices
}

func (child2D *Child2D) AttachPrimitive(p Primitive) {
	child2D.primitive = p.id
	child2D.AttachVertexArray(p.vao, p.numVertices)
}

func (child2D *Child2D) AttachShader(s uint32) {
	child2D.shaderProgram = s
}

func (child2D *Child2D) AttachGroup(group string) {
	child2D.Group = group
}

//  --------------------------------------------------
//  Setters
//  --------------------------------------------------

func (child2D *Child2D) SetVelocity(vx, vy float32) {
	child2D.VX = vx
	child2D.VY = vy
}

func (child2D *Child2D) SetVelocityX(vx float32) {
	child2D.VX = vx
}

func (child2D *Child2D) SetVelocityY(vy float32) {
	child2D.VY = vy
}

func (child2D *Child2D) SetPosition(x, y float32) {
	child2D.X = x
	child2D.Y = y
}

func (child2D *Child2D) SetX(x float32) {
	child2D.X = x
}

func (child2D *Child2D) SetY(y float32) {
	child2D.Y = y
}

func (child2D *Child2D) SetGravity(g float32) {
	child2D.Gravity = g
}

//  --------------------------------------------------
//  Getters
//  --------------------------------------------------

func (child2D *Child2D) GetShaderProgram() uint32 {
	return child2D.shaderProgram
}

func (child2D *Child2D) GetVertexArray() *VertexArray {
	return child2D.vertexArray
}

func (child2D *Child2D) GetNumVertices() int32 {
	return child2D.numVertices
}

func (child2D *Child2D) GetTexture() *uint32 {
	return child2D.texture
}

func (child2D *Child2D) GetCollider() *Collider {
	return &child2D.collider
}

func (child2D *Child2D) GetX() float32 {
	return child2D.X
}

func (child2D *Child2D) GetY() float32 {
	return child2D.Y
}

//  --------------------------------------------------
//  Copying
//  --------------------------------------------------

func (child2D *Child2D) EnableCopying() {
	child2D.copyingEnabled = true
}

func (child2D *Child2D) DisableCopying() {
	child2D.copyingEnabled = false
}

func (child2D *Child2D) AddCopy(config ChildCopy) {
	child2D.copies = append(child2D.copies, config)
}

func (child2D *Child2D) GetCopies() []ChildCopy {
	return child2D.copies
}

func (child2D *Child2D) CheckCopyingEnabled() bool {
	return child2D.copyingEnabled
}

func (child2D *Child2D) AddCurrentCopy(c ChildCopy) {
	child2D.currentCopies = append(child2D.currentCopies, c)
}

func (child2D *Child2D) RemoveCurrentCopies() {
	child2D.currentCopies = []ChildCopy{}
}

func (child2D *Child2D) GetCurrentCopies() []ChildCopy {
	return child2D.currentCopies
}

func ScaleCoordinates(x, y, sw, sh float32) (float32, float32) {
	return 2*(x/float32(sw)) - 1, 2*(y/float32(sh)) - 1
}
