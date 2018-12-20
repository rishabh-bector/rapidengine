package cmd

import (
	"rapidengine/child"
	"rapidengine/ui"
)

type SceneControl struct {
	currentScene *Scene

	scenes []*Scene

	engine *Engine
}

func NewSceneControl() SceneControl {
	return SceneControl{}
}

func (sc *SceneControl) Initialize(engine *Engine) {
	sc.engine = engine
}

func (sc *SceneControl) InstanceScene(scn *Scene) {
	sc.scenes = append(sc.scenes, scn)
}

func (sc *SceneControl) PreRenderChildren() {
	for _, scn := range sc.scenes {
		for _, c := range scn.GetChildren() {
			c.PreRender(sc.engine.Renderer.MainCamera)
		}
	}
}

func (sc *SceneControl) SetCurrentScene(scn *Scene) {
	sc.ClearActivation()
	sc.currentScene = scn
	sc.currentScene.Activate()
}

func (sc *SceneControl) GetCurrentScene() *Scene {
	return sc.currentScene
}

func (sc *SceneControl) GetCurrentChildren() []child.Child {
	return sc.currentScene.GetChildren()
}

func (sc *SceneControl) GetCurrentTexts() []*ui.TextBox {
	return sc.currentScene.GetTexts()
}

func (sc *SceneControl) ClearActivation() {
	for _, scn := range sc.scenes {
		scn.Deactivate()
	}
}

type Scene struct {
	ID string

	children   []child.Child
	uiElements []ui.Element
	texts      []*ui.TextBox

	subscenes []*Scene

	active bool

	automaticRendering bool
}

func (sc *SceneControl) NewScene(id string) *Scene {
	s := &Scene{
		ID:                 id,
		automaticRendering: true,
		active:             true,
		texts:              []*ui.TextBox{},
	}

	if sc.engine.Config.ShowFPS {
		s.InstanceText(sc.engine.FPSBox)
	}

	return s
}

func (s *Scene) InstanceChild(c child.Child) {
	s.children = append(s.children, c)
}

func (s *Scene) InstanceUIElement(e ui.Element) {
	s.uiElements = append(s.uiElements, e)
}

func (s *Scene) InstanceText(t *ui.TextBox) {
	s.texts = append(s.texts, t)
}

func (s *Scene) InstanceSubscene(scn *Scene) {
	s.subscenes = append(s.subscenes, scn)
}

func (s *Scene) Activate() {
	s.active = true
	for _, c := range s.GetChildren() {
		c.Activate()
	}
}

func (s *Scene) Deactivate() {
	for _, c := range s.GetChildren() {
		c.Deactivate()
	}
	for _, scn := range s.subscenes {
		scn.Deactivate()
	}
	s.active = false
}

func (s *Scene) IsActive() bool {
	return s.active
}

func (s *Scene) GetChildren() []child.Child {
	children := []child.Child{}

	// Local children
	if s.active {
		children = s.children
	}

	// UI children
	for _, element := range s.uiElements {
		children = append(children, element.GetChildren()...)
	}

	// Subscene children
	for _, scn := range s.subscenes {
		if scn.active {
			children = append(children, scn.GetChildren()...)
		}
	}

	return children
}

func (s *Scene) GetTexts() []*ui.TextBox {
	texts := s.texts
	for _, scn := range s.subscenes {
		if scn.active {
			texts = append(texts, scn.GetTexts()...)
		}
	}
	return texts
}

func (s *Scene) IsAutomaticRendering() bool {
	return s.automaticRendering
}

func (s *Scene) DisableAutomaticRendering() {
	s.automaticRendering = false
}
