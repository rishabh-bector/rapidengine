#version 410

out vec3 FragPos;
out vec3 TexCoords;
out mat3 TBN;
out vec3 Normal;
out float Visibility;
out vec3 ReflectedVector;
out vec3 RefractedVector;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitTangent;

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

uniform float scale;

uniform vec3 viewPos;
uniform float refractivity;

uniform float totalTime;

uniform sampler2D terrainHeightMap;
uniform sampler2D terrainNormalMap;

uniform float terrainDisplacement;
uniform float terrainWidth;
uniform float terrainLength;

uniform float foliageDisplacement;
uniform float foliageNoiseSeed;
uniform float foliageVariation;

const float fogDensity = 0.01;
const float fogGradient = 1.5;

const float windAmplitude = 0.4;
const float windFrequency = 1;

float rand(vec2 c) {
	return fract(sin(dot(c.xy, vec2(12.9898,78.233))) * 34.5453);
}

float rand2(vec2 c) {
	return fract(cos(dot(c.xy, vec2(4824.135,3.62557))) * 32.5453);
}

float random(float seed, float minimum, float maximum) {
    float initial = rand(vec2(seed, rand2(vec2(foliageNoiseSeed, sin(seed)))));
    return minimum + initial * (maximum - minimum);
}

float getFogVisibility(vec4 mPos) {
    vec4 positionRelativeCamera = viewMtx * mPos;
    float dist = length(positionRelativeCamera.xyz);
    return clamp(exp(-pow((dist * fogDensity), fogGradient)), 0.0, 1.0);
}

vec3 getInstancePosition() {
    float xAdd = random(float(gl_InstanceID) / terrainWidth, 0.0, terrainLength);
    float zAdd = random(float(gl_InstanceID + 100) / terrainLength, 0.0, terrainWidth);
    return vec3(position.x + float(xAdd), position.y, position.z + float(zAdd));
}

float getFoliageHeight(vec4 mPos) {
    return position.y + foliageDisplacement + texture(terrainHeightMap, vec2(mPos.x / terrainWidth, mPos.z / terrainLength)).r * terrainDisplacement;
}

vec3 getTerrainNormal(vec4 mPos) {
    return texture(terrainNormalMap, vec2(mPos.x / terrainWidth, mPos.z / terrainLength)).rgb;
}

void main() {
    vec3 instancePosition = getInstancePosition();
    vec3 finalPosition = vec3(instancePosition.x, instancePosition.y, instancePosition.z);

    // Wind simulation
    if(tex.y == 0) {
       finalPosition.x += sin((totalTime + gl_InstanceID) * windFrequency) * sin(totalTime - gl_InstanceID) * windAmplitude;
       finalPosition.z += sin((totalTime - gl_InstanceID) * windFrequency) * windAmplitude;
       finalPosition.y += random(sqrt(gl_InstanceID), -foliageVariation, foliageVariation);
    }

    vec4 mPos = modelMtx * vec4(finalPosition, 1.0);
    finalPosition.y = getFoliageHeight(mPos);

    // Fragment position
    FragPos = vec3(modelMtx * vec4(finalPosition, 1.0));

     // Normal vector
    Normal = mat3(transpose(inverse(modelMtx))) * getTerrainNormal(mPos);

    //	Vertex position 
    gl_Position = projectionMtx * viewMtx * vec4(FragPos, 1.0);
   
    // Texture coordinates
    TexCoords = vec3(tex.x, tex.y, tex.z) / scale;

    // Normal mapping
    vec3 T = normalize(vec3(modelMtx * vec4(tangent,   0.0)));
    vec3 B = normalize(vec3(modelMtx * vec4(bitTangent, 0.0)));
    vec3 N = normalize(vec3(modelMtx * vec4(normal,    0.0)));
    TBN = mat3(T, B, N);

    // Reflection/refraction
    vec3 viewVector = normalize(finalPosition - viewPos);
    ReflectedVector = reflect(viewVector, normal);
    RefractedVector = refract(viewVector, normal, refractivity);

    Visibility = getFogVisibility(mPos);
}
