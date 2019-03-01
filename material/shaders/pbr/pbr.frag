#version 410

struct PointLight {
    vec3 position;
    
    float constant;
    float linear;
    float quadratic;
    
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

struct DirLight {
    vec3 direction;
    
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

out vec4 FragColor;

in vec3 TexCoords;
in vec3 FragPos;
in vec3 Normal;
in mat3 TBN;
in vec3 ReflectedVector;
in vec3 RefractedVector;

// PBR Material
uniform sampler2D albedoMap;
uniform sampler2D normalMap;
uniform sampler2D heightMap;
uniform sampler2D metallicMap;
uniform sampler2D roughnessMap;
uniform sampler2D aoMap;

uniform float normalScalar;
uniform float metallicScalar;
uniform float roughnessScalar;
uniform float aoScalar;

uniform float parallaxDisplacement;
uniform float vertexDisplacement;
uniform float scale;

// lights
#define MAX_LIGHTS 10
uniform int numPointLights;
uniform PointLight pointLights[MAX_LIGHTS];
uniform DirLight dirLight;

// Camera
uniform vec3 viewPos;

// Constant
const float PI = 3.14159265359;

// ----------------------------------------------------------------------------

float DistributionGGX(vec3 N, vec3 H, float roughness)
{
    float a = roughness*roughness;
    float a2 = a*a;
    float NdotH = max(dot(N, H), 0.0);
    float NdotH2 = NdotH*NdotH;

    float nom   = a2;
    float denom = (NdotH2 * (a2 - 1.0) + 1.0);
    denom = PI * denom * denom;

    return nom / max(denom, 0.001); // prevent divide by zero for roughness=0.0 and NdotH=1.0
}

// ----------------------------------------------------------------------------

float GeometrySchlickGGX(float NdotV, float roughness)
{
    float r = (roughness + 1.0);
    float k = (r*r) / 8.0;

    float nom   = NdotV;
    float denom = NdotV * (1.0 - k) + k;

    return nom / denom;
}

// ----------------------------------------------------------------------------

float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness)
{
    float NdotV = max(dot(N, V), 0.0);
    float NdotL = max(dot(N, L), 0.0);
    float ggx2 = GeometrySchlickGGX(NdotV, roughness);
    float ggx1 = GeometrySchlickGGX(NdotL, roughness);

    return ggx1 * ggx2;
}

// ----------------------------------------------------------------------------

vec3 fresnelSchlick(float cosTheta, vec3 F0)
{
    return F0 + (1.0 - F0) * pow(1.0 - cosTheta, 5.0);
}

// ----------------------------------------------------------------------------

float calculateHeight(vec2 uv) {
    return 1.0 - texture(heightMap, uv).r;
}

vec2 parallaxMapping(vec3 viewDir) {
    // number of depth layers
    float minLayers = 8.0;
    float maxLayers = 32.0;
    float numLayers = 32;//mix(minLayers, maxLayers, abs(dot(vec3(0.0, 0.0, 1.0), viewDir)));  

    // calculate the size of each layer
    float layerDepth = 1.0 / numLayers;

    // depth of current layer
    float currentLayerDepth = 0.0;

    // the amount to shift the texture coordinates per layer (from vector P)
    vec2 P = viewDir.xy * parallaxDisplacement; 
    vec2 deltaTexCoords = P / numLayers;

    // get initial values
    vec2  currentTexCoords     = TexCoords.xy;
    float currentDepthMapValue = calculateHeight(currentTexCoords);
  
    while(currentLayerDepth < currentDepthMapValue) {
        // shift texture coordinates along direction of P
        currentTexCoords -= deltaTexCoords;

        // get depthmap value at current texture coordinates
        currentDepthMapValue = calculateHeight(currentTexCoords); 

        // get depth of next layer
        currentLayerDepth += layerDepth;  
    }

    // get texture coordinates before collision (reverse operations)
    vec2 prevTexCoords = currentTexCoords + deltaTexCoords;

    // get depth after and before collision for linear interpolation
    float afterDepth  = currentDepthMapValue - currentLayerDepth;
    float beforeDepth = calculateHeight(prevTexCoords) - currentLayerDepth + layerDepth;
 
    // interpolation of texture coordinates
    float weight = afterDepth / (afterDepth - beforeDepth);
    vec2 finalTexCoords = prevTexCoords * weight + currentTexCoords * (1.0 - weight);

    return finalTexCoords;  
}

// ----------------------------------------------------------------------------

vec2 getUVs() {
    vec2 uvs = TexCoords.xy;

    /*if(parallaxDisplacement > 0) {
        vec3 tangentViewDir = normalize((viewPos - FragPos) * TBN);
        uvs = parallaxMapping(tangentViewDir);
        
        if(uvs.x > 1.0 / scale || uvs.y > 1.0 / scale || uvs.x < 0.0 || uvs.y < 0.0) {
            discard;
        }
    }*/
    
    return uvs;
}

// ----------------------------------------------------------------------------

vec3 getNormal(vec2 uvs) {
    vec3 norm = vec3(0, 0, 1);

    if(normalScalar > 0) {
        norm = normalize(texture(normalMap, uvs).rgb * normalScalar);
        norm = normalize((norm * 2.0) - 0.5);
        norm = normalize(TBN * norm);
    } else {
        norm = normalize(Normal);
    }

    return norm;
}

// ----------------------------------------------------------------------------

void main() {		
    vec2 uvs = getUVs();

    vec3 albedo = texture(albedoMap, uvs).rgb;
    vec3 normal = getNormal(uvs);

    float metallic = texture(metallicMap, uvs).r * metallicScalar;
    float roughness = (1 - texture(roughnessMap, uvs).r) * roughnessScalar;
    float ao = texture(aoMap, uvs).r * aoScalar;

    //albedo = vec3(0.2, 0.2, 0.8);
    metallic = metallicScalar;
    // roughness = roughnessScalar;
    //ao = aoScalar;

    if(metallicScalar < 0) {
        metallic = abs(metallicScalar);
    }
    if(roughnessScalar < 0) {
        roughness = abs(roughnessScalar);
    }

    vec3 N = normalize(normal);
    vec3 V = normalize(viewPos - FragPos);

    // calculate reflectance at normal incidence; if dia-electric (like plastic) use F0 
    // of 0.04 and if it's a metal, use the albedo color as F0 (metallic workflow)    
    vec3 F0 = vec3(0.04); 
    F0 = mix(F0, albedo, metallic);

    // reflectance equation
    vec3 Lo = vec3(0.0);

    // directional light case
    vec3 radiance = dirLight.diffuse;
    vec3 L = normalize(dirLight.direction);
    vec3 H = normalize(V + L);
    float NDF = DistributionGGX(N, H, roughness);   
    float G   = GeometrySmith(N, V, L, roughness);      
    vec3 F    = fresnelSchlick(clamp(dot(H, V), 0.0, 1.0), F0);
    vec3 nominator    = NDF * G * F; 
    float denominator = 4 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0);
    vec3 specular = nominator / max(denominator, 0.001);
    vec3 kS = F;
    vec3 kD = vec3(1.0) - kS;
    kD *= 1.0 - metallic;	  
    float NdotL = max(dot(N, L), 0.0);        
    Lo += (kD * albedo / PI + specular) * radiance * NdotL;


    for(int i = 0; i < numPointLights; i++) 
    {
        // calculate per-light radiance
        float dist = length(pointLights[i].position - FragPos);
        float attenuation = 1.0 / (dist * dist);
        vec3 radiance = pointLights[i].diffuse * attenuation;

        vec3 L = normalize(pointLights[i].position - FragPos);
        vec3 H = normalize(V + L);

        // Cook-Torrance BRDF
        float NDF = DistributionGGX(N, H, roughness);   
        float G   = GeometrySmith(N, V, L, roughness);      
        vec3 F    = fresnelSchlick(clamp(dot(H, V), 0.0, 1.0), F0);
           
        vec3 nominator    = NDF * G * F; 
        float denominator = 4 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0);
        vec3 specular = nominator / max(denominator, 0.001); // prevent divide by zero for NdotV=0.0 or NdotL=0.0
        
        // kS is equal to Fresnel
        vec3 kS = F;

        // for energy conservation, the diffuse and specular light can't
        // be above 1.0 (unless the surface emits light); to preserve this
        // relationship the diffuse component (kD) should equal 1.0 - kS.
        vec3 kD = vec3(1.0) - kS;

        // multiply kD by the inverse metalness such that only non-metals 
        // have diffuse lighting, or a linear blend if partly metal (pure metals
        // have no diffuse light).
        kD *= 1.0 - metallic;	  

        // scale light by NdotL
        float NdotL = max(dot(N, L), 0.0);        

        // add to outgoing radiance Lo
        Lo += (kD * albedo / PI + specular) * radiance * NdotL;  // note that we already multiplied the BRDF by the Fresnel (kS) so we won't multiply by kS again
    }   
    
    // ambient lighting (note that the next IBL tutorial will replace 
    // this ambient lighting with environment lighting).
    vec3 ambient = vec3(0.03) * albedo * ao;

    vec3 color = ambient + Lo;

    // HDR tonemapping
    color = color / (color + vec3(1.0));

    // gamma correct
    //color = pow(color, vec3(1.0/2.2)); 

    FragColor = vec4(color, 1.0);
}