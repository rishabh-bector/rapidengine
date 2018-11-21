#version 410

out vec3 Normal;
out vec3 FragPos;
out vec3 TexCoords;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;
layout (location = 2) in vec3 normal;

uniform float textureScale;

// Child copying
uniform vec3 copyingEnabled;
layout (location = 3) in vec3 copyPosition;

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

void main() {
    //	Vertex position
    if(copyingEnabled.x > 0) {
        gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position + gl_InstanceID, 1.0);
    } else {
        gl_Position = projectionMtx * viewMtx * modelMtx * vec4(position, 1.0);
    }

    // Normal vector
    Normal = mat3(transpose(inverse(modelMtx))) * normal;

    // Fragment position
    FragPos = vec3(modelMtx * vec4(position, 1.0));

    // Texture coordinates
    TexCoords = tex / textureScale;
}