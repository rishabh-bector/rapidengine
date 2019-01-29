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

in vec3 FragPos;
in vec3 TexCoords;
in mat3 TBN;
in vec3 Normal;
in vec3 ReflectedVector;
in vec3 RefractedVector;

// Standard Material
uniform sampler2D diffuseMap;
uniform sampler2D normalMap;
uniform sampler2D specularMap;

uniform vec4 hue;
uniform float diffuseLevel;
uniform float reflectivity;
uniform float refractLevel;


// Environment Info
uniform samplerCube cubeDiffuseMap;
uniform vec3 viewPos;

uniform DirLight dirLight;

#define MAX_LIGHTS 10
uniform int numPointLights;
uniform PointLight pointLights[MAX_LIGHTS];

uniform vec3 plAmbient[MAX_LIGHTS];

vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir, vec4 inColor);
vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir, vec4 inColor);
vec4 calculateDiffuseColor();

vec3 calculateReflection() {
    return texture(cubeDiffuseMap, ReflectedVector).xyz;
}

vec3 calculateRefraction() {
    return texture(cubeDiffuseMap, RefractedVector).xyz;
}

void main() {    
    //vec3 norm = texture(normalMap, TexCoords.xy).rgb;
    //norm = normalize(norm * 2.0 - 1.0);
    //vec3 norm = normalize(TBN * norm);
    vec3 norm = normalize(Normal);

    vec3 viewDir = normalize(viewPos - FragPos);

    vec4 color = calculateDiffuseColor();
    if(color.a < 0.6) {
        discard;
    }

    // Directional lighting
    vec3 result = CalcDirLight(dirLight, norm, viewDir, color);

    // Point lighting
    for(int i = 0; i < numPointLights; i++) {
        result += CalcPointLight(pointLights[i], norm, FragPos, viewDir, color);    
    }
    
    FragColor = vec4(mix(result, mix(calculateReflection(), calculateRefraction(), refractLevel), reflectivity), 1.0);
}

vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir, vec4 inColor) {
    vec3 lightDir = normalize(-light.direction);

    // diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 0);

    vec3 color = inColor.xyz;
    
    vec3 ambient = light.ambient * color;
    vec3 diffuse = light.diffuse * diff * color;
    vec3 specular = light.specular * spec * color;
    
    return ambient + diffuse + specular;
}

vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir, vec4 inColor) {
    vec3 lightDir = normalize(light.position - fragPos);

    // diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);

    // attenuation
    float d = length(light.position - fragPos);

    if(d > 1.05) {
        //return vec3(0, 0, 0);
    }

    float attenuation = 1.0 / ((light.constant) + (light.linear * d) + (light.quadratic * (d * d))); 

    vec3 color = inColor.xyz;

    vec3 ambient = light.ambient * color;
    vec3 diffuse = light.diffuse * diff * color;
    vec3 specular = light.specular * spec * color * texture(specularMap, TexCoords.xy).x;
    
    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;

    return ambient + diffuse + specular;
}

vec4 calculateDiffuseColor() {
    return mix(hue / 255.0, texture(diffuseMap, TexCoords.xy), diffuseLevel);
}
