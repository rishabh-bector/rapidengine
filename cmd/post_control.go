package cmd

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//  --------------------------------------------------
//  PostControl manages all the post processing effects.
//  These include HDR, Bloom, Reflections, and Water.
//  --------------------------------------------------

// PostControl has 2 ping pong buffers, 1 which hooks into the renderer
// causing the 3D scene to be rendered to it. Then, an effect
// chain is applied while swapping the 2 buffers, and the final buffer
// is rendered to the screen as a 2D quad.
type PostControl struct {
	PostProcessingEnabled bool

	// Post Processing Effects
	hdrEnabled      bool
	gaussianEnabled bool
	bloomEnabled    bool
	waterEnabled    bool

	ScreenChild    *child.Child2D
	ScreenMaterial *material.PostProcessMaterial

	// Ping pong buffers
	PBuffer1            EffectBuffers
	PBuffer2            EffectBuffers
	PIntermediateBuffer EffectBuffers

	// Gaussian Blur
	gaussianIterations int
	gaussianScale      int
	GaussianBuffer1    EffectBuffers
	GaussianBuffer2    EffectBuffers
	GaussianBuffer3    EffectBuffers

	// Bloom
	BloomThreshold float32
	BloomIntensity float32
	BloomOffsetX   int32
	BloomOffsetY   int32
	BloomBuffer1   EffectBuffers

	engine *Engine
}

func NewPostControl() PostControl {
	return PostControl{
		PostProcessingEnabled: false,
	}
}

func (pc *PostControl) Initialize(engine *Engine) {
	pc.engine = engine
}

func (pc *PostControl) EnablePostProcessing() {
	pc.PostProcessingEnabled = true

	// Create buffers
	pc.PBuffer1 = pc.NewEffectBuffers(int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight), true)
	pc.PBuffer2 = pc.NewEffectBuffers(int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight), true)
	pc.PIntermediateBuffer = pc.NewEffectBuffers(int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight), true)

	pc.ScreenMaterial = material.NewPostProcessMaterial(pc.engine.ShaderControl.GetShader("post_final"), &pc.PBuffer2.RenderedTexture)
	pc.ScreenMaterial.FboWidth = float32(pc.engine.Config.ScreenWidth)
	pc.ScreenMaterial.FboHeight = float32(pc.engine.Config.ScreenHeight)

	// Create screen child
	pc.ScreenChild = pc.engine.ChildControl.NewChild2D()
	pc.ScreenChild.AttachMaterial(pc.ScreenMaterial)
	pc.ScreenChild.AttachMesh(geometry.NewScreenQuad())
	pc.ScreenChild.ScaleX = float32(pc.engine.Config.ScreenWidth)
	pc.ScreenChild.ScaleY = float32(pc.engine.Config.ScreenHeight)
	pc.ScreenChild.Static = true
	pc.ScreenChild.SetPosition(0, 0)
	pc.ScreenChild.PreRender(pc.engine.Renderer.MainCamera)
}

func (pc *PostControl) DisablePostProcessing() {
	pc.PostProcessingEnabled = false
}

func (pc *PostControl) IsPostProcessingEnabled() bool {
	return pc.PostProcessingEnabled
}

// UpdateFrameBuffers is called at the beginning of the
// render loop. If post processing is enabled, the output
// framebuffer of the screen will be switched from the
// default framebuffer (0) to the initial buffer of the
// PostControl in preparation for the post processing stage.
func (pc *PostControl) UpdateFrameBuffers() {
	if pc.PostProcessingEnabled {
		gl.BindFramebuffer(gl.FRAMEBUFFER, pc.PBuffer1.FrameBuffer)
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}

// Update applies the post processing effect chain every frame.
// At this point, the scene has been rendered to the PostControl's
// initial buffers. After the effects have been applied, the
// final buffers are rendered onto the screen through the
// ScreenChild and ScreenMaterial
func (pc *PostControl) Update() {
	if !pc.PostProcessingEnabled {
		return
	}

	// Apply HDR
	if pc.hdrEnabled {
		pc.ApplyHDR(&pc.PBuffer1, &pc.PBuffer2)
		pc.SwapPingPongBuffers()
	}

	if pc.bloomEnabled {
		pc.ApplyPreBloom(&pc.PBuffer1, &pc.BloomBuffer1)
		pc.ApplyGaussianBlur(&pc.BloomBuffer1, &pc.PBuffer2)
		pc.ApplyPostBloom(&pc.PBuffer1, &pc.PBuffer2, &pc.BloomBuffer1)

		pc.PIntermediateBuffer = pc.PBuffer2
		pc.PBuffer2 = pc.BloomBuffer1
		pc.BloomBuffer1 = pc.PIntermediateBuffer

		pc.SwapPingPongBuffers()
	}

	if pc.gaussianEnabled {
		//pc.ApplyGaussianBlur(&pc.PBuffer1, &pc.PBuffer2)
		//pc.SwapPingPongBuffers()
	}

	// Render final buffer to screen
	pc.ScreenMaterial.ScreenMap = &pc.PBuffer1.RenderedTexture
	pc.ScreenMaterial.AttachShader(pc.engine.ShaderControl.GetShader("post_final"))

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	pc.engine.Renderer.RenderChild(pc.ScreenChild)
}

// ApplyHDR applies the High Dynamic Range post processing effect.
func (pc *PostControl) ApplyHDR(input, output *EffectBuffers) {
	pc.ScreenMaterial.ScreenMap = &input.RenderedTexture
	pc.ScreenMaterial.AttachShader(pc.engine.ShaderControl.GetShader("post_hdr"))

	output.BindAndClear()

	pc.engine.Renderer.RenderChild(pc.ScreenChild)
}

func (pc *PostControl) ApplyHorizontalGaussian(input, output *EffectBuffers) {
	pc.ScreenMaterial.ScreenMap = &input.RenderedTexture
	pc.ScreenMaterial.AttachShader(pc.engine.ShaderControl.GetShader("post_horizontal"))

	output.BindAndClear()

	pc.engine.Renderer.RenderChild(pc.ScreenChild)
}

func (pc *PostControl) ApplyVerticalGaussian(input, output *EffectBuffers) {
	pc.ScreenMaterial.ScreenMap = &input.RenderedTexture
	pc.ScreenMaterial.AttachShader(pc.engine.ShaderControl.GetShader("post_vertical"))

	output.BindAndClear()

	pc.engine.Renderer.RenderChild(pc.ScreenChild)
}

func (pc *PostControl) ApplyGaussianBlur(input, output *EffectBuffers) {
	gl.Viewport(0, 0, int32(pc.engine.Config.ScreenWidth/pc.gaussianScale), int32(pc.engine.Config.ScreenHeight/pc.gaussianScale))
	pc.ApplyHorizontalGaussian(input, &pc.GaussianBuffer1)
	pc.ApplyVerticalGaussian(&pc.GaussianBuffer1, &pc.GaussianBuffer2)
	pc.swapGaussianPingPongBuffers()

	for i := 0; i < pc.gaussianIterations; i++ {
		pc.ApplyHorizontalGaussian(&pc.GaussianBuffer1, &pc.GaussianBuffer2)
		pc.ApplyVerticalGaussian(&pc.GaussianBuffer2, &pc.GaussianBuffer1)
	}

	gl.Viewport(pc.BloomOffsetX, pc.BloomOffsetY, int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight))
	pc.ApplyHorizontalGaussian(&pc.GaussianBuffer1, &pc.GaussianBuffer3)
	pc.ApplyVerticalGaussian(&pc.GaussianBuffer3, output)
	gl.Viewport(0, 0, int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight))
}

func (pc *PostControl) ApplyPreBloom(input, output *EffectBuffers) {
	pc.ScreenMaterial.ScreenMap = &input.RenderedTexture
	pc.ScreenMaterial.AttachShader(pc.engine.ShaderControl.GetShader("post_prebloom"))
	pc.ScreenMaterial.GetShader().Bind()
	gl.Uniform1f(pc.ScreenMaterial.GetShader().GetUniform("bloomThreshold"), pc.BloomThreshold)

	output.BindAndClear()

	pc.engine.Renderer.RenderChild(pc.ScreenChild)
}

func (pc *PostControl) ApplyPostBloom(mainInput, bloomInput, output *EffectBuffers) {
	pc.ScreenMaterial.ScreenMap = &mainInput.RenderedTexture
	pc.ScreenMaterial.AttachShader(pc.engine.ShaderControl.GetShader("post_postbloom"))

	pc.ScreenMaterial.GetShader().Bind()
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, bloomInput.RenderedTexture)
	gl.Uniform1i(pc.ScreenMaterial.GetShader().GetUniform("bloomInput"), 1)

	gl.Uniform1f(pc.ScreenMaterial.GetShader().GetUniform("bloomIntensity"), pc.BloomIntensity)

	output.BindAndClear()

	pc.engine.Renderer.RenderChild(pc.ScreenChild)
}

// SwapPingPongBuffers swaps PBuffer1 and PBuffer2 so that
// the next effect in the post processing chain will have
// the correct input and output buffers.
func (pc *PostControl) SwapPingPongBuffers() {
	pc.PIntermediateBuffer = pc.PBuffer2
	pc.PBuffer2 = pc.PBuffer1
	pc.PBuffer1 = pc.PIntermediateBuffer
}

func (pc *PostControl) swapGaussianPingPongBuffers() {
	pc.PIntermediateBuffer = pc.GaussianBuffer2
	pc.GaussianBuffer2 = pc.GaussianBuffer1
	pc.GaussianBuffer1 = pc.PIntermediateBuffer
}

func (pc *PostControl) EnableHDR() {
	pc.hdrEnabled = true
}

func (pc *PostControl) EnableGaussianBlur(iterations int, scale int) {
	pc.gaussianEnabled = true
	pc.gaussianIterations = iterations
	pc.gaussianScale = scale

	pc.GaussianBuffer1 = pc.NewEffectBuffers(int32(pc.engine.Config.ScreenWidth/pc.gaussianScale), int32(pc.engine.Config.ScreenHeight/pc.gaussianScale), true)
	pc.GaussianBuffer2 = pc.NewEffectBuffers(int32(pc.engine.Config.ScreenWidth/pc.gaussianScale), int32(pc.engine.Config.ScreenHeight/pc.gaussianScale), true)
	pc.GaussianBuffer3 = pc.NewEffectBuffers(int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight), true)
}

func (pc *PostControl) EnableBloom(blurIterations int, blurScale int) {
	pc.EnableGaussianBlur(blurIterations, blurScale)
	pc.bloomEnabled = true
	pc.BloomThreshold = 0.7
	pc.BloomIntensity = 1
	pc.BloomBuffer1 = pc.NewEffectBuffers(int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight), true)
}

type EffectBuffers struct {
	FrameBuffer uint32

	DepthRenderBuffer uint32

	RenderedTexture uint32
}

func (eb *EffectBuffers) BindAndClear() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, eb.FrameBuffer)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}

func (pc *PostControl) NewEffectBuffers(width, height int32, highPrecision bool) EffectBuffers {
	frameBuffer := uint32(0)
	depthRenderBuffer := uint32(0)
	renderedTexture := uint32(0)

	// Generate frame buffer
	gl.GenFramebuffers(1, &frameBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)

	// Generate rendered texture
	gl.GenTextures(1, &renderedTexture)
	gl.BindTexture(gl.TEXTURE_2D, renderedTexture)

	if highPrecision {
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGBA16F,
			width, height,
			0, gl.RGBA, gl.FLOAT, gl.PtrOffset(0),
		)
	} else {
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGB,
			width, height,
			0, gl.RGB, gl.UNSIGNED_BYTE, gl.PtrOffset(0),
		)
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	// Generate depth buffer
	gl.GenRenderbuffers(1, &depthRenderBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, depthRenderBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, width, height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, depthRenderBuffer)

	// Configure framebuffer
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, renderedTexture, 0)
	drawBuffers := []uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &drawBuffers[0])

	// Check for errors
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("Framebuffer Invalid")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return EffectBuffers{
		FrameBuffer:       frameBuffer,
		DepthRenderBuffer: depthRenderBuffer,
		RenderedTexture:   renderedTexture,
	}
}
