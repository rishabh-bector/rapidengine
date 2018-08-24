package main

//  --------------------------------------------------
//	Primitives.go allows easy creation of basic shapes
//  by automatically creating and binding a VAO to the
//  primitive struct, which can be passed to a Child.
//  --------------------------------------------------

type Primitive struct {
	vao         *VertexArray
	numVertices int32
}

// NormalizeSizes takes in a size in pixels and normalizes to [0, 1]
func NormalizeSizes(x, y float32) (float32, float32) {
	return x / float32(ScreenWidth), y / float32(ScreenHeight)
}

// NewTriangle creates a new triangle based on 3 points and a shaders object
func NewTriangle(points []float32, shaders *Shaders) Primitive {
	indices := []uint32{}
	for i := 0; i < len(points); i++ {
		indices = append(indices, uint32(i))
	}
	t := Primitive{
		NewVertexArray(
			points,
			indices,
		),
		int32(len(indices)),
	}
	return t
}

// NewRectangle creates a rectangle primitive centered around the origin,
// based on a width and height value
func NewRectangle(width, height float32, shaders *Shaders) Primitive {
	w, h := NormalizeSizes(width, height)
	points := []float32{
		0, 0, 0,
		w, 0, 0,
		w, -h, 0,
		0, -h, 0,
	}
	indices := []uint32{
		0, 1, 2,
		2, 0, 3,
	}

	r := Primitive{
		NewVertexArray(
			points,
			indices,
		),
		int32(len(indices)),
	}
	return r
}

// NewCube creates a 3D cube primitive
func NewCube(shaders *Shaders) Primitive {
	indices := []uint32{}
	for i := 0; i < len(CubePoints); i++ {
		indices = append(indices, uint32(i))
	}
	c := Primitive{
		NewVertexArray(
			CubePoints,
			indices,
		),
		int32(108),
	}
	return c
}

var CubePoints = []float32{
	-0.5, -0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, 0.5, -0.5,
	0.5, 0.5, -0.5,
	-0.5, 0.5, -0.5,
	-0.5, -0.5, -0.5,

	-0.5, -0.5, 0.5,
	0.5, -0.5, 0.5,
	0.5, 0.5, 0.5,
	0.5, 0.5, 0.5,
	-0.5, 0.5, 0.5,
	-0.5, -0.5, 0.5,

	-0.5, 0.5, 0.5,
	-0.5, 0.5, -0.5,
	-0.5, -0.5, -0.5,
	-0.5, -0.5, -0.5,
	-0.5, -0.5, 0.5,
	-0.5, 0.5, 0.5,

	0.5, 0.5, 0.5,
	0.5, 0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, -0.5, 0.5,
	0.5, 0.5, 0.5,

	-0.5, -0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, -0.5, 0.5,
	0.5, -0.5, 0.5,
	-0.5, -0.5, 0.5,
	-0.5, -0.5, -0.5,

	-0.5, 0.5, -0.5,
	0.5, 0.5, -0.5,
	0.5, 0.5, 0.5,
	0.5, 0.5, 0.5,
	-0.5, 0.5, 0.5,
	-0.5, 0.5, -0.5,
}

var CubeTextures = []float32{
	0, 0,
	1, 0,
	1, 1,
	1, 1,
	0, 1,
	0, 0,

	0, 0,
	1, 0,
	1, 1,
	1, 1,
	0, 1,
	0, 0,

	1, 0,
	1, 1,
	0, 1,
	0, 1,
	0, 0,
	1, 0,

	1, 0,
	1, 1,
	0, 1,
	0, 1,
	0, 0,
	1, 0,

	0, 1,
	1, 1,
	1, 0,
	1, 0,
	0, 0,
	0, 1,

	0, 1,
	1, 1,
	1, 0,
	1, 0,
	0, 0,
	0, 1,
}
