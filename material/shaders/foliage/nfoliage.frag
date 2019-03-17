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
in float Visibility;
in vec3 ReflectedVector;
in vec3 RefractedVector;

// Standard Material
uniform sampler2D diffuseMap;
uniform sampler2D normalMap;
uniform sampler2D specularMap;
uniform sampler2D heightMap;
uniform sampler2D opacityMap;

uniform float diffuseLevel;
uniform float normalLevel;
uniform float specularLevel;
uniform float heightLevel;

uniform float scale;

uniform vec4 hue;
uniform float reflectivity;
uniform float refractLevel;

// Environment maps
uniform samplerCube cubeDiffuseMap;
uniform vec3 viewPos;

const vec3 skyColor = vec3(0.9, 0.9, 1.0);

// Lights
uniform DirLight dirLight;
#define MAX_LIGHTS 10
uniform int numPointLights;
uniform PointLight pointLights[MAX_LIGHTS];
uniform vec3 plAmbient[MAX_LIGHTS];

vec3 CalcDirLight(DirLight light, vec3 norm, vec3 viewDir, vec4 diffuseColor, vec4 specularColor);
vec3 CalcPointLight(PointLight light, vec3 norm, vec3 viewDir, vec3 fragPos, vec4 diffuseColor, vec4 specularColor);

vec4 calculateDiffuseColor(vec2 uv);
vec4 calculateSpecularColor(vec2 uv);
float calculateHeight(vec2 uv);

vec3 calculateReflection() {
    return texture(cubeDiffuseMap, ReflectedVector).xyz;
}

vec3 calculateRefraction() {
    return texture(cubeDiffuseMap, RefractedVector).xyz;
}

vec2 parallaxMapping(vec3 viewDir) {
    // number of depth layers
    float minLayers = 8.0;
    float maxLayers = 32.0;
    float numLayers = 32.0;//mix(maxLayers, minLayers, abs(dot(vec3(0.0, 0.0, 1.0), viewDir)));  

    // calculate the size of each layer
    float layerDepth = 1.0 / numLayers;

    // depth of current layer
    float currentLayerDepth = 0.0;

    // the amount to shift the texture coordinates per layer (from vector P)
    vec2 P = viewDir.xy * heightLevel; 
    vec2 deltaTexCoords = P / numLayers;

    // get initial values
    vec2  currentTexCoords     = TexCoords.xy;
    float currentDepthMapValue = calculateHeight(currentTexCoords);
  
    while(currentLayerDepth < currentDepthMapValue) {
        // shift texture coordinates along direction of P
        currentTexCoords -= deltaTexCoords;

        // get depthmap value at current texture coordinates
        currentDepthMapValue = calculateHeight(currentTexCoords); 

        // get depth of next layer
        currentLayerDepth += layerDepth;  
    }

    // get texture coordinates before collision (reverse operations)
    vec2 prevTexCoords = currentTexCoords + deltaTexCoords;

    // get depth after and before collision for linear interpolation
    float afterDepth  = currentDepthMapValue - currentLayerDepth;
    float beforeDepth = calculateHeight(prevTexCoords) - currentLayerDepth + layerDepth;
 
    // interpolation of texture coordinates
    float weight = afterDepth / (afterDepth - beforeDepth);
    vec2 finalTexCoords = prevTexCoords * weight + currentTexCoords * (1.0 - weight);

    return finalTexCoords;  
}

void main() {
    if (texture(opacityMap, TexCoords.xy).r < 0.2) {
        discard;
    }

    vec3 norm = vec3(0, 0, 1);

    vec2 uvs = TexCoords.xy;

    if(heightLevel > 0) {
        vec3 tangentViewDir = normalize((TBN * viewPos) - (TBN * FragPos));
        uvs = parallaxMapping(tangentViewDir);
        if(uvs.x > 1.0 / scale || uvs.y > 1.0 / scale || uvs.x < 0.0 || uvs.y < 0.0) {
            discard;
        }
    }

    if(uvs.y < 0.02) {
        discard;
    }

    if(normalLevel == 1) {
        norm = normalize(texture(normalMap, uvs).rgb);
        norm = normalize((norm * 2.0) - 0.5);
        norm = normalize(TBN * norm);
    } else {
        norm = normalize(Normal);
    }

    vec3 viewDir = normalize(viewPos - FragPos);

    vec4 diffuseColor = calculateDiffuseColor(uvs);
    vec4 specularColor = calculateSpecularColor(uvs);

    // Directional lighting
    vec3 result = CalcDirLight(dirLight, norm, viewDir, diffuseColor, specularColor);

    // Point lighting
    for(int i = 0; i < numPointLights; i++) {
        result += CalcPointLight(pointLights[i], norm, viewDir, FragPos, diffuseColor, specularColor);    
    }
    
    vec4 res = vec4(mix(result, mix(calculateReflection(), calculateRefraction(), refractLevel), reflectivity), 1.0);
    FragColor = mix(vec4(skyColor, 1.0), vec4(result, 1.0), Visibility);
}

vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir, vec4 diffuseColor, vec4 specularColor) {
    vec3 lightDir = normalize(-light.direction);
    vec3 halfwayDir = normalize(lightDir + viewDir);

    // diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), 32);

    vec3 color = diffuseColor.xyz;
    
    vec3 ambient = light.ambient * color;
    vec3 diffuse = light.diffuse * diff * color;
    vec3 specular = light.specular * spec * color * specularColor.r;
    
    return ambient + diffuse + specular;
}

vec3 CalcPointLight(PointLight light, vec3 norm, vec3 viewDir, vec3 fragPos, vec4 diffuseColor, vec4 specularColor) {
    vec3 lightDir = normalize(light.position - FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);

    // diffuse shading
    float diff = max(dot(norm, lightDir), 0.0);

    // specular shading
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 32);

    // attenuation
    float d = length(light.position - fragPos);

    if(d > 1.05) {
        //return vec3(0, 0, 0);
    }

    float attenuation = 1.0 / ((light.constant) + (light.linear * d) + (light.quadratic * (d * d))); 

    vec3 color = diffuseColor.xyz;

    vec3 ambient = light.ambient * color;
    vec3 diffuse = light.diffuse * diff * color;
    vec3 specular = light.specular * spec * color * specularColor.r;
    
    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;

    return ambient + diffuse + specular;
}

vec4 calculateDiffuseColor(vec2 uv) {
    return mix(hue / 255.0, texture(diffuseMap, uv), diffuseLevel);
}

vec4 calculateSpecularColor(vec2 uv) {
    if(specularLevel == 0) {
        return vec4(0);
    } else {
        return texture(specularMap, uv);
    }
}

float calculateHeight(vec2 uv) {
    return 1.0 - texture(heightMap, uv).r;
}