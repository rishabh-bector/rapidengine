#version 410

uniform sampler2D texture0;
uniform float darkness;
uniform int transparencyEnabled;
uniform sampler2D transparencyMap;

in vec3 texCoord;

out vec4 outColor;

void main() {

    if(texture(texture0, texCoord.xy).a < 0.5) {
        discard;
    }

    vec3 col = vec3(0, 0, 0);
    float d = darkness;

    if(transparencyEnabled == 1) {
        if(texture(transparencyMap, texCoord.xy).x == 0.0f) {
            discard;
        } else if(texture(transparencyMap, texCoord.xy).x != 1.0f) {
            col = vec3(0, 0, 0);
        } else {
            col = texture(texture0, texCoord.xy).xyz;
        }
    } else {
        col = texture(texture0, texCoord.xy).xyz;
    }

    outColor = vec4(d * col, 1);
}