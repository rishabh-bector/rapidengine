#version 410

uniform sampler2D screen;

in vec2 blurTextureCoords[11];

out vec4 FragColor;

void main() {
    vec4 out_color = vec4(0.0);

    out_color += texture(screen, blurTextureCoords[0]) * 0.0093;
    out_color += texture(screen, blurTextureCoords[1]) * 0.028002;
    out_color += texture(screen, blurTextureCoords[2]) * 0.065984;
    out_color += texture(screen, blurTextureCoords[3]) * 0.121703;
    out_color += texture(screen, blurTextureCoords[4]) * 0.175713;
    out_color += texture(screen, blurTextureCoords[5]) * 0.198596;
    out_color += texture(screen, blurTextureCoords[6]) * 0.175713;
    out_color += texture(screen, blurTextureCoords[7]) * 0.121703;
    out_color += texture(screen, blurTextureCoords[8]) * 0.065984;
    out_color += texture(screen, blurTextureCoords[9]) * 0.028002;
    out_color += texture(screen, blurTextureCoords[10]) * 0.0093;

    FragColor = out_color;
}