#version 410

uniform sampler2D screen;
uniform sampler2D bloomInput;

in vec3 TexCoord;
out vec4 FragColor;

uniform float bloomIntensity;

void main() {
    vec4 texColor = texture(screen, TexCoord.xy);
    vec4 bloomColor = texture(bloomInput, TexCoord.xy);

    FragColor = texColor + bloomColor * bloomIntensity;
}