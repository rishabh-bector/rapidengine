#version 410 

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;

out vec3 texCoord;

void main() {
    texCoord = tex;
    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
}