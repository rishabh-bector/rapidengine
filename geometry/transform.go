package geometry

type Transform struct {
	X float32
	Y float32
	Z float32

	SX float32
	SY float32
	SZ float32
}

func NewTransform(x, y, z, sx, sy, sz float32) Transform {
	return Transform{x, y, z, sx, sy, sz}
}
