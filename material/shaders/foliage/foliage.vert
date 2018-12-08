#version 410

out VS_OUT {
    vec3 FragPos;
    vec3 TexCoords;
    vec3 Normal;
    mat3 TBN;
} vs_out;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitTangent;

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;


void main() {
    vec3 finalPosition = vec3(position.x, position.y, position.z);

    if(tex.y == 1) {
       
    }

    //	Vertex position 
    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(finalPosition, 1.0);

    // Normal vector
    vs_out.Normal = mat3(transpose(inverse(modelMtx))) * normal;

    // Fragment position
    vs_out.FragPos =  vec3(modelMtx * vec4(finalPosition, 1.0));

    // Texture coordinates
    vs_out.TexCoords = tex;

    // Normal Mapping
    vec3 T = normalize(vec3(modelMtx * vec4(tangent,   0.0)));
    vec3 B = normalize(vec3(modelMtx * vec4(bitTangent, 0.0)));
    vec3 N = normalize(vec3(modelMtx * vec4(normal,    0.0)));
    vs_out.TBN = mat3(T, B, N);
}
