package geometry

// A model can be imported from a 3D object file
// format such as OBJ or STL, and can contain multiple
// meshes / shapes.

type Model struct {
	textures map[string]*uint32

	Meshes []Mesh
}

func LoadModel(path string) {
}
