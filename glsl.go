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

var SkyBoxProgram = ShaderProgram{
	vertexShader:   ShaderSkyBoxVertex,
	fragmentShader: ShaderSkyBoxFragment,
}

const ShaderTextureVertex = `

		#version 410

		uniform mat4 modelMtx;
		uniform mat4 viewMtx;
		uniform mat4 projectionMtx;

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec3 tex;

		out vec3 texCoord;
		
		void main() {
			texCoord = tex;
			gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
		}
	
	` + "\x00"

const ShaderTextureFragment = `

		#version 410

		uniform sampler2D texture0;

		in vec3 texCoord;

		out vec4 outColor;

		void main() {
			outColor = texture(texture0, texCoord.xy);
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

	out vec3 Normal;
	out vec3 FragPos;
	out vec3 TexCoords;

	layout (location = 0) in vec3 position;
	layout (location = 1) in vec3 tex;
	layout (location = 2) in vec3 normal;

	uniform mat4 modelMtx;
	uniform mat4 viewMtx;
	uniform mat4 projectionMtx;

	void main() {
		//	Vertex position
		gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);

		// Normal vector
		Normal = mat3(transpose(inverse(modelMtx))) * normal;

		// Fragment position
		FragPos = vec3(modelMtx * vec4(position, 1.0));

		// Texture coordinates
		TexCoords = tex;
	}

` + "\x00"

const ShaderColorLightingFragment = `

	#version 410
	out vec4 FragColor;

	struct DirLight {
		vec3 direction;
		
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};
	
	struct PointLight {
		vec3 position;
		
		float constant;
		float linear;
		float quadratic;
		
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};
	
	struct SpotLight {
		vec3 position;
		vec3 direction;
		float cutOff;
		float outerCutOff;
		
		float constant;
		float linear;
		float quadratic;
		
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;       
	};

	in vec3 Normal;
	in vec3 FragPos;
	in vec3 TexCoords;
	
	// Object Material
	// Solid, Texture, or CubeMap
	uniform vec3 materialType;
	uniform vec3 color;
	uniform sampler2D diffuseMap;
	uniform samplerCube cubeDiffuseMap;
	uniform float shine;

	// Camera position
	uniform vec3 viewPos;

	// Scene Lights
	const int NR_POINT_LIGHTS = 1;
	uniform DirLight dirLight;
	uniform PointLight pointLights[NR_POINT_LIGHTS];
	uniform PointLight lmao;
	uniform SpotLight spotLight;

	// Function prototypes
	vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir);
	vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);
	vec3 CalcSpotLight(SpotLight light, vec3 normal, vec3 fragPos, vec3 viewDir);

	void main() {    
		// properties
		vec3 norm = normalize(Normal);
		vec3 viewDir = normalize(viewPos - FragPos);

		// phase 1: directional lighting
		vec3 result = CalcDirLight(dirLight, norm, viewDir);

		// phase 2: point lights
		for(int i = 0; i < NR_POINT_LIGHTS; i++) {
			result += CalcPointLight(lmao, norm, FragPos, viewDir);    
		}

		// phase 3: spot light
		// result += CalcSpotLight(spotLight, norm, FragPos, viewDir);    
		
		FragColor = vec4(result, 1.0);
	}

	// Calculates the color when using a directional light.
	vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir) {
		vec3 lightDir = normalize(-light.direction);

		// diffuse shading
		float diff = max(dot(normal, lightDir), 0.0);

		// specular shading
		vec3 reflectDir = reflect(-lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), shine);

		// combine results
		vec3 ambient = vec3(1, 1, 1);
		vec3 diffuse = vec3(1, 1, 1);
		vec3 specular = vec3(1, 1, 1);
		
		if(materialType.x > 0) {
			ambient = light.ambient * color;
			diffuse = light.diffuse * diff * color;
			specular = light.specular * spec * color;
		}

		if(materialType.y > 0) {
			ambient = light.ambient * vec3(texture(diffuseMap, TexCoords.xy));
			diffuse = light.diffuse * diff * vec3(texture(diffuseMap, TexCoords.xy));
			specular = light.specular * spec * vec3(texture(diffuseMap, TexCoords.xy));
		}

		if(materialType.z > 0) {
			ambient = light.ambient * vec3(texture(cubeDiffuseMap, TexCoords));
			diffuse = light.diffuse * diff * vec3(texture(cubeDiffuseMap, TexCoords));
			specular = light.specular * spec * vec3(texture(cubeDiffuseMap, TexCoords));
		}

		
		return ambient + diffuse + specular;
	}

	// Calculates the color when using a point light
	vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir) {
		vec3 lightDir = normalize(light.position - fragPos);

		// diffuse shading
		float diff = max(dot(normal, lightDir), 0.0);

		// specular shading
		vec3 reflectDir = reflect(-lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), shine);

		// attenuation
		float distance = length(light.position - fragPos);

		if(distance > 1.05) {
			//return vec3(0, 0, 0);
		}

		float attenuation = 1.0 / ((light.constant) + (light.linear * distance) + (light.quadratic * (distance * distance))); 

		// combine results
		vec3 ambient = vec3(1, 1, 1);
		vec3 diffuse = vec3(1, 1, 1);
		vec3 specular = vec3(1, 1, 1);

		// combine results	
		if(materialType.x > 0) {
			ambient = light.ambient * color;
			diffuse = light.diffuse * diff * color;
			specular = light.specular * spec * color;
		}

		if(materialType.y > 0) {
			ambient = light.ambient * vec3(texture(diffuseMap, TexCoords.xy));
			diffuse = light.diffuse * diff * vec3(texture(diffuseMap, TexCoords.xy));
			specular = light.specular * spec * vec3(texture(diffuseMap, TexCoords.xy));
		}

		if(materialType.z > 0) {
			ambient = light.ambient * vec3(texture(cubeDiffuseMap, TexCoords));
			diffuse = light.diffuse * diff * vec3(texture(cubeDiffuseMap, TexCoords));
			specular = light.specular * spec * vec3(texture(cubeDiffuseMap, TexCoords));
		}

		ambient *= attenuation;
		diffuse *= attenuation;
		specular *= attenuation;

		return ambient + diffuse + specular;
	}

	/*
	// Calculates the color when using a spot light
	vec3 CalcSpotLight(SpotLight light, vec3 normal, vec3 fragPos, vec3 viewDir) {
		vec3 lightDir = normalize(light.position - fragPos);

		// diffuse shading
		float diff = max(dot(normal, lightDir), 0.0);

		// specular shading
		vec3 reflectDir = reflect(-lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), shine);

		// attenuation
		float distance = length(light.position - fragPos);
		float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

		// spotlight intensity
		float theta = dot(lightDir, normalize(-light.direction)); 
		float epsilon = light.cutOff - light.outerCutOff;
		float intensity = clamp((theta - light.outerCutOff) / epsilon, 0.0, 1.0);

		// combine results
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;

		if(textureEnabled == 100) {
			//ambient = light.ambient * vec3(texture(diffuseMap, TexCoords));
			//diffuse = light.diffuse * diff * vec3(texture(diffuseMap, TexCoords));
			//specular = light.specular * spec * vec3(texture(diffuseMap, TexCoords));
		} else {
			//ambient = light.ambient * color;
			//diffuse = light.diffuse * diff * color;
			//specular = light.specular * spec * color;
		}

		ambient *= attenuation * intensity;
		diffuse *= attenuation * intensity;
		specular *= attenuation * intensity;

		return (ambient + diffuse + specular);
	}
	*/

` + "\x00"

const ShaderSkyBoxVertex = `

	#version 410

	layout (location = 0) in vec3 position;

	out vec3 texCoord;

	uniform mat4 viewMtx;
	uniform mat4 projectionMtx;
	uniform mat4 modelMtx;
	
	void main() {
		texCoord = position;
		gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
	}

` + "\x00"

const ShaderSkyBoxFragment = `

	#version 410

	out vec4 FragColor;

	in vec3 texCoord;

	uniform samplerCube cubeDiffuseMap;

	void main() {
		FragColor = texture(cubeDiffuseMap, texCoord);
	}

` + "\x00"
