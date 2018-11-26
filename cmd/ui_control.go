package cmd

import (
	"rapidengine/geometry"
	"rapidengine/input"
	"rapidengine/ui"
)

type UIControl struct {
	Elements []ui.Element

	engine *Engine
}

func NewUIControl() UIControl {
	return UIControl{}
}

func (uiControl *UIControl) Initialize(engine *Engine) {
	uiControl.engine = engine
}

func (uiControl *UIControl) Update(inputs *input.Input) {
	for _, element := range uiControl.Elements {
		element.Update(inputs)
	}
}

func (uiControl *UIControl) InstanceElement(e ui.Element, scene *Scene) {
	for _, c := range e.GetChildren() {
		scene.InstanceChild(c)
	}
	for _, t := range e.GetTextBoxes() {
		if t != nil {
			scene.InstanceText(t)
		}
	}
	uiControl.addElement(e)
}

func (uiControl *UIControl) addElement(e ui.Element) {
	uiControl.Elements = append(uiControl.Elements, e)
}

func (uiControl *UIControl) AlignCenter(e ui.Element) {
	e.SetPosition(float32(uiControl.engine.Config.ScreenWidth/2)-e.GetTransform().SX/2, e.GetTransform().Y)
}

//  --------------------------------------------------
//  Element Constructors
//  --------------------------------------------------

func (uiControl *UIControl) NewUIButton(x, y, width, height float32) *ui.Button {
	button := ui.NewUIButton(x, y, width, height)

	button.ButtonChild = uiControl.engine.ChildControl.NewChild2D()
	button.ButtonChild.AttachMaterial(uiControl.engine.Renderer.DefaultMaterial2)
	button.ButtonChild.AttachMesh(geometry.NewRectangle())

	uiControl.engine.CollisionControl.CreateMouseCollision(button.ButtonChild)

	button.Initialize()

	return &button
}

func (uiControl *UIControl) NewProgressBar() *ui.ProgressBar {
	pb := ui.NewProgressBar(uiControl.engine.Config)

	pb.BackChild = uiControl.engine.ChildControl.NewChild2D()
	pb.BarChild = uiControl.engine.ChildControl.NewChild2D()

	pb.BackChild.AttachMaterial(uiControl.engine.Renderer.DefaultMaterial1)
	pb.BarChild.AttachMaterial(uiControl.engine.Renderer.DefaultMaterial2)

	pb.BarChild.AttachMesh(geometry.NewRectangle())
	pb.BackChild.AttachMesh(geometry.NewRectangle())

	pb.Initialize()

	return &pb
}

func (uiControl *UIControl) NewMenu() *ui.Menu {
	menu := ui.NewMenu(uiControl.engine.Config)

	menu.BackChild = uiControl.engine.ChildControl.NewChild2D()
	menu.BackChild.AttachMaterial(uiControl.engine.Renderer.DefaultMaterial1)
	menu.BackChild.AttachMesh(geometry.NewRectangle())

	menu.Initialize()

	return &menu
}
