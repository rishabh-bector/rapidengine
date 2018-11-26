#version 410 

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

uniform int flipped;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;

out vec3 texCoord;

void main() {
    if(flipped == 0) {
        texCoord = tex;
    } else {
        texCoord = vec3(1 - tex.x, tex.y, tex.z);
    }

    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
}