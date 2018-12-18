package child

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"rapidengine/camera"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/material"
	"rapidengine/physics"
)

type Child3D struct {
	active bool

	numVertices int32

	mesh geometry.Mesh

	Material material.Material

	modelMatrix      mgl32.Mat4
	projectionMatrix mgl32.Mat4

	copies         []ChildCopy
	currentCopies  []ChildCopy
	copyingEnabled bool

	instancingEnabled bool
	numInstances      int

	X float32
	Y float32
	Z float32

	VX float32
	VY float32
	VZ float32

	RX float32
	RY float32
	RZ float32

	ScaleX float32
	ScaleY float32
	ScaleZ float32

	Gravity float32

	Group    string
	collider physics.Collider

	specificRenderDistance float32

	config *configuration.EngineConfig
}

func NewChild3D(config *configuration.EngineConfig) *Child3D {
	return &Child3D{
		modelMatrix: mgl32.Ident4(),
		projectionMatrix: mgl32.Perspective(
			mgl32.DegToRad(45),
			float32(config.ScreenWidth)/float32(config.ScreenHeight),
			0.1, 100000,
		),
		config:                 config,
		Gravity:                0,
		copyingEnabled:         false,
		specificRenderDistance: 0,
		ScaleX:                 1,
		ScaleY:                 1,
		ScaleZ:                 1,
	}
}

func (child3D *Child3D) PreRender(mainCamera camera.Camera) {
	child3D.BindChild()

	gl.UniformMatrix4fv(
		child3D.Material.GetShader().GetUniform("modelMtx"),
		1, false, &child3D.modelMatrix[0],
	)

	gl.UniformMatrix4fv(
		child3D.Material.GetShader().GetUniform("viewMtx"),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		child3D.Material.GetShader().GetUniform("projectionMtx"),
		1, false, &child3D.projectionMatrix[0],
	)

	/*
		if child3D.copyingEnabled {
			vData := []float32{}
			for _, c := range child3D.GetCopies() {
				vData = append(vData, c.X, c.Y, c.Z)
			}
			child3D.vertexArray.AddVertexAttribute(vData, 3, 3)
			gl.VertexAttribDivisor(3, 1)
		}
		gl.BindAttribLocation(child3D.shaderProgram, 3, gl.Str("copyPosition\x00"))*/

	gl.BindVertexArray(0)
}

func (child3D *Child3D) BindChild() {
	gl.BindVertexArray(child3D.mesh.GetVAO().GetID())
	child3D.Material.GetShader().Bind()
}

func (child3D *Child3D) Update(mainCamera camera.Camera, delta float64, totalTime float64) {
	child3D.VY -= child3D.Gravity

	child3D.X += child3D.VX
	child3D.Y += child3D.VY
	child3D.X += child3D.VZ

	child3D.Render(mainCamera, totalTime)
}

func (child3D *Child3D) Render(mainCamera camera.Camera, totalTime float64) {
	child3D.modelMatrix = mgl32.Translate3D(child3D.X, child3D.Y, child3D.Z)
	child3D.modelMatrix = child3D.modelMatrix.Mul4(mgl32.Scale3D(child3D.ScaleX, child3D.ScaleY, child3D.ScaleZ))

	child3D.modelMatrix = child3D.modelMatrix.Mul4(mgl32.HomogRotate3DX(child3D.RX))
	child3D.modelMatrix = child3D.modelMatrix.Mul4(mgl32.HomogRotate3DY(child3D.RY))
	child3D.modelMatrix = child3D.modelMatrix.Mul4(mgl32.HomogRotate3DZ(child3D.RZ))

	gl.UniformMatrix4fv(
		child3D.Material.GetShader().GetUniform("viewMtx"),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		child3D.Material.GetShader().GetUniform("modelMtx"),
		1, false, &child3D.modelMatrix[0],
	)

	child3D.Material.Render(0, 1, totalTime)
}

func (child3D *Child3D) RenderCopy(config ChildCopy, mainCamera camera.Camera) {
	child3D.modelMatrix = mgl32.Translate3D(config.X, config.Y, config.Z)

	gl.UniformMatrix4fv(
		child3D.Material.GetShader().GetUniform("viewMtx"),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		child3D.Material.GetShader().GetUniform("modelMtx"),
		1, false, &child3D.modelMatrix[0],
	)

	c := []float32{1, 0, 0}
	gl.Uniform3fv(
		child3D.Material.GetShader().GetUniform("copyingEnabled"),
		1, &c[0],
	)

	config.Material.Render(0, 1, 0)
}

func (child3D *Child3D) AttachTextureCoords(coords []float32) {
	if child3D.mesh.GetVAO() == nil {
		panic("Cannot attach texture without VertexArray")
	}

	child3D.BindChild()
	child3D.mesh.GetVAO().AddVertexAttribute(coords, 1, 3)
	gl.BindVertexArray(0)
}

func (child3D *Child3D) SetPosition(x, y, z float32) {
	child3D.X = x
	child3D.Y = y
	child3D.Z = z
}

func (child3D *Child3D) AttachMaterial(m material.Material) {
	child3D.Material = m
}

func (child3D *Child3D) AttachMesh(p geometry.Mesh) {
	child3D.mesh = p
	child3D.numVertices = p.GetNumVertices()
	child3D.mesh.GetVAO().AddVertexAttribute(*p.GetNormals(), 2, 3)
	child3D.AttachTextureCoords(*p.GetTexCoords())
}

func (child3D *Child3D) ComputeMeshTangents() {
	child3D.mesh.ComputeTangents()
}

func (child3D *Child3D) Activate() {
	child3D.active = true
}

func (child3D *Child3D) Deactivate() {
	child3D.active = false
}

func (child3D *Child3D) IsActive() bool {
	return child3D.active
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

func (child3D *Child3D) GetShaderProgram() *material.ShaderProgram {
	return child3D.Material.GetShader()
}

func (child3D *Child3D) GetVertexArray() *geometry.VertexArray {
	return child3D.mesh.GetVAO()
}

func (child3D *Child3D) GetNumVertices() int32 {
	return child3D.numVertices
}

func (child3D *Child3D) GetCollider() *physics.Collider {
	return nil
}

//  --------------------------------------------------
//  Copying
//  --------------------------------------------------

func (child3D *Child3D) EnableCopying() {
	child3D.copyingEnabled = true
}

func (child3D *Child3D) DisableCopying() {
	child3D.copyingEnabled = false
}

func (child3D *Child3D) AddCopy(config ChildCopy) {
	child3D.copies = append(child3D.copies, config)
}

func (child3D *Child3D) GetCopies() *[]ChildCopy {
	return &child3D.copies
}

func (child3D *Child3D) GetNumCopies() int {
	return 0
}

func (child3D *Child3D) GetCurrentCopies() []ChildCopy {
	return child3D.currentCopies
}

func (child3D *Child3D) CheckCopyingEnabled() bool {
	return child3D.copyingEnabled
}

func (child3D *Child3D) AddCurrentCopy(c ChildCopy) {
	child3D.currentCopies = append(child3D.currentCopies, c)
}

func (child3D *Child3D) RemoveCurrentCopies() {
	child3D.currentCopies = []ChildCopy{}
}

func (child3D *Child3D) CheckCollision(other Child) int {
	return child3D.collider.CheckCollision(child3D.X, child3D.Y, child3D.VX, child3D.VY, other.GetX(), other.GetY(), other.GetCollider())
}

func (child3D *Child3D) CheckCollisionRaw(otherX, otherY float32, otherCollider *physics.Collider) int {
	return child3D.collider.CheckCollision(child3D.X, child3D.Y, child3D.VX, child3D.VY, otherX, otherY, otherCollider)
}

func (child3D *Child3D) SetSpecificRenderDistance(d float32) {
	child3D.specificRenderDistance = d
}

func (child3D *Child3D) GetSpecificRenderDistance() float32 {
	return child3D.specificRenderDistance
}

func (child3D *Child3D) MouseCollisionFunc(collision bool) {

}

//  --------------------------------------------------
//  GL Instancing
//  --------------------------------------------------

func (child3D *Child3D) CheckInstancingEnabled() bool {
	return child3D.instancingEnabled
}

func (child3D *Child3D) GetNumInstances() int {
	return child3D.numInstances
}

func (child3D *Child3D) EnableGLInstancing(num int) {
	child3D.instancingEnabled = true
	child3D.numInstances = num
}

func (child3D *Child3D) SetInstanceRenderDistance(dist float32) {
	child3D.projectionMatrix = mgl32.Perspective(
		mgl32.DegToRad(45),
		float32(child3D.config.ScreenWidth)/float32(child3D.config.ScreenHeight),
		0.1, dist,
	)
}
