package geometry

//  --------------------------------------------------
//	Mesh.go allows easy creation of basic shapes
//  by automatically creating and binding a VAO to the
//  Mesh struct, which can be passed to a Child.
//  --------------------------------------------------

import (
	"rapidengine/configuration"
	"rapidengine/material"

	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {

	// Mesh type
	id string

	// VAO containing vertices & indices
	vao *VertexArray

	// Texture coordinates
	texCoords *[]float32

	// Face normals
	normals *[]float32

	// Number of vertices
	numVertices int32
}

func (p *Mesh) GetID() string {
	return p.id
}

func (p *Mesh) GetVAO() *VertexArray {
	return p.vao
}

func (p *Mesh) GetNumVertices() int32 {
	return p.numVertices
}

func (p *Mesh) GetTexCoords() *[]float32 {
	return p.texCoords
}

func (p *Mesh) GetNormals() *[]float32 {
	return p.normals
}

// NormalizeSizes takes in a size in pixels and normalizes to [0, 1]
func NormalizeSizes(x, y, sw, sh float32) (float32, float32) {
	return x / float32(sw), y / float32(sh)
}

// GetMeshCoords returns the appropriate texture coordinates
// for each Mesh
func GetMeshCoords(id string) []float32 {
	switch id {
	case "rectangle":
		return RectTextures
	}
	return RectTextures
}

// NewTriangle creates a new triangle Mesh based on 3 points
func NewTriangle(points []float32) Mesh {
	indices := []uint32{}
	for i := 0; i < len(points); i++ {
		indices = append(indices, uint32(i))
	}
	t := Mesh{
		id:          "triangle",
		vao:         NewVertexArray(points, indices),
		texCoords:   &RectTextures,
		normals:     &RectNormals,
		numVertices: int32(len(indices)),
	}
	return t
}

// NewRectangle creates a rectangle Mesh centered around the origin,
// based on a width and height value
func NewRectangle(width, height float32, config *configuration.EngineConfig) Mesh {
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

	r := Mesh{
		id:          "rectangle",
		vao:         NewVertexArray(points, indices),
		texCoords:   &RectTextures,
		normals:     &RectNormals,
		numVertices: int32(len(indices)),
	}
	return r
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
		texCoords:   &CubeTextures,
		normals:     &CubeNormals,
		numVertices: int32(108),
	}
	return c
}

func NewPlane(width, height, density int, heightData [][]float32) Mesh {
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
			if heightData == nil {
				vertices[vertexPointer*3+1] = 0
				normals[vertexPointer*3] = 0
				normals[vertexPointer*3+1] = 1
				normals[vertexPointer*3+2] = 0
			} else {
				vertices[vertexPointer*3+1] = heightData[i][j]
				normal := calculateNormal(i, j, heightData)
				normals[vertexPointer*3] = normal.X()
				normals[vertexPointer*3+1] = normal.Y()
				normals[vertexPointer*3+2] = normal.Z()
			}
			vertices[vertexPointer*3+2] = float32(i) / (float32(yCount) - 1) * float32(height)
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
		texCoords:   &textureCoords,
		normals:     &normals,
		numVertices: int32(len(indices)),
	}
}

func GetHeightMapData(path string, min, max float32) [][]float32 {
	img, err := material.LoadImage(path)
	if err != nil {
		panic(err)
	}

	heightRange := max - min

	data := make([][]float32, img.Bounds().Size().X)
	for i := range data {
		data[i] = make([]float32, img.Bounds().Size().Y)
	}

	for x := 0; x < img.Bounds().Size().X; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			data[x][y] = min + ((float32(r) / 65535) * heightRange)
		}
	}

	return data
}

func calculateNormal(x, z int, heights [][]float32) mgl32.Vec3 {
	if x < 1 || x > len(heights)-2 || z < 1 || z > len(heights[0])-2 {
		return mgl32.Vec3{0, 0, 0}
	}

	L := heights[x-1][z]
	R := heights[x+1][z]
	D := heights[x][z-1]
	U := heights[x][z+1]

	normal := mgl32.Vec3{L - R, 2, D - U}
	normal = normal.Normalize()

	return normal
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
	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,

	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,
}
