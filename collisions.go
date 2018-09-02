package rapidengine

//  --------------------------------------------------
//  Collisions.go contains CollisionControl,
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
	callback func()
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
func (c *CollisionControl) CreateCollision(child Child, group string, callback func()) {
	c.LinkMap[child] = CollisionLink{group, callback}
}

// CheckCollisionWithGroup checks if a child is colliding with
// any of the children in the passed group, including copies currently
// on the screen.
func (c *CollisionControl) CheckCollisionWithGroup(child Child, group string, camX, camY float32) bool {
	for _, other := range c.GroupMap[group] {
		if !other.CheckCopyingEnabled() {
			if child.CheckCollision(other) && child != other {
				return true
			}
		} else {
			for _, c := range other.GetCurrentCopies() {
				if child.CheckCollisionRaw(c.X, c.Y, other.GetCollider()) {
					return true
				}
			}
		}
	}
	return false
}

// Update is called once per frame, and checks for
// collisions of all children in the LinkMap
func (c *CollisionControl) Update(camX, camY float32) {
	for child, link := range c.LinkMap {
		if c.CheckCollisionWithGroup(child, link.group, camX, camY) {
			link.callback()
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
func (collider *Collider) CheckCollision(x, y, otherX, otherY float32, otherCollider *Collider) bool {
	if x+collider.offsetX < otherX+otherCollider.offsetX+otherCollider.width &&
		x+collider.offsetX+collider.width > otherX+otherCollider.offsetX &&
		y+collider.offsetY < otherY+otherCollider.offsetY+otherCollider.height &&
		y+collider.offsetY+collider.height > otherY+otherCollider.offsetY {
		return true
	}
	return false
}
