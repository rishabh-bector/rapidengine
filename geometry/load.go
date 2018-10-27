package geometry

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray struct {
	id            uint32
	vertexBuffer  uint32
	elementBuffer uint32
}

func NewVertexArray(vertices []float32, elements []uint32) *VertexArray {
	var id uint32
	vertexArray := VertexArray{id: id}
	gl.GenVertexArrays(1, &vertexArray.id)
	gl.BindVertexArray(vertexArray.id)
	vertexArray.vertexBuffer = vertexArray.AddVertexAttribute(vertices, 0, 3)
	vertexArray.elementBuffer = vertexArray.AddElementAttribute(elements)
	return &vertexArray
}

func (vertexArray *VertexArray) AddVertexAttribute(data []float32, index, size int32) uint32 {
	gl.BindVertexArray(vertexArray.id)
	vbo := NewVertexBuffer(data)
	gl.VertexAttribPointer(
		uint32(index),
		size,
		gl.FLOAT,
		false,
		0,
		gl.PtrOffset(0),
	)
	return vbo
}

func (vertexArray *VertexArray) AddElementAttribute(data []uint32) uint32 {
	gl.BindVertexArray(vertexArray.id)
	veo := NewElementBuffer(data)
	return veo
}

func NewVertexBuffer(points []float32) uint32 {
	var vertexBufferID uint32
	gl.GenBuffers(1, &vertexBufferID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferID)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	return vertexBufferID
}

func NewElementBuffer(indices []uint32) uint32 {
	var elementBufferID uint32
	gl.GenBuffers(1, &elementBufferID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBufferID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	return elementBufferID
}

func (vertexArray *VertexArray) RebindVertexArray() {
	gl.BindVertexArray(vertexArray.id)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexArray.vertexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertexArray.elementBuffer)
}

func UnbindBuffers() {
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (vertexArray *VertexArray) GetID() uint32 {
	return vertexArray.id
}
