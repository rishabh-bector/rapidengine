#version 410

out vec4 FragColor;

in vec3 texCoord;

uniform samplerCube cubeDiffuseMap;

void main() {
    FragColor = texture(cubeDiffuseMap, texCoord);
}