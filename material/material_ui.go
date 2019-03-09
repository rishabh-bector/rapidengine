package material

type MaterialUI interface {
	Render(delta float64, darkness float32, totalTime float64)
	GetShader() *ShaderProgram

	GetDiffuseScalar() *float32
	GetNormalScalar() *float32
	GetMetallicScalar() *float32
	GetRoughnessScalar() *float32
	GetAOScalar() *float32

	GetVertexDisplacement() *float32
	GetParallaxDisplacement() *float32

	GetScale() *float32

	SetRoughOrSmooth(bool)

	AttachDiffuseMap(m *Texture)
	AttachNormalMap(m *Texture)
	AttachHeightMap(m *Texture)
	AttachMetallicMap(m *Texture)
	AttachRoughnessMap(m *Texture)
	AttachAOMap(m *Texture)
}
