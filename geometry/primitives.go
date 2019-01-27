package geometry

import (
	"math"
	"rapidengine/configuration"
	"rapidengine/material"

	"github.com/go-gl/mathgl/mgl32"
)

// NewTriangle creates a new triangle Mesh based on 3 points
func NewTriangle(points []float32) Mesh {
	indices := []uint32{}
	for i := 0; i < len(points); i++ {
		indices = append(indices, uint32(i))
	}
	t := Mesh{
		id:          "triangle",
		vao:         NewVertexArray(points, indices),
		texCoords:   RectTextures,
		normals:     RectNormals,
		numVertices: int32(len(indices)),
	}
	return t
}

// NewPolygon creates a mesh based on a radius and number of sides
func NewPolygon(radius float32, numSides int, config *configuration.EngineConfig) Mesh {
	vertices := make([]float32, (numSides+1)*3)
	indices := make([]uint32, numSides*3)

	normals := make([]float32, (numSides+1)*3)
	texCoords := make([]float32, (numSides+1)*3)

	sw := float32(config.ScreenWidth)
	sh := float32(config.ScreenHeight)

	vertexPointer := 3
	for i := 0; i < numSides; i++ {
		circleVertex := float64((float32(i) / float32(numSides)) * (2 * math.Pi))
		vertices[vertexPointer] = normalizeX(float32(math.Cos(circleVertex))*radius, sw)
		vertices[vertexPointer+1] = normalizeY(float32(math.Sin(circleVertex))*radius, sh)
		vertices[vertexPointer+2] = 0

		normals[vertexPointer] = 0
		normals[vertexPointer+1] = 0
		normals[vertexPointer+2] = 0

		texCoords[vertexPointer] = 0
		texCoords[vertexPointer+1] = 0
		texCoords[vertexPointer+2] = 0

		vertexPointer += 3
	}

	indexPointer := uint32(0)
	for i := 0; i < numSides-1; i++ {
		indices[indexPointer] = 0
		indices[indexPointer+1] = uint32(i + 1)
		indices[indexPointer+2] = uint32(i + 2)
		indexPointer += 3
	}

	indices[indexPointer] = 0
	indices[indexPointer+1] = uint32(numSides)
	indices[indexPointer+2] = 1

	return Mesh{
		id:          "circle",
		vao:         NewVertexArray(vertices, indices),
		texCoords:   texCoords,
		normals:     normals,
		numVertices: int32(len(indices)),
	}

}

// NewRectangle creates a rectangle Mesh centered around the origin,
// based on a width and height value
func NewRectangle() Mesh {
	points := []float32{
		0, 0, 0,
		1, 0, 0,
		1, 1, 0,
		0, 1, 0,
	}
	indices := []uint32{
		0, 1, 2,
		2, 0, 3,
	}

	r := Mesh{
		id:          "rectangle",
		vao:         NewVertexArray(points, indices),
		texCoords:   RectTextures,
		normals:     RectNormals,
		numVertices: int32(len(indices)),
	}
	return r
}

// NewScreenQuad creates a rectangle mesh that fills the screen.
// This is useful for post processing effects.
func NewScreenQuad() Mesh {
	indices := []uint32{
		0, 1, 2,
		2, 0, 3,
	}

	return Mesh{
		id:          "screen",
		vao:         NewVertexArray(ScreenQuadPoints, indices),
		texCoords:   RectTextures,
		numVertices: int32(len(indices)),
	}
}

// NewCube creates a 3D cube Mesh
func NewCube() Mesh {
	indices := []uint32{}
	for i := 0; i < len(CubePoints); i++ {
		indices = append(indices, uint32(i))
	}
	c := Mesh{
		id:          "cube",
		vao:         NewVertexArray(CubePoints, indices),
		texCoords:   CubeTextures,
		normals:     CubeNormals,
		numVertices: int32(108),
	}
	return c
}

func NewPlane(width, height, density int, heightData [][]float32, scale float32) Mesh {
	//segWidth := float32(width) / float32(xCount)
	//segHeight := float32(height) / float32(yCount)
	xCount := density
	yCount := density

	count := xCount * yCount

	vertices := make([]float32, count*3)
	indices := make([]uint32, 6*(xCount-1)*(yCount-1))
	normals := make([]float32, count*3)
	textureCoords := make([]float32, count*3)

	vertexPointer := 0
	for i := 0; i < xCount; i++ {
		for j := 0; j < yCount; j++ {

			vertices[vertexPointer*3] = float32(j) / (float32(xCount) - 1) * float32(width)
			vertices[vertexPointer*3+2] = float32(i) / (float32(yCount) - 1) * float32(height)

			if heightData == nil {
				vertices[vertexPointer*3+1] = 0

				normals[vertexPointer*3] = 0
				normals[vertexPointer*3+1] = 1
				normals[vertexPointer*3+2] = 0
			} else {
				coordX := float32(i) / float32(xCount) / scale
				coordY := float32(j) / float32(yCount) / scale

				coordX -= float32(int(coordX))
				coordY -= float32(int(coordY))

				vertices[vertexPointer*3+1] = heightData[int(coordX*float32(len(heightData)))][int(coordY*float32(len(heightData[0])))]

				normal := calculateNormal(i, j, heightData)
				normals[vertexPointer*3] = normal.X()
				normals[vertexPointer*3+1] = normal.Y()
				normals[vertexPointer*3+2] = normal.Z()
			}

			textureCoords[vertexPointer*3] = float32(j) / (float32(xCount) - 1)
			textureCoords[vertexPointer*3+1] = float32(i) / (float32(yCount) - 1)
			textureCoords[vertexPointer*3+2] = 0
			vertexPointer++
		}
	}

	pointer := 0
	for gz := 0; gz < xCount-1; gz++ {
		for gx := 0; gx < yCount-1; gx++ {
			topLeft := uint32((gz * xCount) + gx)
			topRight := topLeft + 1
			bottomLeft := uint32(((gz + 1) * yCount) + gx)
			bottomRight := bottomLeft + 1

			indices[pointer] = topLeft
			indices[pointer+1] = bottomLeft
			indices[pointer+2] = topRight
			indices[pointer+3] = topRight
			indices[pointer+4] = bottomLeft
			indices[pointer+5] = bottomRight

			pointer += 6
		}
	}

	return Mesh{
		id:          "plane",
		vao:         NewVertexArray(vertices, indices),
		texCoords:   textureCoords,
		normals:     normals,
		numVertices: int32(len(indices)),
	}
}

var Offset = float32(1.0)

func NewBillBoard() Mesh {

	basePoints := []mgl32.Vec4{
		mgl32.Vec4{0, 0, 0, 1},
		mgl32.Vec4{1, 0, 0, 1},
		mgl32.Vec4{1, 1, 0, 1},
		mgl32.Vec4{0, 1, 0, 1},
	}

	baseTex := []mgl32.Vec4{
		mgl32.Vec4{RectTextures[0], RectTextures[1], RectTextures[2]},
		mgl32.Vec4{RectTextures[3], RectTextures[4], RectTextures[5]},
		mgl32.Vec4{RectTextures[6], RectTextures[7], RectTextures[8]},
		mgl32.Vec4{RectTextures[9], RectTextures[10], RectTextures[11]},
	}

	baseNorm := []mgl32.Vec4{
		mgl32.Vec4{RectNormals[0], RectNormals[1], RectNormals[2]},
		mgl32.Vec4{RectNormals[3], RectNormals[4], RectNormals[5]},
		mgl32.Vec4{RectNormals[6], RectNormals[7], RectNormals[8]},
		mgl32.Vec4{RectNormals[9], RectNormals[10], RectNormals[11]},
	}

	rot30 := mgl32.HomogRotate3D(Offset*(math.Pi/180.0), mgl32.Vec3{0, 1, 0})
	rot60 := mgl32.HomogRotate3D(-45.0*(math.Pi/180.0), mgl32.Vec3{0, 1, 0})

	rot30 = mgl32.Ident4()
	rot60 = mgl32.Ident4()

	bill2 := []mgl32.Vec4{
		rot30.Mul4x1(basePoints[0]),
		rot30.Mul4x1(basePoints[1]),
		rot30.Mul4x1(basePoints[2]),
		rot30.Mul4x1(basePoints[3]),
	}

	bill3 := []mgl32.Vec4{
		rot60.Mul4x1(basePoints[0]),
		rot60.Mul4x1(basePoints[1]),
		rot60.Mul4x1(basePoints[2]),
		rot60.Mul4x1(basePoints[3]),
	}

	bill2Tex := []mgl32.Vec4{
		rot30.Mul4x1(baseTex[0]),
		rot30.Mul4x1(baseTex[1]),
		rot30.Mul4x1(baseTex[2]),
		rot30.Mul4x1(baseTex[3]),
	}

	bill3Tex := []mgl32.Vec4{
		rot60.Mul4x1(baseTex[0]),
		rot60.Mul4x1(baseTex[1]),
		rot60.Mul4x1(baseTex[2]),
		rot60.Mul4x1(baseTex[3]),
	}

	bill2Norm := []mgl32.Vec4{
		rot30.Mul4x1(baseNorm[0]),
		rot30.Mul4x1(baseNorm[1]),
		rot30.Mul4x1(baseNorm[2]),
		rot30.Mul4x1(baseNorm[3]),
	}

	bill3Norm := []mgl32.Vec4{
		rot60.Mul4x1(baseNorm[0]),
		rot60.Mul4x1(baseNorm[1]),
		rot60.Mul4x1(baseNorm[2]),
		rot60.Mul4x1(baseNorm[3]),
	}

	points := []float32{
		basePoints[0].X(), basePoints[0].Y(), basePoints[0].Z(),
		basePoints[1].X(), basePoints[1].Y(), basePoints[1].Z(),
		basePoints[2].X(), basePoints[2].Y(), basePoints[2].Z(),
		basePoints[3].X(), basePoints[3].Y(), basePoints[3].Z(),

		bill2[0].X(), bill2[0].Y(), bill2[0].Z(),
		bill2[1].X(), bill2[1].Y(), bill2[1].Z(),
		bill2[2].X(), bill2[2].Y(), bill2[2].Z(),
		bill2[3].X(), bill2[3].Y(), bill2[3].Z(),

		bill3[0].X(), bill3[0].Y(), bill3[0].Z(),
		bill3[1].X(), bill3[1].Y(), bill3[1].Z(),
		bill3[2].X(), bill3[2].Y(), bill3[2].Z(),
		bill3[3].X(), bill3[3].Y(), bill3[3].Z(),
	}

	indices := []uint32{
		0, 1, 2,
		2, 0, 3,
		4, 5, 6,
		6, 4, 7,
		8, 9, 10,
		10, 8, 11,
	}

	textures := []float32{
		baseTex[0].X(), baseTex[0].Y(), baseTex[0].Z(),
		baseTex[1].X(), baseTex[1].Y(), baseTex[1].Z(),
		baseTex[2].X(), baseTex[2].Y(), baseTex[2].Z(),
		baseTex[3].X(), baseTex[3].Y(), baseTex[3].Z(),

		bill2Tex[0].X(), bill2Tex[0].Y(), bill2Tex[0].Z(),
		bill2Tex[1].X(), bill2Tex[1].Y(), bill2Tex[1].Z(),
		bill2Tex[2].X(), bill2Tex[2].Y(), bill2Tex[2].Z(),
		bill2Tex[3].X(), bill2Tex[3].Y(), bill2Tex[3].Z(),

		bill3Tex[0].X(), bill3Tex[0].Y(), bill3Tex[0].Z(),
		bill3Tex[1].X(), bill3Tex[1].Y(), bill3Tex[1].Z(),
		bill3Tex[2].X(), bill3Tex[2].Y(), bill3Tex[2].Z(),
		bill3Tex[3].X(), bill3Tex[3].Y(), bill3Tex[3].Z(),
	}

	normals := []float32{
		baseNorm[0].X(), baseNorm[0].Y(), baseNorm[0].Z(),
		baseNorm[1].X(), baseNorm[1].Y(), baseNorm[1].Z(),
		baseNorm[2].X(), baseNorm[2].Y(), baseNorm[2].Z(),
		baseNorm[3].X(), baseNorm[3].Y(), baseNorm[3].Z(),

		bill2Norm[0].X(), bill2Norm[0].Y(), bill2Norm[0].Z(),
		bill2Norm[1].X(), bill2Norm[1].Y(), bill2Norm[1].Z(),
		bill2Norm[2].X(), bill2Norm[2].Y(), bill2Norm[2].Z(),
		bill2Norm[3].X(), bill2Norm[3].Y(), bill2Norm[3].Z(),

		bill3Norm[0].X(), bill3Norm[0].Y(), bill3Norm[0].Z(),
		bill3Norm[1].X(), bill3Norm[1].Y(), bill3Norm[1].Z(),
		bill3Norm[2].X(), bill3Norm[2].Y(), bill3Norm[2].Z(),
		bill3Norm[3].X(), bill3Norm[3].Y(), bill3Norm[3].Z(),
	}

	return Mesh{
		id:          "billboard",
		vao:         NewVertexArray(points, indices),
		numVertices: int32(len(indices)),
		texCoords:   textures,
		normals:     normals,
	}
}

func GetHeightMapData(path string, max float32) [][]float32 {
	img, err := material.LoadImageFullDepth(path)
	if err != nil {
		panic(err)
	}

	img.At(0, 0).RGBA()

	data := make([][]float32, img.Bounds().Size().X)
	for i := range data {
		data[i] = make([]float32, img.Bounds().Size().Y)
	}

	for x := 0; x < img.Bounds().Size().X; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			r, _, _, _ := img.At(y, x).RGBA()
			data[x][y] = ((float32(r) / 65535.0) * max)
		}
	}

	for x := 0; x < len(data); x++ {
		for y := 0; y < len(data[0]); y++ {
			temp := data[x][y]
			data[x][y] = data[y][x]
			data[y][x] = temp
		}
	}

	return data
}

func calculateNormal(x, z int, heights [][]float32) mgl32.Vec4 {
	if x < 1 || x > len(heights)-2 || z < 1 || z > len(heights[0])-2 {
		return mgl32.Vec4{0, 0, 0}
	}

	L := heights[x-1][z]
	R := heights[x+1][z]
	D := heights[x][z-1]
	U := heights[x][z+1]

	normal := mgl32.Vec4{2 * (R - L), 4, 2 * (D - U)}
	return normal.Normalize()
}
