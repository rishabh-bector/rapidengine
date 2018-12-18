#version 410

out vec3 WorldPos_CS_in;
out vec3 TexCoord_CS_in;
out vec3 MatCoord_CS_in;
out vec3 Normal_CS_in;

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

uniform vec3 viewPos;

float getDisplacement();
float getTerrainDisplacement();

vec3 getTerrainNormal();

void main() {
    vec3 finalPosition = vec3(position.x, position.y + getTerrainDisplacement(), position.z);
    vec4 mPos = modelMtx * vec4(finalPosition, 1.0);

    //	Vertex position 
    //gl_Position = projectionMtx * viewMtx * modelMtx * vec4(finalPosition, 1.0);
    WorldPos_CS_in = position; 

    // Normal vector
    Normal = mat3(transpose(inverse(modelMtx))) * getTerrainNormal();
    Normal_CS_in = normal;

    // Fragment position
    FragPos =  vec3(modelMtx * vec4(finalPosition, 1.0));

    // Texture coordinates
    MatCoord_CS_in = tex / scale;
    TexCoord_CS_in = tex;

    // Normal Mapping
    vec3 T = normalize(vec3(modelMtx * vec4(tangent,   0.0)));
    vec3 B = normalize(vec3(modelMtx * vec4(bitTangent, 0.0)));
    vec3 N = normalize(vec3(modelMtx * vec4(normal,    0.0)));
    TBN = mat3(T, B, N);
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
