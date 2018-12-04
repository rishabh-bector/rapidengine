#version 410

uniform sampler2D screen;

in vec3 TexCoord;
out vec4 FragColor;

void main() {
    vec4 texColor = texture(screen, TexCoord.xy);
    //FragColor = vec4(0.1, 0.2, 0.5, 1.0);


    FragColor = texColor;
}