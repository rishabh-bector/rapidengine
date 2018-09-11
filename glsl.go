package rapidengine

var TextureProgram = ShaderProgram{
	vertexShader:   ShaderTextureVertex,
	fragmentShader: ShaderTextureFragment,
}

var ColorProgram = ShaderProgram{
	vertexShader:   ShaderColorVertex,
	fragmentShader: ShaderColorFragment,
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

const ShaderColorVertex = `

	#version 410

	uniform mat4 modelMtx;
	uniform mat4 viewMtx;
	uniform mat4 projectionMtx;

	layout (location = 0) in vec3 position;

	void main() {
		gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
	}

` + "\x00"

const ShaderColorFragment = `

	#version 410

	uniform vec3 color;

	out vec4 outColor;

	void main() {
		outColor = vec4(color, 1.0);
	}

` + "\x00"

const ShaderColorLightingVertex = `

		#version 410

		uniform mat4 modelMtx;
		uniform mat4 viewMtx;
		uniform mat4 projectionMtx;

		layout (location = 0) in vec3 position;
		layout (location = 2) in vec3 normal;

		out vec3 Normal;
		out vec3 FragPos;

		void main() {
			gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
			Normal = mat3(transpose(inverse(modelMtx))) * normal;  
			FragPos = vec3(modelMtx * vec4(position, 1.0));
		}

` + "\x00"

const ShaderColorLightingFragment = `

		#version 410

		uniform vec3 color;
		uniform vec3 light;
		uniform vec3 lightPos;
		uniform vec3 viewPos;

		in vec3 Normal;
		in vec3 FragPos;

		out vec4 outColor;

		vec3 lightColor = vec3(0.5, 0.5, 0.5);

		void main() {

			vec3 norm = normalize(Normal);
			vec3 lightDir = normalize(lightPos - FragPos);

			// Ambient Lighting

			vec3 ambient = light[0] * lightColor;

			// Diffuse Lighting

			float diff = max(dot(norm, lightDir), 0.0);
			vec3 diffuse = diff * lightColor;

			// Specular Lighting

			vec3 viewDir = normalize(viewPos - FragPos);
			vec3 reflectDir = reflect(-lightDir, norm);
			float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
			vec3 specular = light[2] * spec * lightColor;  


			outColor = vec4((ambient + diffuse + specular) * color, 1.0);
		}

` + "\x00"
