#version 410

uniform sampler2D screen;
uniform sampler2D scatterInput;

in vec3 TexCoord;
out vec4 FragColor;

void main() {
    vec4 texColor = texture(screen, TexCoord.xy);
    vec4 scatterColor = texture(scatterInput, TexCoord.xy);

    FragColor = mix(texColor, scatterColor, 0.5);
}