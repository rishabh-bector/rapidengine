package material

type CustomMaterial struct {
	shader     *ShaderProgram
	RenderFunc func(delta float64, darkness float32, totalTime float64)
}

func NewCustomMaterial(shader *ShaderProgram) *CustomMaterial {
	return &CustomMaterial{
		shader: shader,
	}
}

func (sm *CustomMaterial) Render(delta float64, darkness float32, totalTime float64) {
	sm.RenderFunc(delta, darkness, totalTime)
}

func (sm *CustomMaterial) GetShader() *ShaderProgram {
	return sm.shader
}
