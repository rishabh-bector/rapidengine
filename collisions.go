package rapidengine

//  --------------------------------------------------
//  Collisions.go contains CollisionControl,
//  which manages what group each Child is in and
//  whether linked children are colliding or not. It
//  also contains the Collider, which defines a collision
//  rectangle for a child and can check for collision
//  against other children.
//  --------------------------------------------------

type CollisionControl struct {
	GroupMap map[string][]Child
	LinkMap  map[Child]CollisionLink
}

type CollisionLink struct {
	group    string
	callback func()
}

func NewCollisionControl() CollisionControl {
	return CollisionControl{
		make(map[string][]Child),
		make(map[Child]CollisionLink),
	}
}

func (c *CollisionControl) CreateGroup(group string) {
	c.GroupMap[group] = []Child{}
}

func (c *CollisionControl) AddChildToGroup(child Child, group string) {
	c.GroupMap[group] = append(c.GroupMap[group], child)
}

func (c *CollisionControl) CreateCollision(child Child, group string, callback func()) {
	c.LinkMap[child] = CollisionLink{group, callback}
}

func (c *CollisionControl) CheckCollisionWithGroup(child Child, group string) bool {
	for _, other := range c.GroupMap[group] {
		if child.CheckCollision(other) && child != other {
			return true
		}
	}
	return false
}

func (c *CollisionControl) Update() {
	for child, link := range c.LinkMap {
		if c.CheckCollisionWithGroup(child, link.group) {
			link.callback()
		}
	}
}

type Collider struct {
	offsetX float32
	offsetY float32
	width   float32
	height  float32
}

func NewCollider(x, y, w, h float32) Collider {
	return Collider{x, y, w, h}
}

func (collider *Collider) CheckCollision(x, y float32, other Child) bool {
	if x+collider.offsetX < other.GetX()+other.GetCollider().offsetX+other.GetCollider().width &&
		x+collider.offsetX+collider.width > other.GetX()+other.GetCollider().offsetX &&
		y+collider.offsetY < other.GetY()+other.GetCollider().offsetY+other.GetCollider().height &&
		y+collider.offsetY+collider.height > other.GetY()+other.GetCollider().offsetY {
		return true
	}
	return false
}
