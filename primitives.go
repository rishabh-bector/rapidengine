package rapidengine

//  --------------------------------------------------
//	Primitives.go allows easy creation of basic shapes
//  by automatically creating and binding a VAO to the
//  primitive struct, which can be passed to a Child.
//  --------------------------------------------------

import (
	"rapidengine/configuration"
)

type Primitive struct {
	id          string
	vao         *VertexArray
	numVertices int32
}

// NormalizeSizes takes in a size in pixels and normalizes to [0, 1]
func NormalizeSizes(x, y, sw, sh float32) (float32, float32) {
	return x / float32(sw), y / float32(sh)
}

// GetPrimitiveCoords returns the appropriate texture coordinates
// for each primitive
func GetPrimitiveCoords(id string) []float32 {
	switch id {
	case "rectangle":
		return RectTextures
	}
	return RectTextures
}

// NewTriangle creates a new triangle based on 3 points and a shaders object
func NewTriangle(points []float32) Primitive {
	indices := []uint32{}
	for i := 0; i < len(points); i++ {
		indices = append(indices, uint32(i))
	}
	t := Primitive{
		"triangle",
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
func NewRectangle(width, height float32, config *configuration.EngineConfig) Primitive {
	w, h := NormalizeSizes(width, height, float32(config.ScreenWidth), float32(config.ScreenHeight))
	points := []float32{
		0, 0, 0,
		w * 2, 0, 0,
		w * 2, h * 2, 0,
		0, h * 2, 0,
	}
	indices := []uint32{
		0, 1, 2,
		2, 0, 3,
	}

	r := Primitive{
		"rectangle",
		NewVertexArray(
			points,
			indices,
		),
		int32(len(indices)),
	}
	return r
}

// NewCube creates a 3D cube primitive
func NewCube() Primitive {
	indices := []uint32{}
	for i := 0; i < len(CubePoints); i++ {
		indices = append(indices, uint32(i))
	}
	c := Primitive{
		"cube",
		NewVertexArray(
			CubePoints,
			indices,
		),
		int32(108),
	}
	return c
}

var RectTextures = []float32{
	0, 1, //0,
	1, 1, //0,
	1, 0, //0,
	0, 0, //0,
}

var RectNormals = []float32{
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
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

var CubeNormals = []float32{
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,

	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,

	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,

	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,

	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,

	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
}

var CubeTextures = []float32{
	0, 0, 0,
	1, 0, 0,
	1, 1, 0,
	1, 1, 0,
	0, 1, 0,
	0, 0, 0,

	0, 0, 0,
	1, 0, 0,
	1, 1, 0,
	1, 1, 0,
	0, 1, 0,
	0, 0, 0,

	1, 0, 0,
	1, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 0, 0,
	1, 0, 0,

	1, 0, 0,
	1, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 0, 0,
	1, 0, 0,

	0, 1, 0,
	1, 1, 0,
	1, 0, 0,
	1, 0, 0,
	0, 0, 0,
	0, 1, 0,

	0, 1, 0,
	1, 1, 0,
	1, 0, 0,
	1, 0, 0,
	0, 0, 0,
	0, 1, 0,
}

var CubeMapTextures = []float32{
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
}
