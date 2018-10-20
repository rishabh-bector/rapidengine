package rapidengine

import (
	"rapidengine/configuration"
	"rapidengine/input"
)

//  --------------------------------------------------
//  Collision.go contains CollisionControl,
//  which manages what group each Child is in and
//  whether linked children are colliding or not. It
//  also contains the Collider, which defines a collision
//  rectangle for a child and can check for collision
//  against other children.
//  --------------------------------------------------

// CollisionControl contains a map of groupnames -> children
// and a map of children -> collisionlinks. This way, each child
// with a collision link can be checked against all the children in
// the group of the link.
type CollisionControl struct {
	GroupMap map[string][]Child
	LinkMap  map[Child]CollisionLink

	MouseChildren    map[int]Child
	NumMouseChildren int
	MouseCollider    Collider

	config *configuration.EngineConfig
}

// CollisionLink contains the data for a single collision
// between a child and a group, the groupname which the child
// should collide with, and a callback function to call when
// this collision happens.
type CollisionLink struct {
	group    string
	callback func([]bool)
}

// NewCollisionControl creates a new CollisionControl
func NewCollisionControl(config *configuration.EngineConfig) CollisionControl {
	return CollisionControl{
		GroupMap:         make(map[string][]Child),
		LinkMap:          make(map[Child]CollisionLink),
		MouseChildren:    make(map[int]Child),
		NumMouseChildren: 0,
		MouseCollider:    Collider{0, 0, 5, 5},
		config:           config,
	}
}

// CreateGroup creates an empty child collision group
func (c *CollisionControl) CreateGroup(group string) {
	c.GroupMap[group] = []Child{}
}

// AddChildToGroup adds a child to a collision group
func (c *CollisionControl) AddChildToGroup(child Child, group string) {
	c.GroupMap[group] = append(c.GroupMap[group], child)
}

// CreateCollision adds a child/collisionlink pair to the LinkMap, so that
// collision will be checked for in Update()
func (c *CollisionControl) CreateCollision(child Child, group string, callback func([]bool)) {
	c.LinkMap[child] = CollisionLink{group, callback}
}

// CreateMouseCollision adds a child to the MouseChildren list to be checked against mouse coordinates
func (c *CollisionControl) CreateMouseCollision(child Child) {
	c.MouseChildren[c.NumMouseChildren] = child
	c.NumMouseChildren++
}

// CheckCollisionWithGroup checks if a child is colliding with
// any of the children in the passed group, including copies currently
// on the screen.
func (c *CollisionControl) CheckCollisionWithGroup(child Child, group string, camX, camY float32) []bool {
	out := []bool{false, false, false, false}
	for _, other := range c.GroupMap[group] {
		if !other.CheckCopyingEnabled() {
			if col := child.CheckCollision(other); col != 0 && child != other {
				out[col-1] = true
			}
		} else {
			for _, cpy := range other.GetCurrentCopies() {
				if col := child.CheckCollisionRaw(cpy.X, cpy.Y, other.GetCollider()); col != 0 {
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
func (c *CollisionControl) Update(camX, camY float32, inputs *input.Input) {
	for child, link := range c.LinkMap {
		if col := c.CheckCollisionWithGroup(child, link.group, camX, camY); col != nil {
			link.callback(col)
		}
	}
	mx, my := c.ScaleMouseCoords(inputs.MouseX, inputs.MouseY, camX, camY)
	for _, child := range c.MouseChildren {
		child.MouseCollisionFunc(child.CheckCollisionRaw(mx, -my, &c.MouseCollider) != 0)
	}
}

func (c *CollisionControl) ScaleMouseCoords(x, y float64, camX, camY float32) (float32, float32) {
	return float32(x) + camX - float32(c.config.ScreenWidth/2), (float32(y) - camY - float32(c.config.ScreenHeight/2))
}

// Collider contains data about a collision rect.
// All children with collision detection need one of these.
type Collider struct {
	offsetX float32
	offsetY float32
	width   float32
	height  float32
}

// NewCollider creates a new collision rect
func NewCollider(x, y, w, h float32) Collider {
	return Collider{x, y, w, h}
}

// CheckCollision checks for collision between 2 collision rects
// 0 - None
// 1 - Right
// 2 - Top
// 3 - Left
// 4 - Bottom
func (collider *Collider) CheckCollision(x, y, vx, vy, otherX, otherY float32, otherCollider *Collider) int {
	if x+collider.offsetX+collider.width > otherX &&
		x+collider.offsetX < otherX+otherCollider.width &&
		y+collider.offsetY+collider.height+vy > otherY &&
		y+collider.offsetY+vy < otherY+otherCollider.height {
		if vy < 0 {
			return 4
		}
		return 2
	}

	if x+collider.offsetX+collider.width+(-1*vx) > otherX &&
		x+collider.offsetX+(-1*vx) < otherX+otherCollider.width &&
		y+collider.offsetY+collider.height > otherY &&
		y+collider.offsetY < otherY+otherCollider.height {
		if vx < 0 {
			return 1
		}
		return 3
	}

	return 0
}
