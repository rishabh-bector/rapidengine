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

// Function prototypes
vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir);
vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);

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