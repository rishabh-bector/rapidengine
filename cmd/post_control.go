package cmd

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type PostControl struct {
	PostProcessingEnabled bool

	ScreenChild *child.Child2D

	FrameBuffer uint32

	DepthRenderBuffer uint32

	RenderedTexture uint32

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

	// Generate frame buffer
	gl.GenFramebuffers(1, &pc.FrameBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, pc.FrameBuffer)

	// Generate rendered texture
	gl.GenTextures(1, &pc.RenderedTexture)
	gl.BindTexture(gl.TEXTURE_2D, pc.RenderedTexture)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGB,
		int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight),
		0, gl.RGB, gl.UNSIGNED_BYTE, gl.PtrOffset(0),
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	// Generate depth buffer
	gl.GenRenderbuffers(1, &pc.DepthRenderBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, pc.DepthRenderBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(pc.engine.Config.ScreenWidth), int32(pc.engine.Config.ScreenHeight))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, pc.DepthRenderBuffer)

	// Configure framebuffer
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, pc.RenderedTexture, 0)
	drawBuffers := []uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &drawBuffers[0])

	// Check for errors
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("Framebuffer Invalid")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// Create screen child
	pc.ScreenChild = pc.engine.ChildControl.NewChild2D()
	pc.ScreenChild.AttachMaterial(material.NewPostProcessMaterial(pc.engine.ShaderControl.GetShader("postprocessing"), &pc.RenderedTexture))
	pc.ScreenChild.AttachMesh(geometry.NewRectangle())
	pc.ScreenChild.ScaleX = float32(pc.engine.Config.ScreenWidth)
	pc.ScreenChild.ScaleY = float32(pc.engine.Config.ScreenHeight)
	pc.ScreenChild.Static = true
	pc.ScreenChild.SetPosition(0, 0)
	pc.ScreenChild.PreRender(pc.engine.Renderer.MainCamera)
}

func (pc *PostControl) UpdateFrameBuffers() {
	if pc.PostProcessingEnabled {
		gl.BindFramebuffer(gl.FRAMEBUFFER, pc.FrameBuffer)
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}

func (pc *PostControl) Update() {
	if !pc.PostProcessingEnabled {
		return
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	pc.engine.Renderer.RenderChild(pc.ScreenChild)
}
