#version 410

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;

out vec2 blurTextureCoords[11];

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

uniform float fboHeight;

void main() {
    gl_Position = vec4(position.xy, 0, 1.0);
    
    float pixelSize = 1.0 / fboHeight;

    for(int i = -5; i <= 5; i++) {
        blurTextureCoords[i+5] = vec2(tex.x, 1 - tex.y) + vec2(0, pixelSize * i);
    }
}