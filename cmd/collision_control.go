package cmd

import (
	"rapidengine/child"
	"rapidengine/configuration"
	"rapidengine/input"
	"rapidengine/physics"
)

//  --------------------------------------------------
//  CollisionControl manages what group each Child
//  is in and whether linked children are colliding or not.
//  --------------------------------------------------

// CollisionControl contains a map of groupnames -> children
// and a map of children -> collisionlinks. This way, each child
// with a collision link can be checked against all the children in
// the group of the link.
type CollisionControl struct {
	GroupMap map[string][]child.Child
	LinkMap  map[child.Child]physics.CollisionLink

	MouseChildren    map[int]child.Child
	NumMouseChildren int
	MouseCollider    physics.Collider

	config *configuration.EngineConfig
}

// NewCollisionControl creates a new CollisionControl
func NewCollisionControl(config *configuration.EngineConfig) CollisionControl {
	return CollisionControl{
		GroupMap:         make(map[string][]child.Child),
		LinkMap:          make(map[child.Child]physics.CollisionLink),
		MouseChildren:    make(map[int]child.Child),
		NumMouseChildren: 0,
		MouseCollider: physics.Collider{
			OffsetX: 0,
			OffsetY: 0,
			Width:   5,
			Height:  5,
		},
		config: config,
	}
}

// CreateGroup creates an empty child collision group
func (collisionControl *CollisionControl) CreateGroup(group string) {
	collisionControl.GroupMap[group] = []child.Child{}
}

// AddChildToGroup adds a child to a collision group
func (collisionControl *CollisionControl) AddChildToGroup(c child.Child, group string) {
	collisionControl.GroupMap[group] = append(collisionControl.GroupMap[group], c)
}

// CreateCollision adds a child/collisionlink pair to the LinkMap, so that
// collision will be checked for in Update()
func (collisionControl *CollisionControl) CreateCollision(c child.Child, group string, callback func([]bool)) {
	collisionControl.LinkMap[c] = physics.CollisionLink{group, callback}
}

// CreateMouseCollision adds a child to the MouseChildren list to be checked against mouse coordinates
func (collisionControl *CollisionControl) CreateMouseCollision(c child.Child) {
	collisionControl.MouseChildren[collisionControl.NumMouseChildren] = c
	collisionControl.NumMouseChildren++
}

// CheckCollisionWithGroup checks if a child is colliding with
// any of the children in the passed group, including copies currently
// on the screen.
func (collisionControl *CollisionControl) CheckCollisionWithGroup(c child.Child, group string, camX, camY float32) []bool {
	out := []bool{false, false, false, false}
	for _, other := range collisionControl.GroupMap[group] {
		if !other.CheckCopyingEnabled() {
			if col := c.CheckCollision(other); col != 0 && c != other {
				out[col-1] = true
			}
		} else {
			for _, cpy := range other.GetCurrentCopies() {
				if col := c.CheckCollisionRaw(cpy.X, cpy.Y, other.GetCollider()); col != 0 {
					out[col-1] = true
				}
			}
		}
	}
	return out
}

// Update is called once per frame, and checks for
// collisions of all children in the LinkMap. It also
// checks for collisions with the mouse with all
// children in the MouseChildren map
func (collisionControl *CollisionControl) Update(camX, camY float32, inputs *input.Input) {
	for c, link := range collisionControl.LinkMap {
		if col := collisionControl.CheckCollisionWithGroup(c, link.Group, camX, camY); col != nil {
			link.Callback(col)
		}
	}
	mx, my := collisionControl.ScaleMouseCoords(inputs.MouseX, inputs.MouseY, camX, camY)
	for _, c := range collisionControl.MouseChildren {
		c.MouseCollisionFunc(c.CheckCollisionRaw(mx, -my, &collisionControl.MouseCollider) != 0)
	}
}

func (collisionControl *CollisionControl) ScaleMouseCoords(x, y float64, camX, camY float32) (float32, float32) {
	return float32(x) + camX - float32(collisionControl.config.ScreenWidth/2), (float32(y) - camY - float32(collisionControl.config.ScreenHeight/2))
}
