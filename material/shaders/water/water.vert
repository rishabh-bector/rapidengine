#version 410

out vec3 FragPos;
out vec3 TexCoords;
out mat3 TBN;
out vec3 Normal;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitTangent;

uniform float scale;
uniform float displacement;

uniform sampler2D heightMap;

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;


uniform float totalTime;
const float waterAmplitude = 0.5;
const float waterFrequency = 1;

const float pi = 3.14159;
const float waterHeight = 0.2;
const int numWaves = 4;
const float amplitude[4] = float[](0.1, 0.05, 0.025, 0.01);
const float wavelength[4] = float[](4, 3, 2, 1);
const float speed[4] = float[](0.1, 0.2, 0.3, 0.4);
const vec2 direction[4] = vec2[](vec2(cos(0.1), sin(0.1)), vec2(cos(0.5), sin(0.5)), vec2(cos(1.4), sin(1.4)), vec2(cos(2.9), sin(2.9)));

float getDisplacement();

float wave(int i, float x, float y) {
    float frequency = 2*pi/wavelength[i];
    float phase = speed[i] * frequency;
    float theta = dot(direction[i], vec2(x, y));
    return amplitude[i] * sin(theta * frequency + totalTime * phase);
}

float waveHeight(float x, float y) {
    float height = 0.0;
    for (int i = 0; i < numWaves; ++i)
        height += wave(i, x, y);
    return height;
}

float dWavedx(int i, float x, float y) {
    float frequency = 2*pi/wavelength[i];
    float phase = speed[i] * frequency;
    float theta = dot(direction[i], vec2(x, y));
    float A = amplitude[i] * direction[i].x * frequency;
    return A * cos(theta * frequency + totalTime * phase);
}

float dWavedy(int i, float x, float y) {
    float frequency = 2*pi/wavelength[i];
    float phase = speed[i] * frequency;
    float theta = dot(direction[i], vec2(x, y));
    float A = amplitude[i] * direction[i].y * frequency;
    return A * cos(theta * frequency + totalTime * phase);
}

vec3 waveNormal(float x, float y) {
    float dx = 0.0;
    float dy = 0.0;
    for (int i = 0; i < numWaves; ++i) {
        dx += dWavedx(i, x, y);
        dy += dWavedy(i, x, y);
    }
    vec3 n = vec3(-dx, -dy, 1.0);
    return normalize(n);
}

void main() {
    vec3 finalPosition = vec3(position.x, position.y, position.z);

    //finalPosition.y += waterAmplitude * sin((totalTime + position.x) * waterFrequency);
    //finalPosition.y += waterAmplitude * cos((totalTime + position.z) * waterFrequency);
    finalPosition.y += 1 + waveHeight(position.x, position.z);

    //	Vertex position 
    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(finalPosition, 1.0);

    // Normal vector
    Normal = waveNormal(position.x, position.z);

    // Fragment position
    FragPos =  vec3(modelMtx * vec4(finalPosition, 1.0));

    // Texture coordinates
    TexCoords = tex / scale;

    // Normal Mapping
    vec3 T = normalize(vec3(modelMtx * vec4(tangent,   0.0)));
    vec3 B = normalize(vec3(modelMtx * vec4(bitTangent, 0.0)));
    vec3 N = normalize(vec3(modelMtx * vec4(normal,    0.0)));
    TBN = mat3(T, B, N);
}

float getDisplacement() {
    return texture(heightMap, tex.xy / scale).x * displacement;
}
