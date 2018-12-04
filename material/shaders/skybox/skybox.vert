#version 410

layout (location = 0) in vec3 position;

out vec3 texCoord;

uniform mat4 viewMtx;
uniform mat4 projectionMtx;
uniform mat4 modelMtx;

void main() {
    texCoord = position;
    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
}