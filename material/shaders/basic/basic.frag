#version 410

uniform float diffuseLevel;

uniform vec4 hue;
uniform float darkness;

uniform sampler2D diffuseMap;
uniform float diffuseMapScale;

uniform sampler2D alphaMap;
uniform float alphaMapLevel;

in vec3 texCoord;
out vec4 outColor;

vec4 calculateDiffuseColor();
vec4 calculateHueColor();
float calculateAlpha();

void main() {
    vec4 finalHue = calculateHueColor();
    vec4 finalDiffuse = calculateDiffuseColor();
    vec4 finalColor = mix(finalHue, finalDiffuse, diffuseLevel);
    
    float finalAlpha = calculateAlpha();

    if(finalColor.a == 0 || finalAlpha < 0.15) {
        discard;
    }

    outColor = vec4(darkness * finalColor.xyz, finalColor.a * finalAlpha);
}

vec4 calculateDiffuseColor() {
    return texture(diffuseMap, texCoord.xy);
}

vec4 calculateHueColor() {
    return (hue / 255);
}

float calculateAlpha() {
    if(alphaMapLevel == 0) {
        return 1;
    }
    return texture(alphaMap, texCoord.xy).x;
}