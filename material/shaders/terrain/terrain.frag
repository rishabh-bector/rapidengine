#version 410

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

out vec4 FragColor;

in vec3 FragPos;
in vec3 TexCoords;
in mat3 TBN;
in vec3 Normal;
in float Visibility;
const vec3 skyColor = vec3(0.9, 0.9, 1.0);

in vec3 TexCoord_FS_in;
in vec3 MatCoord_FS_in;
in float Visibility_FS_in;

uniform sampler2D diffuseMap;
uniform sampler2D normalMap;
uniform sampler2D heightMap;

uniform sampler2D terrainNormalMap;

uniform vec3 viewPos;

uniform DirLight dirLight;

#define MAX_LIGHTS 10
uniform int numPointLights;
uniform PointLight pointLights[MAX_LIGHTS];

vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir);
vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);
vec4 calculateDiffuseColor();

void main() {    
    //vec3 norm = texture(normalMap, TexCoords.xy).rgb;
    //norm = normalize(norm * 2.0 - 1.0);
    //norm = normalize(TBN * norm);
    //vec3 norm = normalize(Normal);
    vec3 norm = texture(terrainNormalMap, TexCoord_FS_in.xy).rgb;

    vec3 viewDir = normalize(viewPos - FragPos);

    // Directional lighting
    vec3 result = CalcDirLight(dirLight, norm, viewDir);

    // Point lighting
    for(int i = 0; i < numPointLights; i++) {
        //result += CalcPointLight(pointLights[i], norm, FragPos, viewDir);    
    }
    
    FragColor = mix(vec4(skyColor, 1.0), vec4(result, 1.0), Visibility_FS_in);
    //FragColor = vec4(0.5, 1.0, 0.0, 1.0);
}

vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir) {
    vec3 lightDir = normalize(-light.direction);

    // Diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // Specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 0.5);

    vec3 color = calculateDiffuseColor().xyz;
    vec3 ambient = light.ambient * color;
    vec3 diffuse = light.diffuse * diff * color;
    vec3 specular = light.specular * spec * color;
    
    return ambient + diffuse + specular;
}

vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir) {
    vec3 lightDir = normalize(light.position - fragPos);

    // Diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // Specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 0.5);

    // Attenuation
    float distance = length(light.position - fragPos);

    if(distance > 1.05) {
        //return vec3(0, 0, 0);
    }

    float attenuation = 1.0 / ((light.constant) + (light.linear * distance) + (light.quadratic * (distance * distance))); 

    vec3 color = calculateDiffuseColor().xyz;

    vec3 ambient = light.ambient * color;
    vec3 diffuse = light.diffuse * diff * color;
    vec3 specular = light.specular * spec * color;
    
    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;

    return ambient + diffuse + specular;
}

vec4 calculateDiffuseColor() {
    return texture(diffuseMap, MatCoord_FS_in.xy);
}

