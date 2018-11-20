package cmd

import (
	"rapidengine/child"
	"rapidengine/geometry"
)

type ChildControl struct {
	scene string

	scenes map[string]*Scene

	engine *Engine
}

func NewChildControl() ChildControl {
	return ChildControl{
		scenes: make(map[string]*Scene),
	}
}

func (cc *ChildControl) Initialize(engine *Engine) {
	cc.engine = engine
	cc.NewScene("scene1")
	cc.SetScene("scene1")
}

func (cc *ChildControl) PreRenderChildren() {
	for _, scn := range cc.scenes {
		for _, c := range scn.children {
			c.PreRender(cc.engine.Renderer.MainCamera)
		}
	}
}

func (cc *ChildControl) InstanceChild(c child.Child, scene string) {
	cc.scenes[scene].InstanceChild(c)
}

func (cc *ChildControl) NewChild2D() *child.Child2D {
	c := child.NewChild2D(cc.engine.Config)
	c.AttachMaterial(&cc.engine.Renderer.DefaultMaterial1)
	c.AttachMesh(geometry.NewRectangle())
	return c
}

func (cc *ChildControl) NewChild3D() *child.Child3D {
	return child.NewChild3D(cc.engine.Config)
}

func (cc *ChildControl) NewScene(id string) {
	cc.scenes[id] = NewScene()
}

func (cc *ChildControl) SetScene(id string) {
	cc.scene = id
}

func (cc *ChildControl) GetSceneChildren() []child.Child {
	return cc.scenes[cc.scene].GetChildren()
}

func (cc *ChildControl) GetScene() string {
	return cc.scene
}

func (cc *ChildControl) IsAutomaticRendering() bool {
	return cc.scenes[cc.scene].IsAutomaticRendering()
}

func (cc *ChildControl) DisableAutomaticRendering(scene string) {
	cc.scenes[scene].DisableAutomaticRendering()
}

type Scene struct {
	children []child.Child

	automaticRendering bool
}

func NewScene() *Scene {
	return &Scene{
		automaticRendering: true,
	}
}

func (s *Scene) InstanceChild(c child.Child) {
	s.children = append(s.children, c)
}

func (s *Scene) GetChildren() []child.Child {
	return s.children
}

func (s *Scene) IsAutomaticRendering() bool {
	return s.automaticRendering
}

func (s *Scene) DisableAutomaticRendering() {
	s.automaticRendering = false
}
