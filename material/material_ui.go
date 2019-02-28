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

	AttachDiffuseMap(m *uint32)
	AttachNormalMap(m *uint32)
	AttachHeightMap(m *uint32)
	AttachMetallicMap(m *uint32)
	AttachRoughnessMap(m *uint32)
	AttachAOMap(m *uint32)
}
