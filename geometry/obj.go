package geometry

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

func LoadObj(path string) Mesh {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	vertices := []mgl32.Vec3{}
	normals := []mgl32.Vec3{}
	textures := []mgl32.Vec2{}

	line := strings.Split(scanner.Text(), " ")

	// Load vertex/normal/texture data into slices. This is randomly ordered.
	for scanner.Scan() {
		if line[0] == "v" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			z, _ := strconv.ParseFloat(line[3], 32)

			vertices = append(vertices, mgl32.Vec3{float32(x), float32(y), float32(z)})
		}

		if line[0] == "vn" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			z, _ := strconv.ParseFloat(line[3], 32)

			normals = append(normals, mgl32.Vec3{float32(x), float32(y), float32(z)})
		}

		if line[0] == "vt" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)

			textures = append(textures, mgl32.Vec2{float32(x), float32(y)})
		}

		if line[0] == "f" {
			break
		}

		line = strings.Split(scanner.Text(), " ")
	}

	indicesArray := []uint32{}

	var verticesArray []float32
	var normalsArray []float32
	var texturesArray []float32

	verticesArray = make([]float32, len(vertices)*3)
	normalsArray = make([]float32, len(vertices)*3)
	texturesArray = make([]float32, len(vertices)*3)

	// Load faces into arrays
	for scanner.Scan() {
		if line[0] == "f" {
			v1 := strings.Split(line[1], "/")
			v2 := strings.Split(line[2], "/")
			v3 := strings.Split(line[3], "/")

			// Vertex 1
			vertexIndex, _ := strconv.ParseInt(v1[0], 10, 32)
			textureIndex, _ := strconv.ParseInt(v1[1], 10, 32)
			normalIndex, _ := strconv.ParseInt(v1[2], 10, 32)

			vertexIndex--
			textureIndex--
			normalIndex--

			indicesArray = append(indicesArray, uint32(vertexIndex))

			currentTexture := textures[textureIndex]
			currentNormal := normals[normalIndex]

			texturesArray[vertexIndex*2] = currentTexture.X()
			texturesArray[vertexIndex*2+1] = 1 - currentTexture.Y()

			normalsArray[vertexIndex*3] = currentNormal.X()
			normalsArray[vertexIndex*3+1] = currentNormal.Y()
			normalsArray[vertexIndex*3+2] = currentNormal.Z()

			// Vertex 2
			vertexIndex, _ = strconv.ParseInt(v2[0], 10, 32)
			textureIndex, _ = strconv.ParseInt(v2[1], 10, 32)
			normalIndex, _ = strconv.ParseInt(v2[2], 10, 32)

			vertexIndex--
			textureIndex--
			normalIndex--

			indicesArray = append(indicesArray, uint32(vertexIndex))

			currentTexture = textures[textureIndex]
			currentNormal = normals[normalIndex]

			texturesArray[vertexIndex*2] = currentTexture.X()
			texturesArray[vertexIndex*2+1] = 1 - currentTexture.Y()

			normalsArray[vertexIndex*3] = currentNormal.X()
			normalsArray[vertexIndex*3+1] = currentNormal.Y()
			normalsArray[vertexIndex*3+2] = currentNormal.Z()

			// Vertex 3
			vertexIndex, _ = strconv.ParseInt(v3[0], 10, 32)
			textureIndex, _ = strconv.ParseInt(v3[1], 10, 32)
			normalIndex, _ = strconv.ParseInt(v3[2], 10, 32)

			vertexIndex--
			textureIndex--
			normalIndex--

			indicesArray = append(indicesArray, uint32(vertexIndex))

			//println(textureIndex)

			currentTexture = textures[textureIndex]
			currentNormal = normals[normalIndex]

			texturesArray[vertexIndex*2] = currentTexture.X()
			texturesArray[vertexIndex*2+1] = 1 - currentTexture.Y()

			normalsArray[vertexIndex*3] = currentNormal.X()
			normalsArray[vertexIndex*3+1] = currentNormal.Y()
			normalsArray[vertexIndex*3+2] = currentNormal.Z()
		}

		if line[0] == "usemtl" {
			break
		}

		line = strings.Split(scanner.Text(), " ")
	}

	for i, v := range vertices {
		verticesArray[i*3] = v.X()
		verticesArray[i*3+1] = v.Y()
		verticesArray[i*3+2] = v.Z()
	}

	return Mesh{
		id:          path,
		vao:         NewVertexArray(verticesArray, indicesArray),
		normals:     &normalsArray,
		texCoords:   &texturesArray,
		numVertices: int32(len(indicesArray)),
	}
}
