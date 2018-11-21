package material

type Material interface {
	Render(delta float64, darkness float32)

	GetShader() *ShaderProgram
}
