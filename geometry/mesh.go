package geometry

//  --------------------------------------------------
//	Mesh.go allows easy creation of basic shapes
//  by automatically creating and binding a VAO to the
//  Mesh struct, which can be passed to a Child.
//  --------------------------------------------------

import (
	"math"
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

	// Normal Mapping
	tangents   []float32
	bitangents []float32
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
	return x / sw, y / sh
}

func normalizeX(x, sw float32) float32 {
	return x / sw
}

func normalizeY(y, sh float32) float32 {
	return y / sh
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
		texCoords:   &texCoords,
		normals:     &normals,
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
			vertices[vertexPointer*3+2] = float32(i) / (float32(yCount) - 1) * float32(height)

			if heightData == nil {
				vertices[vertexPointer*3+1] = 0

				normals[vertexPointer*3] = 0
				normals[vertexPointer*3+1] = 1
				normals[vertexPointer*3+2] = 0
			} else {
				vertices[vertexPointer*3+1] = heightData[int((float32(i)/float32(xCount))*float32(len(heightData)))][int((float32(j)/float32(yCount))*float32(len(heightData[0])))]

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
		texCoords:   &textureCoords,
		normals:     &normals,
		numVertices: int32(len(indices)),
	}
}

func (m *Mesh) ComputeTangents() {
	for i := 0; i < len(m.vao.vertices); i++ {
		m.tangents = append(m.tangents, 0)
		m.bitangents = append(m.bitangents, 0)
	}

	for i := 0; i < int(len(m.vao.indices)); i += 3 {
		v0 := mgl32.Vec3{
			m.vao.vertices[m.vao.indices[i]*3],
			m.vao.vertices[m.vao.indices[i]*3+1],
			m.vao.vertices[m.vao.indices[i]*3+2],
		}
		v1 := mgl32.Vec3{
			m.vao.vertices[m.vao.indices[i+1]*3],
			m.vao.vertices[m.vao.indices[i+1]*3+1],
			m.vao.vertices[m.vao.indices[i+1]*3+2],
		}
		v2 := mgl32.Vec3{
			m.vao.vertices[m.vao.indices[i+2]*3],
			m.vao.vertices[m.vao.indices[i+2]*3+1],
			m.vao.vertices[m.vao.indices[i+2]*3+2],
		}

		texCoords := *m.texCoords
		uv0 := mgl32.Vec2{
			texCoords[m.vao.indices[i]*3],
			texCoords[m.vao.indices[i]*3+1],
		}
		uv1 := mgl32.Vec2{
			texCoords[m.vao.indices[i+1]*3],
			texCoords[m.vao.indices[i+1]*3+1],
		}
		uv2 := mgl32.Vec2{
			texCoords[m.vao.indices[i+2]*3],
			texCoords[m.vao.indices[i+2]*3+1],
		}

		e1 := v1.Sub(v0)
		e2 := v2.Sub(v0)

		deltaUV1 := uv1.Sub(uv0)
		deltaUV2 := uv2.Sub(uv0)

		r := float32(1) / (deltaUV1.X()*deltaUV2.Y() - deltaUV1.Y()*deltaUV2.X())

		tangent := (e1.Mul(deltaUV2.Y()).Sub(e2.Mul(deltaUV1.Y()))).Mul(r)
		bitangent := (e2.Mul(deltaUV1.X()).Sub(e1.Mul(deltaUV2.X()))).Mul(r)

		m.tangents[m.vao.indices[i]*3] += tangent.X()
		m.tangents[m.vao.indices[i]*3+1] += tangent.Y()
		m.tangents[m.vao.indices[i]*3+2] += tangent.Z()

		m.tangents[m.vao.indices[i+1]*3] += tangent.X()
		m.tangents[m.vao.indices[i+1]*3+1] += tangent.Y()
		m.tangents[m.vao.indices[i+1]*3+2] += tangent.Z()

		m.tangents[m.vao.indices[i+2]*3] += tangent.X()
		m.tangents[m.vao.indices[i+2]*3+1] += tangent.Y()
		m.tangents[m.vao.indices[i+2]*3+2] += tangent.Z()

		m.bitangents[m.vao.indices[i]*3] += bitangent.X()
		m.bitangents[m.vao.indices[i]*3+1] += bitangent.Y()
		m.bitangents[m.vao.indices[i]*3+2] += bitangent.Z()

		m.bitangents[m.vao.indices[i+1]*3] += bitangent.X()
		m.bitangents[m.vao.indices[i+1]*3+1] += bitangent.Y()
		m.bitangents[m.vao.indices[i+1]*3+2] += bitangent.Z()

		m.bitangents[m.vao.indices[i+2]*3] += bitangent.X()
		m.bitangents[m.vao.indices[i+2]*3+1] += bitangent.Y()
		m.bitangents[m.vao.indices[i+2]*3+2] += bitangent.Z()
	}

	m.GetVAO().AddVertexAttribute(m.tangents, 3, 3)
	m.GetVAO().AddVertexAttribute(m.bitangents, 4, 3)

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

	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
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

func calculateNormal(x, z int, heights [][]float32) mgl32.Vec3 {
	if x < 1 || x > len(heights)-2 || z < 1 || z > len(heights[0])-2 {
		return mgl32.Vec3{0, 0, 0}
	}

	L := heights[x-1][z]
	R := heights[x+1][z]
	D := heights[x][z-1]
	U := heights[x][z+1]

	normal := mgl32.Vec3{2 * (R - L), 4, 2 * (D - U)}
	return normal.Normalize()
}
