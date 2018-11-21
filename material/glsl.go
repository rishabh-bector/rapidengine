package material

var TextureProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/glsl/basic.frag",
	fragmentShader: "../rapidengine/material/glsl/basic.vert",
	uniformLocations: map[string]int32{
		"texture0":            0,
		"modelMtx":            0,
		"viewMtx":             0,
		"projectionMtx":       0,
		"darkness":            0,
		"transparencyEnabled": 0,
		"transparencyMap":     0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      1,
	},
}

var ColorLightingProgram = ShaderProgram{
	vertexShader:   ShaderColorLightingVertex,
	fragmentShader: ShaderColorLightingFragment,
	uniformLocations: map[string]int32{
		"textureScale":        0,
		"copyingEnabled":      0,
		"transparency":        0,
		"modelMtx":            0,
		"viewMtx":             0,
		"projectionMtx":       0,
		"materialType":        0,
		"diffuseMap":          0,
		"cubeDiffuseMap":      0,
		"darkness":            0,
		"color":               0,
		"shine":               0,
		"transparencyEnabled": 0,
		"transparencyMap":     0,
		"viewPos":             0,

		"dirLight.direction": 0,
		"dirLight.ambient":   0,
		"dirLight.diffuse":   0,
		"dirLight.specular":  0,

		"pointLights": 0,
		"lmao":        0,
		"spotLight":   0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      1,
		"normal":   2,
	},
}

var SkyBoxProgram = ShaderProgram{
	vertexShader:   ShaderSkyBoxVertex,
	fragmentShader: ShaderSkyBoxFragment,
	uniformLocations: map[string]int32{
		"color":         0,
		"transparency":  0,
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,
		"darkness":      0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
	},
}

const ShaderColorLightingVertex = `

	#version 410

	out vec3 Normal;
	out vec3 FragPos;
	out vec3 TexCoords;

	layout (location = 0) in vec3 position;
	layout (location = 1) in vec3 tex;
	layout (location = 2) in vec3 normal;

	uniform float textureScale;
	
	// Child copying
	uniform vec3 copyingEnabled;
	layout (location = 3) in vec3 copyPosition;

	uniform mat4 modelMtx;
	uniform mat4 viewMtx;
	uniform mat4 projectionMtx;

	void main() {
		//	Vertex position
		if(copyingEnabled.x > 0) {
			gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position + gl_InstanceID, 1.0);
		} else {
			gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
		}

		// Normal vector
		Normal = mat3(transpose(inverse(modelMtx))) * normal;

		// Fragment position
		FragPos = vec3(modelMtx * vec4(position, 1.0));

		// Texture coordinates
		TexCoords = tex / textureScale;
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
	uniform sampler2D diffuseMap;
	uniform samplerCube cubeDiffuseMap;

	uniform vec3 color;
	uniform float shine;
	uniform float darkness;

	uniform int transparencyEnabled;
	uniform sampler2D transparencyMap;

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
			//result += CalcPointLight(lmao, norm, FragPos, viewDir);    
		}

		result *= darkness;
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
			if(texture(diffuseMap, TexCoords.xy).a != 1.0f) {
				discard;
			}

			if(transparencyEnabled == 1) {
				if(texture(transparencyMap, TexCoords.xy).x == 0.0f) {
					discard;
				} else if (texture(transparencyMap, TexCoords.xy).x != 1.0f) {
					return vec3(0, 0, 0);
				}
			}

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
			if(transparencyEnabled == 1) {
				if(texture(transparencyMap, TexCoords.xy).x == 0.0f) {
					discard;
				} else if(texture(transparencyMap, TexCoords.xy).x != 1.0f) {
					return vec3(0, 0, 0);
				}
			}
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
