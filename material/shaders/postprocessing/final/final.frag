#version 410

uniform sampler2D screen;

in vec3 TexCoord;
out vec4 FragColor;

void main() {
    vec4 texColor = texture(screen, TexCoord.xy);

    FragColor = texColor;
}