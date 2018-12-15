#version 410

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;

out vec3 TexCoord;

void main() {
    gl_Position = modelMtx * vec4(position.xy, 0, 1.0);
    TexCoord = vec3(tex.x, 1 - tex.y, 0);
}