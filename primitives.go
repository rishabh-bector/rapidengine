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

/*func NewCube(shaders *Shaders) Prism {
	points := []float32{
		// Top
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
	}
	p := Prism{
		NewVertexArray(
			NewVertexBuffer(points),
			NewElementBuffer(indices),
			shaders.idList[0],
			shaders.idList[1],
			AttribConfig{
				name:   "position",
				index:  0,
				size:   3,
				xtype:  gl.FLOAT,
				stride: 6 * 4,
				offset: 0 * 4,
			},
		),
	}

	return p
}*/
