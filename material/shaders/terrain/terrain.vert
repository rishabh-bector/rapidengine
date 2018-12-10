#version 410

out vec3 FragPos;
out vec3 TexCoords;
out mat3 TBN;
out vec3 Normal;
out float Visibility;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitTangent;

uniform float scale;
uniform float displacement;

uniform sampler2D heightMap;

uniform sampler2D terrainHeightMap;
uniform sampler2D terrainNormalMap;
uniform float terrainDisplacement;

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

const float fogDensity = 0.007;
const float fogGradient = 1.5;

float getDisplacement();
float getTerrainDisplacement();

vec3 getTerrainNormal();

float getFogVisibility(vec4 mPos) {
    vec4 positionRelativeCamera = viewMtx * mPos;
    float dist = length(positionRelativeCamera.xyz);
    return clamp(exp(-pow((dist * fogDensity), fogGradient)), 0.0, 1.0);
}

void main() {
    vec3 finalPosition = vec3(position.x, position.y + getTerrainDisplacement(), position.z);
    vec4 mPos = modelMtx * vec4(finalPosition, 1.0);

    //	Vertex position 
    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(finalPosition, 1.0);

    // Normal vector
    Normal = mat3(transpose(inverse(modelMtx))) * getTerrainNormal();

    // Fragment position
    FragPos =  vec3(modelMtx * vec4(finalPosition, 1.0));

    // Texture coordinates
    TexCoords = tex / scale;

    // Normal Mapping
    vec3 T = normalize(vec3(modelMtx * vec4(tangent,   0.0)));
    vec3 B = normalize(vec3(modelMtx * vec4(bitTangent, 0.0)));
    vec3 N = normalize(vec3(modelMtx * vec4(normal,    0.0)));
    TBN = mat3(T, B, N);

    // Fog
    Visibility = getFogVisibility(mPos);
}

float getTerrainDisplacement() {
    return texture(terrainHeightMap, tex.xy).x * terrainDisplacement;
}

vec3 getTerrainNormal() {
    return texture(terrainNormalMap, tex.xy).rgb;
}

float getDisplacement() {
    return texture(heightMap, tex.xy / scale).x * displacement;
}
