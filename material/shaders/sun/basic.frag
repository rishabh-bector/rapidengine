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
    vec4 col = calculateDiffuseColor();

    if(col.a < 0.8) {
        discard;
    }
 
    outColor = col;
    scatterColor = col;
}

vec4 calculateDiffuseColor() {
    return texture(diffuseMap, texCoord.xy);
}
