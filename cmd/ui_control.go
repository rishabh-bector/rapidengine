package cmd

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/input"
	"rapidengine/ui"
)

//  --------------------------------------------------
//  UIControl manages the user interface system within
//  RapidEngine. The system is set up as a tree (one per scene):
//
//  - Each scene has an associated RootElement
//
//  - The RootElement has any number of child elements,
//    and those child elements can also have children
//
//  - To initialize, all the elements in each RootElement
//    are instanced into their respective scenes, and all
//    children of each element are prerendered.
//
//	- If the UI tree changes, such as if a new element is instanced,
//    that element is initialized and inserted into the tree.
//  --------------------------------------------------

type UIControl struct {
	roots map[*Scene]*ui.RootElement

	engine *Engine
}

func NewUIControl() UIControl {
	return UIControl{}
}

func (uiControl *UIControl) Initialize(engine *Engine) {
	uiControl.engine = engine
	uiControl.roots = make(map[*Scene]*ui.RootElement)
}

func (uiControl *UIControl) Update(inputs *input.Input) {
	if currentRoot, ok := uiControl.roots[uiControl.engine.SceneControl.GetCurrentScene()]; ok {
		currentRoot.Update(inputs)
	} else {
		panic("UIControl: Could not determine root element")
	}
}

func (uiControl *UIControl) InstanceElement(e ui.Element, scene *Scene) {
	if root, ok := uiControl.roots[scene]; ok {
		root.InstanceElement(e)
		e.GetTransform().Parent = root
	}
}

// InitializeTrees traverses all UI trees and prerenders all children,
// instances each element into it's scene, and instances each textbox
// into it's scene.
func (uiControl *UIControl) InitializeTrees() {
	for scene, root := range uiControl.roots {
		for _, element := range root.GetElements() {
			for _, c := range element.GetChildren() {
				c.PreRender(uiControl.engine.Renderer.MainCamera)
			}

			scene.InstanceUIElement(element)

			for _, t := range element.GetTextBoxes() {
				if t != nil {
					scene.InstanceText(t)
				}
			}
		}
	}
}

func (uiControl *UIControl) SetupChild(loc *child.Child2D) {
	loc.AttachMaterial(uiControl.engine.Renderer.DefaultMaterial1)
	loc.AttachMesh(geometry.NewRectangle())
}

func (uiControl *UIControl) SceneSetup(scene *Scene) {
	uiControl.roots[scene] = uiControl.NewRootElement(
		float32(uiControl.engine.Config.ScreenWidth),
		float32(uiControl.engine.Config.ScreenHeight),
	)
}

//  --------------------------------------------------
//  Element Constructors
//  --------------------------------------------------

func (uiControl *UIControl) NewRootElement(width, height float32) *ui.RootElement {
	root := ui.NewRootElement(width, height)
	return &root
}

func (uiControl *UIControl) NewUIButton(x, y, width, height float32) *ui.Button {
	button := ui.NewUIButton(x, y, width, height)

	button.ButtonChild = uiControl.engine.ChildControl.NewChild2D()
	uiControl.SetupChild(button.ButtonChild)

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
