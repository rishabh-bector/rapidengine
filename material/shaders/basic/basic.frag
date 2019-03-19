#version 410

uniform float diffuseLevel;

uniform vec4 hue;
uniform float darkness;

uniform sampler2D diffuseMap;
uniform float diffuseMapScale;

uniform sampler2D alphaMap;
uniform float alphaMapLevel;

uniform float scatterLevel;

in vec3 texCoord;

layout(location = 0) out vec4 outColor;
layout(location = 1) out vec4 scatterColor;

vec4 calculateDiffuseColor();
vec4 calculateHueColor();
float calculateAlpha();

void main() {
    vec4 finalHue = calculateHueColor();
    vec4 finalDiffuse = calculateDiffuseColor();
    vec4 finalColor = mix(finalHue, finalDiffuse, diffuseLevel);
    
    float finalAlpha = calculateAlpha();

    if(finalColor.a < 0.01 || finalAlpha < 0.2) {
        discard;
    }

    outColor = vec4(darkness * finalColor.xyz, finalColor.a);
    scatterColor = outColor * scatterLevel;
}

vec4 calculateDiffuseColor() {
    return texture(diffuseMap, texCoord.xy);
}

vec4 calculateHueColor() {
    return (hue / 255);
}

float calculateAlpha() {
    if(alphaMapLevel == 0) {
        return 1.0;
    }
    return texture(alphaMap, texCoord.xy).x;
}