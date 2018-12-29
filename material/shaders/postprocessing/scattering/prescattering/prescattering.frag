#version 410

uniform sampler2D screen;

in vec3 TexCoord;
out vec4 FragColor;

uniform vec2 lightPos;

uniform float decay;
uniform float density;
uniform float weight;
uniform float exposure;

const int samples = 100;

void main() {
    vec2 deltaTexCoord = vec2(TexCoord.xy - lightPos);
    vec2 texCoord = TexCoord.xy;
    
    deltaTexCoord *= 1.0 / float(samples) * density;
    
    float illuminationDecay = 2.0;

    for(int i = 0; i < samples; i++) {
        texCoord -= deltaTexCoord;
        vec4 smpl = texture(screen, texCoord);

        smpl *= illuminationDecay * weight;

        FragColor += smpl;

        illuminationDecay *= decay;
    }

    FragColor *= exposure;
}