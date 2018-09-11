package rapidengine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"rapidengine/camera"
	"rapidengine/configuration"
)

type Child3D struct {
	vertexArray *VertexArray
	numVertices int32

	primitive string

	shaderProgram uint32

	textureEnabled bool
	texture        *uint32

	colorEnabled bool
	color        []float32

	modelMatrix      mgl32.Mat4
	projectionMatrix mgl32.Mat4

	copies         []ChildCopy
	currentCopies  []ChildCopy
	copyingEnabled bool

	X float32
	Y float32
	Z float32

	VX float32
	VY float32
	VZ float32

	Gravity float32

	Group    string
	collider Collider

	config           *configuration.EngineConfig
	collisionControl *CollisionControl
}

func NewChild3D(config *configuration.EngineConfig, collision *CollisionControl) Child3D {
	return Child3D{
		modelMatrix: mgl32.Ident4(),
		projectionMatrix: mgl32.Perspective(
			mgl32.DegToRad(45),
			float32(config.ScreenWidth)/float32(config.ScreenHeight),
			0.1, 100,
		),
		config:           config,
		Gravity:          0,
		textureEnabled:   false,
		copyingEnabled:   false,
		collisionControl: collision,
	}
}

func (child3D *Child3D) PreRender(mainCamera camera.Camera) {
	child3D.BindChild()

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child3D.modelMatrix[0],
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("projectionMtx\x00")),
		1, false, &child3D.projectionMatrix[0],
	)

	gl.BindAttribLocation(child3D.shaderProgram, 0, gl.Str("position\x00"))

	if child3D.textureEnabled {
		gl.BindAttribLocation(child3D.shaderProgram, 1, gl.Str("tex\x00"))
	}

	gl.BindVertexArray(0)
}

func (child3D *Child3D) BindChild() {
	gl.BindVertexArray(child3D.vertexArray.id)
	gl.UseProgram(child3D.shaderProgram)
}

func (child3D *Child3D) Update(mainCamera camera.Camera, delta float64, lastFrame float64) {
	child3D.VY -= child3D.Gravity

	child3D.X += child3D.VX
	child3D.Y += child3D.VY
	child3D.X += child3D.VZ

	child3D.Render(mainCamera)
}

func (child3D *Child3D) Render(mainCamera camera.Camera) {
	child3D.modelMatrix = mgl32.Translate3D(child3D.X, child3D.Y, child3D.Z)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child3D.modelMatrix[0],
	)

	if child3D.colorEnabled {
		gl.Uniform3fv(gl.GetUniformLocation(child3D.shaderProgram, gl.Str("color\x00")), 1, &child3D.color[0])
	}
}

func (child3D *Child3D) AttachTexture(coords []float32, texture *uint32) {
	if child3D.vertexArray == nil {
		panic("Cannot attach texture without VertexArray")
	}
	if child3D.shaderProgram == 0 {
		panic("Cannot attach texture without shader program")
	}

	gl.BindVertexArray(child3D.vertexArray.id)
	gl.UseProgram(child3D.shaderProgram)

	child3D.vertexArray.AddVertexAttribute(coords, 1, 2)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *texture)

	loc1 := gl.GetUniformLocation(child3D.shaderProgram, gl.Str("texture0\x00"))
	gl.Uniform1i(loc1, int32(0))
	CheckError("IGNORE")

	child3D.textureEnabled = true
	child3D.texture = texture
	gl.BindVertexArray(0)
}

func (child3D *Child3D) AttachColor(rgb []float32) {
	child3D.color = rgb
}

func (child3D *Child3D) EnableColor() {
	child3D.colorEnabled = true
}

func (child3D *Child3D) SetPosition(x, y, z float32) {
	child3D.X = x
	child3D.Y = y
	child3D.Z = z
}

func (child3D *Child3D) AttachVertexArray(vao *VertexArray, numVertices int32) {
	child3D.vertexArray = vao
	child3D.numVertices = numVertices
}

func (child3D *Child3D) AttachPrimitive(p Primitive) {
	child3D.AttachVertexArray(p.vao, p.numVertices)
	child3D.vertexArray.AddVertexAttribute(CubeNormals, 2, 3)
}

func (child3D *Child3D) AttachShader(s uint32) {
	child3D.shaderProgram = s
}

func (child3D *Child3D) GetX() float32 {
	return child3D.X
}

func (child3D *Child3D) GetY() float32 {
	return child3D.Y
}

func (child3D *Child3D) GetZ() float32 {
	return child3D.Z
}

func (child3D *Child3D) GetShaderProgram() uint32 {
	return child3D.shaderProgram
}

func (child3D *Child3D) GetVertexArray() *VertexArray {
	return child3D.vertexArray
}

func (child3D *Child3D) GetNumVertices() int32 {
	return child3D.numVertices
}

func (child3D *Child3D) GetTexture() *uint32 {
	return child3D.texture
}

func (child3D *Child3D) GetTextureEnabled() bool {
	return child3D.textureEnabled
}

func (child3D *Child3D) GetCollider() *Collider {
	return nil
}

func (child3D *Child3D) RenderCopy(c ChildCopy, cam camera.Camera)       {}
func (child3D *Child3D) AddCurrentCopy(c ChildCopy)                      {}
func (child3D *Child3D) GetCurrentCopies() []ChildCopy                   { return nil }
func (child3D *Child3D) GetCopies() []ChildCopy                          { return nil }
func (child3D *Child3D) RemoveCurrentCopies()                            {}
func (child3D *Child3D) CheckCopyingEnabled() bool                       { return false }
func (child3D *Child3D) CheckCollision(c Child) int                      { return 0 }
func (child3D *Child3D) CheckCollisionRaw(x, y float32, c *Collider) int { return 0 }
