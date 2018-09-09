package rapidengine

var TextureProgram = ShaderProgram{
	vertexShader:   ShaderTextureVertex,
	fragmentShader: ShaderTextureFragment,
}

var ColorLightingProgram = ShaderProgram{
	vertexShader:   ShaderColorLightingVertex,
	fragmentShader: ShaderColorLightingFragment,
}

const ShaderTextureVertex = `

		#version 410

		uniform mat4 modelMtx;
		uniform mat4 viewMtx;
		uniform mat4 projectionMtx;

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 tex;

		out vec2 texCoord;
		
		void main() {
			texCoord = tex;
			gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
		}
	
	` + "\x00"

const ShaderTextureFragment = `

		#version 410

		uniform sampler2D texture0;

		in vec2 texCoord;

		out vec4 outColor;

		void main() {
			outColor = texture(texture0, texCoord);
		}
		
	` + "\x00"

const ShaderColorLightingVertex = `

		#version 410

		uniform mat4 modelMtx;
		uniform mat4 viewMtx;
		uniform mat4 projectionMtx;

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec4 color;

		out vec4 inColor;

		void main() {
			gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
			inColor = color;
		}

` + "\x00"

const ShaderColorLightingFragment = `

		#version 410

		in vec4 inColor;
		out vec4 outColor;

		void main() {
			outColor = inColor;
		}

` + "\x00"
