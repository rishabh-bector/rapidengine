package material

type Material interface {
	Render(delta float64, darkness float32, totalTime float64)

	GetShader() *ShaderProgram
}
