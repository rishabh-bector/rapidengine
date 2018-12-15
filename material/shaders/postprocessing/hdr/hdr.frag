#version 410

uniform sampler2D screen;

in vec3 TexCoord;
out vec4 FragColor;

void main() {
    const float exposure = 1.0;

    // HDR Color from unclamped floating point buffer
    vec3 hdrColor = texture(screen, TexCoord.xy).rgb;

    // Reinhard tone mapping
    vec3 mapped = vec3(1.0) - exp(-hdrColor * exposure);

    // Gamma correction
    //mapped = pow(mapped, vec3(1.0 / gamma));

    FragColor = vec4(mapped, 1.0);
}