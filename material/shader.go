package material

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderProgram struct {
	id uint32

	vertexShader   string
	fragmentShader string

	geometryShader string

	controlShader string
	evalShader    string

	uniformLocations   map[string]int32
	attributeLocations map[string]uint32
}

func (shaderProgram *ShaderProgram) Bind() {
	b := shaderProgram.id
	gl.UseProgram(b)
}

func (shaderProgram *ShaderProgram) RebindAttribLocations() {
	for attrib, location := range shaderProgram.attributeLocations {
		gl.BindAttribLocation(shaderProgram.id, location, gl.Str(attrib+"\x00"))
	}
}

func (shaderProgram *ShaderProgram) GetUniform(name string) int32 {
	return shaderProgram.uniformLocations[name]
}

func (shaderProgram *ShaderProgram) GetID() uint32 {
	return shaderProgram.id
}

func (shaderProgram *ShaderProgram) Compile() {
	vert, err := ioutil.ReadFile(shaderProgram.vertexShader)
	if err != nil {
		panic(err)
	}
	vertexShader, err := CompileShader(string(vert)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	frag, err := ioutil.ReadFile(shaderProgram.fragmentShader)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := CompileShader(string(frag)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shaderProgram.id = gl.CreateProgram()
	gl.AttachShader(shaderProgram.id, vertexShader)
	gl.AttachShader(shaderProgram.id, fragmentShader)

	if shaderProgram.geometryShader != "" {
		geom, err := ioutil.ReadFile(shaderProgram.geometryShader)
		if err != nil {
			panic(err)
		}
		geometryShader, err := CompileShader(string(geom)+"\x00", gl.GEOMETRY_SHADER)
		if err != nil {
			panic(err)
		}
		gl.AttachShader(shaderProgram.id, geometryShader)
	}

	// Tesselation shaders
	if shaderProgram.controlShader != "" {
		println("Compiling tesselation shaders")
		cont, err := ioutil.ReadFile(shaderProgram.controlShader)
		if err != nil {
			panic(err)
		}
		controlShader, err := CompileShader(string(cont)+"\x00", gl.TESS_CONTROL_SHADER)
		if err != nil {
			panic(err)
		}
		gl.AttachShader(shaderProgram.id, controlShader)

		eval, err := ioutil.ReadFile(shaderProgram.evalShader)
		if err != nil {
			panic(err)
		}
		evalShader, err := CompileShader(string(eval)+"\x00", gl.TESS_EVALUATION_SHADER)
		if err != nil {
			panic(err)
		}
		gl.AttachShader(shaderProgram.id, evalShader)
	}

	gl.LinkProgram(shaderProgram.id)

	for uni := range shaderProgram.uniformLocations {
		shaderProgram.uniformLocations[uni] = gl.GetUniformLocation(shaderProgram.id, gl.Str(uni+"\x00"))
	}

	for attrib, location := range shaderProgram.attributeLocations {
		gl.BindAttribLocation(shaderProgram.id, location, gl.Str(attrib+"\x00"))
	}
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

//  --------------------------------------------------
//  Shader Programs
//  --------------------------------------------------

var BasicProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/basic/basic.vert",
	fragmentShader: "../rapidengine/material/shaders/basic/basic.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"flipped": 0,

		// Basic Material
		"diffuseLevel": 0,

		"hue":      0,
		"darkness": 0,

		"diffuseMap": 0,
		"scale":      0,

		"alphaMapLevel": 0,
		"alphaMap":      0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      1,
	},
}

var StandardProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/standard/standard.vert",
	fragmentShader: "../rapidengine/material/shaders/standard/standard.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		// Standard Material
		"diffuseMap":   0,
		"normalMap":    0,
		"heightMap":    0,
		"hue":          0,
		"diffuseLevel": 0,
		"displacement": 0,
		"scale":        0,
		"reflectivity": 0,
		"refractivity": 0,
		"refractLevel": 0,

		"cubeDiffuseMap": 0,

		// Lighting
		"dirLight.direction": 0,
		"dirLight.ambient":   0,
		"dirLight.diffuse":   0,
		"dirLight.specular":  0,

		"viewPos": 0,

		"numPointLights": 0,
	},
	attributeLocations: map[string]uint32{
		"position":   0,
		"tex":        1,
		"normal":     2,
		"tangent":    3,
		"bitTangent": 4,
	},
}

var TerrainProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/terrain/terrain.vert",
	fragmentShader: "../rapidengine/material/shaders/terrain/terrain.frag",
	controlShader:  "../rapidengine/material/shaders/terrain/terrain.cont",
	evalShader:     "../rapidengine/material/shaders/terrain/terrain.eval",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		// Standard Material
		"diffuseMap":   0,
		"normalMap":    0,
		"heightMap":    0,
		"displacement": 0,
		"scale":        0,

		// Terrain Data
		"terrainHeightMap":    0,
		"terrainNormalMap":    0,
		"terrainDisplacement": 0,

		// Lighting
		"dirLight.direction": 0,
		"dirLight.ambient":   0,
		"dirLight.diffuse":   0,
		"dirLight.specular":  0,
		"viewPos":            0,
		"numPointLights":     0,
	},
	attributeLocations: map[string]uint32{
		"position":   0,
		"tex":        1,
		"normal":     2,
		"tangent":    3,
		"bitTangent": 4,
	},
}

var FoliageProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/foliage/foliage.vert",
	fragmentShader: "../rapidengine/material/shaders/foliage/foliage.frag",
	geometryShader: "../rapidengine/material/shaders/foliage/foliage.geom",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		// Standard Material
		"diffuseMap": 0,
		"normalMap":  0,
		"heightMap":  0,
		"opacityMap": 0,

		// Terrain Info
		"terrainHeightMap":    0,
		"terrainNormalMap":    0,
		"terrainDisplacement": 0,

		"terrainWidth":  0,
		"terrainLength": 0,

		"foliageDisplacement": 0,
		"foliageNoiseSeed":    0,
		"foliageVariation":    0,

		"totalTime": 0,

		// Lighting
		"dirLight.direction": 0,
		"dirLight.ambient":   0,
		"dirLight.diffuse":   0,
		"dirLight.specular":  0,
		"viewPos":            0,
		"numPointLights":     0,
	},
	attributeLocations: map[string]uint32{
		"position":   0,
		"tex":        1,
		"normal":     2,
		"tangent":    3,
		"bitTangent": 4,
	},
}

var WaterProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/water/water.vert",
	fragmentShader: "../rapidengine/material/shaders/water/water.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		// Standard Material
		"diffuseMap": 0,
		"normalMap":  0,

		"heightMap":    0,
		"displacement": 0,

		"scale": 0,

		// Lighting
		"dirLight.direction": 0,
		"dirLight.ambient":   0,
		"dirLight.diffuse":   0,
		"dirLight.specular":  0,

		"viewPos": 0,

		"numPointLights": 0,

		"totalTime": 0,
	},
	attributeLocations: map[string]uint32{
		"position":   0,
		"tex":        1,
		"normal":     2,
		"tangent":    3,
		"bitTangent": 4,
	},
}

var SkyBoxProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/skybox/skybox.vert",
	fragmentShader: "../rapidengine/material/shaders/skybox/skybox.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"cubeDiffuseMap": 0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      0,
	},
}

var PostFinalProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/postprocessing/final/final.vert",
	fragmentShader: "../rapidengine/material/shaders/postprocessing/final/final.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"screen":    0,
		"fboWidth":  0,
		"fboHeight": 0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      0,
	},
}

var PostHDRProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/postprocessing/hdr/hdr.vert",
	fragmentShader: "../rapidengine/material/shaders/postprocessing/hdr/hdr.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"screen":    0,
		"fboWidth":  0,
		"fboHeight": 0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      0,
	},
}

var PostHorizontalProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/postprocessing/blur/horizontal/horizontal.vert",
	fragmentShader: "../rapidengine/material/shaders/postprocessing/blur/horizontal/horizontal.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"screen":    0,
		"fboWidth":  0,
		"fboHeight": 0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      0,
	},
}

var PostVerticalProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/postprocessing/blur/vertical/vertical.vert",
	fragmentShader: "../rapidengine/material/shaders/postprocessing/blur/vertical/vertical.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"screen":    0,
		"fboWidth":  0,
		"fboHeight": 0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      0,
	},
}

var PostPreBloomProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/postprocessing/bloom/prebloom/prebloom.vert",
	fragmentShader: "../rapidengine/material/shaders/postprocessing/bloom/prebloom/prebloom.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"screen":    0,
		"fboWidth":  0,
		"fboHeight": 0,

		"bloomThreshold": 0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      0,
	},
}

var PostPostBloomProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/shaders/postprocessing/bloom/postbloom/postbloom.vert",
	fragmentShader: "../rapidengine/material/shaders/postprocessing/bloom/postbloom/postbloom.frag",
	uniformLocations: map[string]int32{
		// Vertices
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"screen":    0,
		"fboWidth":  0,
		"fboHeight": 0,

		"bloomInput":     0,
		"bloomIntensity": 0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      0,
	},
}
