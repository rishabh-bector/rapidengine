package main

type Primitive struct {
	vao         *VertexArray
	numVertices int32
}

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

func NewRectangle(points []float32, shaders *Shaders) Primitive {
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
