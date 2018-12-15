#version 410

out VS_OUT {
    vec3 FragPos;
    vec3 TexCoords;
    vec3 Normal;
    mat3 TBN;

    float Visibility;
} vs_out;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitTangent;

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

uniform float totalTime;

uniform sampler2D terrainHeightMap;
uniform sampler2D terrainNormalMap;

uniform float terrainDisplacement;
uniform float terrainWidth;
uniform float terrainLength;

uniform float foliageDisplacement;
uniform float foliageNoiseSeed;
uniform float foliageVariation;

const float fogDensity = 0.007;
const float fogGradient = 1.5;

const float windAmplitude = 0.4;
const float windFrequency = 1;

float rand(vec2 c) {
	return fract(sin(dot(c.xy, vec2(12.9898,78.233))) * 437.5453);
}

float random(float seed, float minimum, float maximum) {
    //float initial = fract(sin(seed) * 100000.0);
    float initial = rand(vec2(seed, foliageNoiseSeed));
    return minimum + initial * (maximum - minimum);
}

float getFoliageHeight(vec4 mPos) {
    return position.y + foliageDisplacement + texture(terrainHeightMap, vec2(mPos.x / terrainWidth, mPos.z / terrainLength)).r * terrainDisplacement;
}

vec3 getTerrainNormal(vec4 mPos) {
    return texture(terrainNormalMap, vec2(mPos.x / terrainWidth, mPos.z / terrainLength)).rgb;
}

vec3 getInstancePosition() {
    float xAdd = random(float(gl_InstanceID) / terrainWidth, 0.0, terrainLength);
    float zAdd = random(float(gl_InstanceID + 100) / terrainLength, 0.0, terrainWidth);
    return vec3(position.x + float(xAdd), position.y, position.z + float(zAdd));
}

float getFogVisibility(vec4 mPos) {
    vec4 positionRelativeCamera = viewMtx * mPos;
    float dist = length(positionRelativeCamera.xyz);
    return clamp(exp(-pow((dist * fogDensity), fogGradient)), 0.0, 1.0);
}

void main() {
    vec3 instancePosition = getInstancePosition();
    vec4 mPos = modelMtx * vec4(instancePosition, 1.0);
    vec3 finalPosition = vec3(instancePosition.x, getFoliageHeight(mPos), instancePosition.z);

    // Wind simulation
    if(tex.y == 0) {
       finalPosition.x += sin((totalTime + gl_InstanceID) * windFrequency) * sin(totalTime - gl_InstanceID) * windAmplitude;
       finalPosition.z += sin((totalTime - gl_InstanceID) * windFrequency) * windAmplitude;
       finalPosition.y += random(sqrt(gl_InstanceID), -foliageVariation, foliageVariation);
    }

    //	Vertex position 
    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(finalPosition, 1.0);

    // Normal vector
    vs_out.Normal = mat3(transpose(inverse(modelMtx))) * getTerrainNormal(mPos);

    // Fragment position
    vs_out.FragPos =  vec3(modelMtx * vec4(finalPosition, 1.0));

    // Texture coordinates
    vs_out.TexCoords = tex;

    // Normal Mapping
    vec3 T = normalize(vec3(modelMtx * vec4(tangent,   0.0)));
    vec3 B = normalize(vec3(modelMtx * vec4(bitTangent, 0.0)));
    vec3 N = normalize(vec3(modelMtx * vec4(normal,    0.0)));
    vs_out.TBN = mat3(T, B, N);

    // Fog
    vs_out.Visibility = getFogVisibility(mPos);
}