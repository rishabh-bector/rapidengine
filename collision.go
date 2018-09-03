package rapidengine

import (
	"math"
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
}

// CollisionLink contains the data for a single collision
// between a child and a group, the groupname which the child
// should collide with, and a callback function to call when
// this collision happens.
type CollisionLink struct {
	group    string
	callback func(int)
}

// NewCollisionControl creates a new CollisionControl
func NewCollisionControl() CollisionControl {
	return CollisionControl{
		make(map[string][]Child),
		make(map[Child]CollisionLink),
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
func (c *CollisionControl) CreateCollision(child Child, group string, callback func(int)) {
	c.LinkMap[child] = CollisionLink{group, callback}
}

// CheckCollisionWithGroup checks if a child is colliding with
// any of the children in the passed group, including copies currently
// on the screen.
func (c *CollisionControl) CheckCollisionWithGroup(child Child, group string, camX, camY float32) int {
	for _, other := range c.GroupMap[group] {
		if !other.CheckCopyingEnabled() {
			if col := child.CheckCollision(other); col != 0 && child != other {
				return col
			}
		} else {
			for _, cpy := range other.GetCurrentCopies() {
				if col := child.CheckCollisionRaw(cpy.X, cpy.Y, other.GetCollider()); col != 0 {
					return col
				}
			}
		}
	}
	return 0
}

// Update is called once per frame, and checks for
// collisions of all children in the LinkMap
func (c *CollisionControl) Update(camX, camY float32) {
	for child, link := range c.LinkMap {
		if col := c.CheckCollisionWithGroup(child, link.group, camX, camY); col != 0 {
			link.callback(col)
		}
	}
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
func (collider *Collider) CheckCollision(x, y, otherX, otherY float32, otherCollider *Collider) int {
	dx := float64((x + collider.offsetX + collider.width/2) - (otherX + otherCollider.width/2))
	dy := float64((y + collider.offsetY + collider.height/2) - (otherY + otherCollider.height/2))
	width := float64((collider.width + otherCollider.width) / 2)
	height := float64((collider.height + otherCollider.height) / 2)
	crossWidth := width * dy
	crossHeight := height * dx
	if math.Abs(dx) <= width && math.Abs(dy) <= height {
		if crossWidth > crossHeight {
			if crossWidth > -crossHeight {
				return 4
			}
			return 3
		}
		if crossWidth > -crossHeight {
			return 1
		}
		return 2
	}

	/*if x+collider.offsetX < otherX+otherCollider.offsetX+otherCollider.width &&
		x+collider.offsetX+collider.width > otherX+otherCollider.offsetX &&
		y+collider.offsetY < otherY+otherCollider.offsetY+otherCollider.height &&
		y+collider.offsetY+collider.height > otherY+otherCollider.offsetY {
		return true
	}*/
	return 0
}
