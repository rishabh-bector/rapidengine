package cmd

import (
	"rapidengine/child"
	"rapidengine/geometry"
)

type ChildControl struct {
	engine *Engine
}

func NewChildControl() ChildControl {
	return ChildControl{}
}

func (cc *ChildControl) Initialize(engine *Engine) {
	cc.engine = engine
}

func (cc *ChildControl) NewChild2D() *child.Child2D {
	c := child.NewChild2D(cc.engine.Config)
	c.AttachMaterial(cc.engine.Renderer.DefaultMaterial1)
	c.AttachMesh(geometry.NewRectangle())
	return c
}

func (cc *ChildControl) NewChild3D() *child.Child3D {
	return child.NewChild3D(cc.engine.Config)
}
