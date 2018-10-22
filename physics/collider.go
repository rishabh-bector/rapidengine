package physics

//  --------------------------------------------------
//  Collider defines a collision
//  rectangle for a child and can check for collision
//  against other children. CollisionLink contains the data
//  for a single collision between a child and a group,
//  the groupname which the child should collide with,
//  and a callback function to call when
//  this collision happens.
//  --------------------------------------------------

// CollisionLink defines a collision between a child and a group
type CollisionLink struct {
	Group    string
	Callback func([]bool)
}

// Collider contains data about a collision rect.
// All children with collision detection need one of these.
type Collider struct {
	OffsetX float32
	OffsetY float32
	Width   float32
	Height  float32
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
	if x+collider.OffsetX+collider.Width > otherX &&
		x+collider.OffsetX < otherX+otherCollider.Width &&
		y+collider.OffsetY+collider.Height+vy > otherY &&
		y+collider.OffsetY+vy < otherY+otherCollider.Height {
		if vy < 0 {
			return 4
		}
		return 2
	}

	if x+collider.OffsetX+collider.Width+(-1*vx) > otherX &&
		x+collider.OffsetX+(-1*vx) < otherX+otherCollider.Width &&
		y+collider.OffsetY+collider.Height > otherY &&
		y+collider.OffsetY < otherY+otherCollider.Height {
		if vx < 0 {
			return 1
		}
		return 3
	}

	return 0
}
