#version 410

uniform sampler2D screen;

in vec3 TexCoord;
out vec4 FragColor;

uniform float bloomThreshold;

void main() {
    vec4 texColor = texture(screen, TexCoord.xy);
    float brightness = (texColor.r * 0.2126) + (texColor.g * 0.7152) + (texColor.b * 0.0722);

    if(brightness > bloomThreshold) {
        FragColor = texColor;
    } else {
        FragColor = vec4(0.0, 0.0, 0.0, 1.0);
    }
    //FragColor = texColor * brightness;
}